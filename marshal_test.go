package md_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/cugu/md"
)

type MarshalTestCase[T any] struct {
	name    string
	args    MarshalTestArgs[T]
	want    string
	wantErr error
}

func (tc MarshalTestCase[T]) Name() string {
	return tc.name
}

func (tc MarshalTestCase[T]) Run(t *testing.T) {
	runMarshalTest(t, tc.args.v, tc.want, tc.wantErr)
}

func runMarshalTest(t *testing.T, v any, want string, wantErr error) {
	b, err := md.Marshal(v)

	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("Marshal() error = %v, wantErr %v", err, wantErr)
	}

	if !bytes.Equal(b, []byte(want)) {
		t.Errorf("Marshal() = %v, want %v", string(b), want)
	}
}

type MarshalTestArgs[T any] struct {
	v T
}

func TestMarshalNil(t *testing.T) {
	runMarshalTest(t, nil, "", errors.New("v must be a non-nil pointer"))
}

func TestMarshalNonPointer(t *testing.T) {
	i := 0

	runMarshalTest(t, &i, "", errors.New("v must be a pointer to a struct"))
}

func TestMarshal(t *testing.T) {
	type OnlyTitle struct {
		Title string `md:"heading1"`
	}

	type NoTag struct {
		Title string
	}

	type Empty struct{}

	type AllTags struct {
		Blockquote    string `md:"blockquote"`
		CodeBlock     string `md:"code_block"`
		Title         string `md:"heading1"`
		Title2        string `md:"heading2"`
		Title3        string `md:"heading3"`
		Title4        string `md:"heading4"`
		Title5        string `md:"heading5"`
		Title6        string `md:"heading6"`
		Paragraph     string `md:"paragraph"`
		ThematicBreak string `md:"thematic_break"`
	}

	type UnknownTag struct {
		Unknown string `md:"unknow"`
	}

	type Omitempty struct {
		Title string `md:"heading1,omitempty"`
	}

	type NonString struct {
		Title int `md:"heading1"`
	}

	tests := []Test{
		MarshalTestCase[*OnlyTitle]{
			name:    "use custom parser",
			args:    MarshalTestArgs[*OnlyTitle]{&OnlyTitle{Title: "Title"}},
			want:    "# Title\n\n",
			wantErr: nil,
		},
		MarshalTestCase[*NoTag]{
			name:    "no tag",
			args:    MarshalTestArgs[*NoTag]{&NoTag{Title: "Title"}},
			want:    "",
			wantErr: nil,
		},
		MarshalTestCase[*Empty]{
			name:    "empty struct",
			args:    MarshalTestArgs[*Empty]{&Empty{}},
			want:    "",
			wantErr: nil,
		},
		MarshalTestCase[*AllTags]{
			name: "all tags",
			args: MarshalTestArgs[*AllTags]{&AllTags{
				Blockquote:    "Blockquote",
				CodeBlock:     "CodeBlock",
				Title:         "Title",
				Title2:        "Title2",
				Title3:        "Title3",
				Title4:        "Title4",
				Title5:        "Title5",
				Title6:        "Title6",
				Paragraph:     "Paragraph",
				ThematicBreak: "ThematicBreak",
			}},
			want:    "> Blockquote\n\n```\nCodeBlock\n```\n\n# Title\n\n## Title2\n\n### Title3\n\n#### Title4\n\n##### Title5\n\n###### Title6\n\nParagraph\n\n---\n\n",
			wantErr: nil,
		},
		MarshalTestCase[*UnknownTag]{
			name:    "unknow tag",
			args:    MarshalTestArgs[*UnknownTag]{&UnknownTag{Unknown: "Unknow"}},
			want:    "",
			wantErr: errors.New("unsupported tag: unknow"),
		},
		MarshalTestCase[*Omitempty]{
			name:    "omitempty",
			args:    MarshalTestArgs[*Omitempty]{&Omitempty{Title: ""}},
			want:    "",
			wantErr: nil,
		},
		MarshalTestCase[*NonString]{
			name:    "non string",
			args:    MarshalTestArgs[*NonString]{&NonString{Title: 1}},
			want:    "",
			wantErr: errors.New("field Title must be a string"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name(), tt.Run)
	}
}
