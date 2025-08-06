package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type ModOpCode struct{}

func (*ModOpCode) Execute(v *vm.DebuggerVM) error {
	// MOD requires two values on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Perform the modulo operation.
	return v.Push(new(uint256.Int).Mod(a, b))
}
