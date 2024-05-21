package md_test

import (
	"fmt"

	"github.com/cugu/md"
)

type Text struct {
	Title       string `md:"heading"`
	Description string `md:"paragraph"`
}

const example = `
# Title

A short description.
`

func ExampleUnmarshal() {
	var text Text
	if err := md.Unmarshal([]byte(example), &text); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(text.Title)
	fmt.Println(text.Description)
	// Output:
	// Title
	// A short description.
}

type ExtendedText struct {
	Title       string `md:"heading"`
	Quote       string `md:"blockquote"`
	Description string `md:"paragraph"`
	Break       string `md:"thematic_break"`
	CodeBlock   string `md:"code_block"`
}

const exampleExtended = "\n# Title\n\n> A quote.\n\nA short description.\n\n---\n\n```\ncode block\n```\n"

func ExampleUnmarshal_extended() {
	var text ExtendedText
	if err := md.Unmarshal([]byte(exampleExtended), &text); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(text.Title)
	fmt.Println(text.Quote)
	fmt.Println(text.Description)
	fmt.Println(text.Break)
	fmt.Println(text.CodeBlock)
	// Output:
	// Title
	// A quote.
	// A short description.
	//
	// code block
}

type OmitemptyText struct {
	Title       string `md:"heading"`
	Quote       string `md:"blockquote,omitempty"`
	Description string `md:"paragraph"`
}

const exampleOmitempty = `
# Title

A short description.
`

func ExampleUnmarshal_omitempty() {
	var text OmitemptyText
	if err := md.Unmarshal([]byte(exampleOmitempty), &text); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(text.Title)
	fmt.Println(text.Quote)
	fmt.Println(text.Description)
	// Output:
	// Title
	//
	// A short description.
}
