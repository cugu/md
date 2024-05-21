package md

import "github.com/yuin/goldmark/parser"

const mdTag = "md"

type config struct {
	disallowUnknownFields bool
	parser                parser.Parser
}

// Option is a functional option type for the Unmarshal function.
type Option func(*config)

// WithParser is a functional option that allows you to set the parser to be used by Unmarshal.
func WithParser(p parser.Parser) Option {
	return func(m *config) {
		m.parser = p
	}
}

// WithDisallowUnknownFields is a functional option that allows you to disallow unknown fields in the markdown.
func WithDisallowUnknownFields() Option {
	return func(m *config) {
		m.disallowUnknownFields = true
	}
}
