package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type MulOpCode struct{}

func (*MulOpCode) Execute(v *vm.DebuggerVM) error {
	// MUL requires two values on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Perform the multiplication operation.
	return v.Push(new(big.Int).Mul(a, b))
}
