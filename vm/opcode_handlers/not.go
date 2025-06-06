package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type NotOpCode struct{}

func (*NotOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(1); err != nil {
		return err
	}
	x, err := v.Stack.Pop()
	if err != nil {
		return err
	}

	// EVM word size is 256 bits: apply ^ over 32-byte mask
	mask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	return v.Push(new(big.Int).Xor(x, mask))
}
