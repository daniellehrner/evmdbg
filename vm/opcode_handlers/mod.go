package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
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

	// If the divisor is zero, return zero as per EVM rules.
	if b.Sign() == 0 {
		return v.Push(new(big.Int)) // EVM: mod by zero = 0
	}

	// Perform the modulo operation.
	return v.Push(new(big.Int).Mod(a, b))
}
