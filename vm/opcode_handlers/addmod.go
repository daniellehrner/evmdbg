package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
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
		return v.Push(big.NewInt(0))
	}

	// Perform the addition and then take the modulus.
	sum := new(big.Int).Add(a, b)
	sum.Mod(sum, m)

	// Push the result back onto the stack.
	return v.Push(sum)
}
