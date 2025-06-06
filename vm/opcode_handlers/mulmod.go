package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
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

	// If the modulus is zero, return zero as per EVM rules.
	if m.Sign() == 0 {
		return v.Push(big.NewInt(0))
	}

	// Perform the multiplication and then take modulo.
	prod := new(big.Int).Mul(a, b)
	prod.Mod(prod, m)

	return v.Push(prod)
}
