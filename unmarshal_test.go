package md_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yuin/goldmark"

	"github.com/cugu/md"
)

type UnmarshalTestCase[T any] struct {
	name    string
	args    UnmarshalTestArgs[T]
	want    T
	wantErr error
}

func (tc UnmarshalTestCase[T]) Name() string {
	return tc.name
}

func (tc UnmarshalTestCase[T]) Run(t *testing.T) {
	runUnmarshalTest(t, tc.args.md, tc.args.v, tc.args.options, tc.want, tc.wantErr)
}

func runUnmarshalTest(t *testing.T, markdown string, v any, options []md.Option, want any, wantErr error) {
	err := md.Unmarshal([]byte(markdown), v, options...)

	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("Unmarshal() error = %v, wantErr %v", err, wantErr)
	}

	if !reflect.DeepEqual(v, want) {
		t.Errorf("Unmarshal() = %v, want %v", v, want)
	}
}

type UnmarshalTestArgs[T any] struct {
	md      string
	v       T
	options []md.Option
}

func TestUnmarshalNil(t *testing.T) {
	runUnmarshalTest(t, "test", nil, nil, nil, errors.New("v must be a non-nil pointer"))
}

func TestUnmarshal(t *testing.T) {
	type OnlyTitle struct {
		Title string `md:"heading1"`
	}

	type TitleAndDescription struct {
		Title       string `md:"heading1"`
		Description string `md:"paragraph"`
	}

	type UntaggedTitle struct {
		Title string
	}

	type MissingFields struct{}

	type OptionalTitle struct {
		Title string `md:"heading1,omitempty"`
	}

	tests := []Test{
		UnmarshalTestCase[*OnlyTitle]{
			name:    "use custom parser",
			args:    UnmarshalTestArgs[*OnlyTitle]{"# Title", &OnlyTitle{}, []md.Option{md.WithParser(goldmark.DefaultParser())}},
			want:    &OnlyTitle{Title: "Title"},
			wantErr: nil,
		},
		UnmarshalTestCase[*TitleAndDescription]{
			name:    "disallow unknown fields",
			args:    UnmarshalTestArgs[*TitleAndDescription]{"\n# Title\n\n- A list item.\n\nA short description.\n", &TitleAndDescription{}, []md.Option{md.WithDisallowUnknownFields()}},
			want:    &TitleAndDescription{Title: "Title"},
			wantErr: errors.New("unexpected element: List"),
		},
		UnmarshalTestCase[*TitleAndDescription]{
			name:    "allow unknown fields",
			args:    UnmarshalTestArgs[*TitleAndDescription]{"\n# Title\n\n- A list item.\n\nA short description.\n", &TitleAndDescription{}, nil},
			want:    &TitleAndDescription{Title: "Title", Description: "A short description."},
			wantErr: nil,
		},
		UnmarshalTestCase[*UntaggedTitle]{
			name:    "untagged struct",
			args:    UnmarshalTestArgs[*UntaggedTitle]{"# Title", &UntaggedTitle{}, nil},
			want:    &UntaggedTitle{},
			wantErr: nil,
		},
		UnmarshalTestCase[*OnlyTitle]{
			name:    "mismatch between field and markdown block element",
			args:    UnmarshalTestArgs[*OnlyTitle]{"A short description.", &OnlyTitle{}, nil},
			want:    &OnlyTitle{},
			wantErr: errors.New("unexpected paragraph: A short description."),
		},
		UnmarshalTestCase[*MissingFields]{
			name:    "more markdown blocks than fields",
			args:    UnmarshalTestArgs[*MissingFields]{"# Title", &MissingFields{}, nil},
			want:    &MissingFields{},
			wantErr: errors.New("extra blocks: heading1"),
		},
		UnmarshalTestCase[*OnlyTitle]{
			name:    "more fields than blocks",
			args:    UnmarshalTestArgs[*OnlyTitle]{"", &OnlyTitle{}, nil},
			want:    &OnlyTitle{},
			wantErr: errors.New("missing heading1"),
		},
		UnmarshalTestCase[*OptionalTitle]{
			name:    "more fields than blocks, but omitempty tag set",
			args:    UnmarshalTestArgs[*OptionalTitle]{"", &OptionalTitle{}, nil},
			want:    &OptionalTitle{},
			wantErr: nil,
		},
		UnmarshalTestCase[*OnlyTitle]{
			name:    "more markdown blocks than fields, and disallow unknown fields",
			args:    UnmarshalTestArgs[*OnlyTitle]{"\n# Title\n\n- A list item.\n", &OnlyTitle{}, []md.Option{md.WithDisallowUnknownFields()}},
			want:    &OnlyTitle{Title: "Title"},
			wantErr: errors.New("unexpected element: List"),
		},
		UnmarshalTestCase[*OnlyTitle]{
			name:    "more fields than blocks, and allow unknown fields",
			args:    UnmarshalTestArgs[*OnlyTitle]{"\n# Title\n\n- A list item.\n", &OnlyTitle{}, nil},
			want:    &OnlyTitle{Title: "Title"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name(), tt.Run)
	}
}
