// Package md implements functions to parse markdown into Go structs, similar to how JSON is parsed into Go structs.
package md

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const mdTag = "md"

type config struct {
	disallowUnknownFields bool
	parser                parser.Parser
}

// Option is a functional option type for the Unmarshal function.
type Option func(*config)

// WithParser is a functional option that allows you to set the parser to be used by Unmarshal.
func WithParser(p parser.Parser) Option {
	return func(m *config) {
		m.parser = p
	}
}

// WithDisallowUnknownFields is a functional option that allows you to disallow unknown fields in the markdown.
func WithDisallowUnknownFields() Option {
	return func(m *config) {
		m.disallowUnknownFields = true
	}
}

// Unmarshal takes a byte slice of markdown (md) and a non-nil pointer (v) as arguments.
// It parses the markdown into sections and assigns the content of each section to the corresponding field in v.
// The function uses struct tags to map markdown sections to fields in v.
// If a required section (one without the 'omitempty' option in its tag) is missing from the markdown, the function returns an error.
// If an unexpected section kind is encountered, the function also returns an error.
func Unmarshal(md []byte, v any, option ...Option) error {
	config := &config{
		parser:                goldmark.DefaultParser(),
		disallowUnknownFields: false,
	}

	for _, opt := range option {
		opt(config)
	}

	return unmarshal(md, v, config)
}

func unmarshal(md []byte, v any, settings *config) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("v must be a non-nil pointer")
	}

	sv := rv.Elem()
	t := reflect.TypeOf(v).Elem()

	txt, section := firstSection(settings.parser, md)

	for i := range sv.NumField() {
		tag, omitempty, ok := tag(t.Field(i))
		if !ok {
			continue
		}

		if section == nil {
			if omitempty {
				continue
			}

			return fmt.Errorf("missing %s", tag)
		}

		kind := section.Kind()

		kindMap := map[ast.NodeKind]string{
			ast.KindBlockquote:      "blockquote",
			ast.KindFencedCodeBlock: "code_block",
			ast.KindHeading:         "heading",
			ast.KindParagraph:       "paragraph",
			ast.KindThematicBreak:   "thematic_break",
		}

		tagName, ok := kindMap[kind]
		if !ok {
			if settings.disallowUnknownFields {
				return fmt.Errorf("unexpected kind: %s", kind.String())
			}

			continue
		}

		if tag == tagName {
			sv.Field(i).SetString(content(section, txt))
		} else if omitempty {
			continue
		} else {
			return fmt.Errorf("unexpected %s: %s", tagName, content(section, txt))
		}

		section = section.NextSibling()
	}

	return nil
}

func content(v ast.Node, txt text.Reader) string {
	if v.Type() != ast.TypeBlock {
		return ""
	}

	s := ""

	for i := range v.Lines().Len() {
		line := v.Lines().At(i)
		s += string(line.Value(txt.Source()))
	}

	for c := v.FirstChild(); c != nil; c = c.NextSibling() {
		s += content(c, txt)
	}

	return s
}

func firstSection(p parser.Parser, md []byte) (text.Reader, ast.Node) {
	txt := text.NewReader(md)
	node := p.Parse(txt)
	return txt, node.FirstChild()
}

func tag(field reflect.StructField) (value string, omitempty bool, exists bool) {
	tag, ok := field.Tag.Lookup(mdTag)
	if !ok {
		return "", false, false
	}

	if strings.HasSuffix(tag, ",omitempty") {
		omitempty = true

		tag = strings.TrimSuffix(tag, ",omitempty")
	}

	return tag, omitempty, true
}
