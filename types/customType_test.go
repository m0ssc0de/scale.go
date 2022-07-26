package types_test

import (
	"github.com/m0ssc0de/scale.go/types"
	"testing"
)

func TestRegCustomTypes(t *testing.T) {
	types.RuntimeType{}.Reg()
}
