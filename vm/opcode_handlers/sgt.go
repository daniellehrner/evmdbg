package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type SGTOpCode struct{}

func (*SGTOpCode) Execute(v *vm.DebuggerVM) error {
	// The SGT opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	if a.Sgt(b) {
		return v.Push(uint256.NewInt(1))
	}

	return v.Push(uint256.NewInt(0))
}
