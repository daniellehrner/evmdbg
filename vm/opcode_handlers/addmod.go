package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type AddModOpCode struct{}

func (*AddModOpCode) Execute(v *vm.DebuggerVM) error {
	// The ADDMOD opcode requires at least three items on the stack:
	if err := v.RequireStack(3); err != nil {
		return err
	}

	// Pop the top three items from the stack.
	a, b, m, err := v.Pop3()
	if err != nil {
		return err
	}

	// If the modulus is zero, return zero as per EVM specification.
	if m.Sign() == 0 {
		return v.Push(uint256.NewInt(0))
	}

	// Push the result back onto the stack.
	return v.Push(new(uint256.Int).AddMod(a, b, m))
}
