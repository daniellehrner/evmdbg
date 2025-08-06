package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type DivOpCode struct{}

func (*DivOpCode) Execute(v *vm.DebuggerVM) error {
	// The DIV opcode divides the top two stack items.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Perform the division and push the result back onto the stack.
	return v.Push(new(uint256.Int).Div(a, b))
}
