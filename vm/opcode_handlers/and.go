package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type AndOpCode struct{}

func (*AndOpCode) Execute(v *vm.DebuggerVM) error {
	// The AND opcode requires at least two items on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Perform the bitwise AND operation and push the result back onto the stack.
	return v.Push(new(uint256.Int).And(a, b))
}
