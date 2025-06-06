package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
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

	// If the divisor is zero, return 0 as per EVM specification
	if b.Sign() == 0 {
		return v.Push(big.NewInt(0))
	}

	// Define the half of the uint256 value for signed comparison
	sa := new(big.Int).Set(a)
	if sa.Cmp(uint256Half) >= 0 {
		sa.Sub(sa, uint256)
	}

	// Handle the case where the divisor is negative
	sb := new(big.Int).Set(b)
	if sb.Cmp(uint256Half) >= 0 {
		sb.Sub(sb, uint256)
	}

	// Perform the signed modulo operation
	mod := new(big.Int).Mod(sa, sb)
	if mod.Sign() < 0 {
		mod.Add(mod, uint256)
	}

	return v.Push(mod)
}
