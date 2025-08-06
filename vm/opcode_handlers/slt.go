package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type SLTOpCode struct{}

func (*SLTOpCode) Execute(v *vm.DebuggerVM) error {
	// The SLT opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Check if a < b when interpreted as signed 256-bit integers
	if a.Slt(b) {
		return v.Push(uint256.NewInt(1))
	}

	// If a is not less than b, push 0 onto the stack
	return v.Push(new(uint256.Int).Clear())
}
