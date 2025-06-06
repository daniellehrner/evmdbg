package opcode_handlers

import (
	"math/big"

	"github.com/daniellehrner/evmdbg/vm"
)

type ByteOpCode struct{}

func (*ByteOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	shift, word, err := v.Pop2()
	if err != nil {
		return err
	}

	if !shift.IsUint64() || shift.Uint64() >= 32 {
		return v.Push(big.NewInt(0))
	}

	// BYTE returns the (31 - shift)th byte
	pos := 31 - shift.Uint64()
	byteVal := word.Bytes()
	if len(byteVal) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(byteVal):], byteVal)
		byteVal = padded
	}
	return v.Push(new(big.Int).SetUint64(uint64(byteVal[pos])))
}
