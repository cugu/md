# md

Parse markdown into Go structs, similar to encoding/json.

## Example

Parse this README.md file into a struct.

```go
package main

import (
	_ "embed"
	"fmt"

	"github.com/cugu/md"
)

type Text struct {
	Title           string `md:"heading"`
	Description     string `md:"paragraph"`
	OptionalTitle   string `md:"heading,omitempty"`
	mainDescription string `md:"paragraph,omitempty"`
}

func main() {
	markdown := []byte("# Title\n\nA short description.")

	var text Text
	if err := md.Unmarshal(markdown, &text); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(text.Description)
	// Output: A short description.
}
```