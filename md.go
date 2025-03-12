// Package md implements functions to parse markdown into Go structs, similar to how JSON is parsed into Go structs.
package md

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func tagName(node ast.Node) (string, bool) {
	switch n := node.(type) {
	case *ast.Blockquote:
		return "blockquote", true
	case *ast.FencedCodeBlock:
		return "code_block", true
	case *ast.Heading:
		return "heading" + strconv.Itoa(n.Level), true
	case *ast.Paragraph:
		return "paragraph", true
	case *ast.ThematicBreak:
		return "thematic_break", true
	default:
		return "", false
	}
}

func handleAdditionalFields(fieldIndex int, dst reflect.Value, src reflect.Type) error {
	var missingFields []string

	for ; fieldIndex < dst.NumField(); fieldIndex++ {
		tag, omitempty, _ := tag(src.Field(fieldIndex))
		if omitempty {
			continue
		}

		missingFields = append(missingFields, tag)
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing %s", strings.Join(missingFields, ", "))
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

func firstBlock(p parser.Parser, md []byte) (text.Reader, ast.Node) {
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
