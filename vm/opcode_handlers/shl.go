package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ShlOpCode struct{}

func (*ShlOpCode) Execute(v *vm.DebuggerVM) error {
	// The SHL opcode requires two values on the stack
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

	// EVM word size is 256 bits: apply shift left operation
	n := shift.Uint64()
	if n >= 256 {
		return v.Push(new(big.Int)) // shift too far, zero
	}

	// Perform the left shift operation
	res := new(big.Int).Lsh(value, uint(n))

	// mask to 256 bits
	res.And(res, uint256Mask)

	return v.Push(res)
}
