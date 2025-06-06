package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ChainIdOpCode struct{}

func (*ChainIdOpCode) Execute(v *vm.DebuggerVM) error {
	// If the block is nil or the chain ID is nil, push a zero value.
	if v.Context.Block == nil || v.Context.Block.ChainID == nil {
		return v.Push(new(big.Int))
	}

	// Otherwise, push the chain ID onto the stack.
	return v.Push(v.Context.Block.ChainID)
}
