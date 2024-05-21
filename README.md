# md

Parse markdown into Go structs, similar to encoding/json.

Supports parsing markdown into structs with the following fields:

- `heading`: A Markdown heading (e.g. `# Title`).
- `paragraph`: A Markdown paragraph (e.g. `A short description.`).
- `blockquote`: A Markdown blockquote (e.g. `> A blockquote.`).
- `thematic_break`: A Markdown thematic break (e.g. `---`).
- `code_block`: A Markdown code block (e.g. ```` ```go\nfunc main() {}``` ````).

## Example

Parse a markdown file into a struct.

```go
package main

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

func main() {
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
```