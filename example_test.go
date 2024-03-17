package md_test

import (
	_ "embed"
	"fmt"

	"github.com/cugu/md"
)

type Text struct {
	Title               string `md:"heading"`
	Description         string `md:"paragraph"`
	OptionalTitle       string `md:"heading,omitempty"`
	OptionalDescription string `md:"paragraph,omitempty"`
}

func ExampleUnmarshal() {
	markdown := []byte("# Title\n\nA short description.")

	var text Text
	if err := md.Unmarshal(markdown, &text); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(text.Description)
	// Output: A short description.
}
