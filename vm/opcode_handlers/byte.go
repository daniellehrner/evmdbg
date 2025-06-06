package opcode_handlers

import (
	"math/big"

	"github.com/daniellehrner/evmdbg/vm"
)

type ByteOpCode struct{}

func (*ByteOpCode) Execute(v *vm.DebuggerVM) error {
	// The BYTE opcode requires at least two items on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	shift, word, err := v.Pop2()
	if err != nil {
		return err
	}

	// If the shift is not a valid uint64 or is greater than or equal to 32,
	if !shift.IsUint64() || shift.Uint64() >= 32 {
		return v.Push(big.NewInt(0))
	}

	// BYTE returns the (31 - shift)th byte
	pos := 31 - shift.Uint64()
	byteVal := word.Bytes()

	// If the byte value is less than 32 bytes, pad it with leading zeros
	if len(byteVal) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(byteVal):], byteVal)
		byteVal = padded
	}

	return v.Push(new(big.Int).SetUint64(uint64(byteVal[pos])))
}
