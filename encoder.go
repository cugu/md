package md

import (
	"bytes"
	"fmt"
	"reflect"
)

type encoder struct{}

func newEncoder() *encoder {
	return &encoder{}
}

func (e *encoder) marshal(v any) ([]byte, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return nil, fmt.Errorf("v must be a non-nil pointer")
	}

	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("v must be a pointer to a struct")
	}

	var buf bytes.Buffer
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() != reflect.String {
			return nil, fmt.Errorf("field %s must be a string", typ.Field(i).Name)
		}

		tag, omitempty, exists := tag(typ.Field(i))

		if !exists {
			continue
		}

		if omitempty && field.String() == "" {
			continue
		}

		switch tag {
		case "blockquote":
			buf.WriteString(fmt.Sprintf("> %s\n\n", field.String()))
		case "code_block":
			buf.WriteString(fmt.Sprintf("```\n%s\n```\n\n", field.String()))
		case "heading1":
			buf.WriteString(fmt.Sprintf("# %s\n\n", field.String()))
		case "heading2":
			buf.WriteString(fmt.Sprintf("## %s\n\n", field.String()))
		case "heading3":
			buf.WriteString(fmt.Sprintf("### %s\n\n", field.String()))
		case "heading4":
			buf.WriteString(fmt.Sprintf("#### %s\n\n", field.String()))
		case "heading5":
			buf.WriteString(fmt.Sprintf("##### %s\n\n", field.String()))
		case "heading6":
			buf.WriteString(fmt.Sprintf("###### %s\n\n", field.String()))
		case "paragraph":
			buf.WriteString(fmt.Sprintf("%s\n\n", field.String()))
		case "thematic_break":
			buf.WriteString("---\n\n")
		default:
			return nil, fmt.Errorf("unsupported tag: %s", tag)
		}
	}

	return buf.Bytes(), nil
}
