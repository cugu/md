// This example directory demonstrates how to use `github.com/cugu/md` to write test cases for a markdown parser in
// markdown.
// Each test case consists of a markdown file with a description of the test and the expected output in markdown and HTML.
package markdowntest

import (
	"bytes"
	"embed"
	"io/fs"
	"reflect"
	"testing"

	"github.com/yuin/goldmark"

	"github.com/cugu/md"
)

//go:embed test_*.md
var TestFiles embed.FS

type TestCase struct {
	Title           string `md:"heading"`
	Description     string `md:"paragraph"`
	MarkdownHeading string `md:"heading"`
	MarkdownText    string `md:"code_block"`
	HTMLHeading     string `md:"heading"`
	HTMLText        string `md:"code_block"`
}

func TestMarkdownToHTML(t *testing.T) {
	// list all test files
	entries, err := fs.ReadDir(TestFiles, ".")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			// read markdown file
			markdown, err := fs.ReadFile(TestFiles, entry.Name())
			if err != nil {
				t.Fatal(err)
			}

			// unmarshal markdown file using github.com/cugu/md package
			var testCase TestCase
			if err := md.Unmarshal(markdown, &testCase); err != nil {
				t.Errorf("Unmarshal() error = %v", err)
			}

			// convert markdown to HTML
			buf := &bytes.Buffer{}
			if err := goldmark.Convert([]byte(testCase.MarkdownText), buf); err != nil {
				return
			}

			// compare the HTML output with the expected HTML output
			if !reflect.DeepEqual(buf.String(), testCase.HTMLText) {
				t.Errorf("HTML = %v, want %v", buf.String(), testCase.HTMLText)
			}
		})
	}
}
