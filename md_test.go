package md_test

import (
	"testing"
)

type Test interface {
	Name() string
	Run(t *testing.T)
}
