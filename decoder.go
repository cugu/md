package md

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type decoder struct {
	config *config
}

func newDecoder(option ...Option) *decoder {
	config := &config{
		parser:                goldmark.DefaultParser(),
		disallowUnknownFields: false,
	}

	for _, opt := range option {
		opt(config)
	}

	return &decoder{config: config}
}

func (d *decoder) unmarshal(md []byte, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("v must be a non-nil pointer")
	}

	dst := rv.Elem()
	src := reflect.TypeOf(v).Elem()

	txt, block := firstBlock(d.config.parser, md)

	fieldIndex := 0

	for block != nil {
		if fieldIndex >= dst.NumField() {
			return d.handleExtraMarkdownBlocks(block)
		}

		var err error
		if fieldIndex, block, err = d.assignBlockContentToField(fieldIndex, src, block, txt, dst); err != nil {
			return err
		}
	}

	return handleAdditionalFields(fieldIndex, dst, src)
}

func (d *decoder) assignBlockContentToField(fieldIndex int, src reflect.Type, block ast.Node, txt text.Reader, dst reflect.Value) (nextFieldIndex int, nextBlock ast.Node, err error) {
	tag, omitempty, exists := tag(src.Field(fieldIndex))

	// if the field is not tagged, skip it.
	if !exists {
		return fieldIndex + 1, block.NextSibling(), nil
	}

	element := block.Kind()

	tagName, ok := elementMap[element]

	// if the field is tagged but the element is not supported
	if !ok {
		if d.config.disallowUnknownFields {
			return 0, nil, fmt.Errorf("unexpected element: %s", element.String())
		}

		return fieldIndex, block.NextSibling(), nil
	}

	if tag == tagName {
		dst.Field(fieldIndex).SetString(content(block, txt))
	} else { // the field is tagged but the element is not supported
		if omitempty { // if the field has omitempty tag, skip it.
			return fieldIndex + 1, block, nil
		} else {
			return 0, nil, fmt.Errorf("unexpected %s: %s", tagName, content(block, txt))
		}
	}

	return fieldIndex + 1, block.NextSibling(), nil
}

func (d *decoder) handleExtraMarkdownBlocks(block ast.Node) error {
	var extraBlocks []string

	for ; block != nil; block = block.NextSibling() {
		element := block.Kind()
		tagName, ok := elementMap[element]
		if !ok {
			if d.config.disallowUnknownFields {
				return fmt.Errorf("unexpected element: %s", element.String())
			}

			continue
		}

		extraBlocks = append(extraBlocks, tagName)
	}

	if len(extraBlocks) > 0 {
		return fmt.Errorf("extra blocks: %s", strings.Join(extraBlocks, ", "))
	}

	return nil
}
