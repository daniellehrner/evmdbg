package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ChainIdOpCode struct{}

func (*ChainIdOpCode) Execute(v *vm.DebuggerVM) error {
	if v.Context.Block == nil || v.Context.Block.ChainID == nil {
		return v.Push(new(big.Int))
	}
	return v.Push(v.Context.Block.ChainID)
}
