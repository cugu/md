package md

// Marshal takes any value (v) as an argument and returns a byte slice of markdown (md).
// It uses the struct tags of v to generate markdown blocks.
// The function returns an error if the value is not a struct or if a field in the struct is not supported.
func Marshal(v any) ([]byte, error) {
	return newEncoder().marshal(v)
}
