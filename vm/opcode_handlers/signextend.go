package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type SignExtendOpCode struct{}

func (*SignExtendOpCode) Execute(v *vm.DebuggerVM) error {
	// The SIGNEXTEND opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	k, val, err := v.Pop2()
	if err != nil {
		return err
	}

	return v.Push(new(uint256.Int).ExtendSign(val, k))
}
