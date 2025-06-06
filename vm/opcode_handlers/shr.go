package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ShrOpCode struct{}

func (*ShrOpCode) Execute(v *vm.DebuggerVM) error {
	// The SHR opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	shift, value, err := v.Pop2()
	if err != nil {
		return err
	}

	// if the shift is greater than 256 bits, the result is zero
	if shift.BitLen() > 256 {
		return v.Push(new(big.Int)) // result is zero
	}

	// EVM word size is 256 bits: apply shift right operation
	n := shift.Uint64()

	// If n is greater than or equal to 256, we can directly return zero.
	if n >= 256 {
		return v.Push(new(big.Int)) // shift too far, zero
	}

	// Perform the right shift operation
	res := new(big.Int).Rsh(value, uint(n))

	return v.Push(res)
}
