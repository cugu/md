# Markdown Test Example

This example directory demonstrates how to use `github.com/cugu/md` to write test cases for a markdown parser in
markdown.
Each test case consists of a markdown file with a description of the test and the expected output in markdown and HTML.
See [test_1.md](test_1.md) and [test_2.md](test_2.md) for examples.

[convert_test.go](convert_test.go) contains the test code that reads the markdown files and compares the parsed output
with the expected output.