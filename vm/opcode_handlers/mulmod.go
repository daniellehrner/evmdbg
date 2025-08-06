package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type MulModOpCode struct{}

func (*MulModOpCode) Execute(v *vm.DebuggerVM) error {
	// MULMOD requires three values on the stack.
	if err := v.RequireStack(3); err != nil {
		return err
	}

	// Pop the top three items from the stack.
	a, b, m, err := v.Pop3()
	if err != nil {
		return err
	}

	return v.Push(new(uint256.Int).MulMod(a, b, m))
}
