package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type SModOpCode struct{}

func (*SModOpCode) Execute(v *vm.DebuggerVM) error {
	// The SMOD opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	return v.Push(new(uint256.Int).SMod(a, b))
}
