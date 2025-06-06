package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SubOpCode struct{}

func (*SubOpCode) Execute(v *vm.DebuggerVM) error {
	// The SUB opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Perform the subtraction
	return v.Push(new(big.Int).Sub(a, b))
}
