package md_test

import (
	"reflect"
	"testing"

	"github.com/cugu/md"
)

type T1 struct {
	Title string `md:"heading"`
}

type T2 struct {
	Title string `md:"heading"`
	Body  string `md:"paragraph,omitempty"`
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		md []byte
		v  any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{"nil", args{[]byte("test"), nil}, nil, true},
		{"non-nil", args{[]byte("test"), &struct{}{}}, &struct{}{}, false},
		{"test struct", args{[]byte("# test"), &T1{}}, &T1{Title: "test"}, false},
		{"omitempty", args{[]byte("# test\n\nbody"), &T2{}}, &T2{Title: "test", Body: "body"}, false},
		{"omitempty empty", args{[]byte("# test"), &T2{}}, &T2{Title: "test"}, false},
		{"link", args{[]byte("# [test](/test)"), &T1{}}, &T1{Title: "[test](/test)"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := tt.args.v

			if err := md.Unmarshal(tt.args.md, in); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(in, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", in, tt.want)
			}
		})
	}
}
