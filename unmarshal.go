package md

// Unmarshal takes a byte slice of markdown (md) and a non-nil pointer (v) as arguments.
// It parses the markdown into blocks and assigns the content of each block to the corresponding field in v.
// The function uses struct tags to map markdown block to fields in v.
// If a required block (one without the 'omitempty' option in its tag) is missing from the markdown, the function returns an error.
// If an unexpected block element is encountered, the function also returns an error.
func Unmarshal(md []byte, v any, option ...Option) error {
	return newDecoder(option...).unmarshal(md, v)
}
