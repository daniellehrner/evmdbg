package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ExpOpCode struct{}

func (*ExpOpCode) Execute(v *vm.DebuggerVM) error {
	// The EXP opcode requires at least two items on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the base and exponent from the stack.
	base, exponent, err := v.Pop2()
	if err != nil {
		return err
	}

	result := new(big.Int).Exp(base, exponent, uint256)
	return v.Push(result)
}
