package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type CoinbaseOpCode struct{}

func (*CoinbaseOpCode) Execute(v *vm.DebuggerVM) error {
	if v.Context.Block == nil {
		return v.Push(new(big.Int))
	}
	return v.PushBytes(v.Context.Block.Coinbase[:])
}
