package source_test

import (
	"github.com/m0ssc0de/scale.go/source"
	"testing"
)

func TestLoadTypeRegistry(t *testing.T) {
	source.LoadTypeRegistry([]byte(source.BaseType))
}
