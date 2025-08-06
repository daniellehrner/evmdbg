package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"

	"github.com/holiman/uint256"
)

type AddOpCode struct{}

func (*AddOpCode) Execute(v *vm.DebuggerVM) error {
	// The ADD opcode requires at least two items on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	return v.Push(new(uint256.Int).Add(a, b))
}
