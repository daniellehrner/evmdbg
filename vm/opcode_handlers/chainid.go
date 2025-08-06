package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type ChainIdOpCode struct{}

func (*ChainIdOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("chain id op code requires the execution context to be set")
	}

	// If the block is nil or the chain ID is nil, push a zero value.
	if v.Context.Block == nil || v.Context.Block.ChainID == nil {
		return v.Push(new(uint256.Int))
	}

	// Otherwise, push the chain ID onto the stack.
	return v.Push(v.Context.Block.ChainID)
}
