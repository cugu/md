package md_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/yuin/goldmark"

	"github.com/cugu/md"
)

type Test interface {
	Name() string
	Run(t *testing.T)
}

type TestCase[T any] struct {
	name    string
	args    TestArgs[T]
	want    T
	wantErr error
}

func (tc TestCase[T]) Name() string {
	return tc.name
}

func (tc TestCase[T]) Run(t *testing.T) {
	runTest(t, tc.args.md, tc.args.v, tc.args.options, tc.want, tc.wantErr)
}

func runTest(t *testing.T, markdown string, v any, options []md.Option, want any, wantErr error) {
	err := md.Unmarshal([]byte(markdown), v, options...)

	if !reflect.DeepEqual(err, wantErr) {
		t.Errorf("Unmarshal() error = %v, wantErr %v", err, wantErr)
	}

	if !reflect.DeepEqual(v, want) {
		t.Errorf("Unmarshal() = %v, want %v", v, want)
	}
}

type TestArgs[T any] struct {
	md      string
	v       T
	options []md.Option
}

func TestUnmarshalNil(t *testing.T) {
	runTest(t, "test", nil, nil, nil, errors.New("v must be a non-nil pointer"))
}

func TestUnmarshal(t *testing.T) {
	type OnlyTitle struct {
		Title string `md:"heading"`
	}

	type TitleAndDescription struct {
		Title       string `md:"heading"`
		Description string `md:"paragraph"`
	}

	type UntaggedTitle struct {
		Title string
	}

	type MissingFields struct{}

	type OptionalTitle struct {
		Title string `md:"heading,omitempty"`
	}

	tests := []Test{
		TestCase[*OnlyTitle]{
			name:    "use custom parser",
			args:    TestArgs[*OnlyTitle]{"# Title", &OnlyTitle{}, []md.Option{md.WithParser(goldmark.DefaultParser())}},
			want:    &OnlyTitle{Title: "Title"},
			wantErr: nil,
		},
		TestCase[*TitleAndDescription]{
			name:    "disallow unknown fields",
			args:    TestArgs[*TitleAndDescription]{"\n# Title\n\n- A list item.\n\nA short description.\n", &TitleAndDescription{}, []md.Option{md.WithDisallowUnknownFields()}},
			want:    &TitleAndDescription{Title: "Title"},
			wantErr: errors.New("unexpected element: List"),
		},
		TestCase[*TitleAndDescription]{
			name:    "allow unknown fields",
			args:    TestArgs[*TitleAndDescription]{"\n# Title\n\n- A list item.\n\nA short description.\n", &TitleAndDescription{}, nil},
			want:    &TitleAndDescription{Title: "Title", Description: "A short description."},
			wantErr: nil,
		},
		TestCase[*UntaggedTitle]{
			name:    "untagged struct",
			args:    TestArgs[*UntaggedTitle]{"# Title", &UntaggedTitle{}, nil},
			want:    &UntaggedTitle{},
			wantErr: nil,
		},
		TestCase[*OnlyTitle]{
			name:    "mismatch between field and markdown block element",
			args:    TestArgs[*OnlyTitle]{"A short description.", &OnlyTitle{}, nil},
			want:    &OnlyTitle{},
			wantErr: errors.New("unexpected paragraph: A short description."),
		},
		TestCase[*MissingFields]{
			name:    "more markdown blocks than fields",
			args:    TestArgs[*MissingFields]{"# Title", &MissingFields{}, nil},
			want:    &MissingFields{},
			wantErr: errors.New("extra blocks: heading"),
		},
		TestCase[*OnlyTitle]{
			name:    "more fields than blocks",
			args:    TestArgs[*OnlyTitle]{"", &OnlyTitle{}, nil},
			want:    &OnlyTitle{},
			wantErr: errors.New("missing heading"),
		},
		TestCase[*OptionalTitle]{
			name:    "more fields than blocks, but omitempty tag set",
			args:    TestArgs[*OptionalTitle]{"", &OptionalTitle{}, nil},
			want:    &OptionalTitle{},
			wantErr: nil,
		},
		TestCase[*OnlyTitle]{
			name:    "more markdown blocks than fields, and disallow unknown fields",
			args:    TestArgs[*OnlyTitle]{"\n# Title\n\n- A list item.\n", &OnlyTitle{}, []md.Option{md.WithDisallowUnknownFields()}},
			want:    &OnlyTitle{Title: "Title"},
			wantErr: errors.New("unexpected element: List"),
		},
		TestCase[*OnlyTitle]{
			name:    "more fields than blocks, and allow unknown fields",
			args:    TestArgs[*OnlyTitle]{"\n# Title\n\n- A list item.\n", &OnlyTitle{}, nil},
			want:    &OnlyTitle{Title: "Title"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name(), tt.Run)
	}
}
