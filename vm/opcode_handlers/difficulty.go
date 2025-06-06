package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type DifficultyOpCode struct{}

func (*DifficultyOpCode) Execute(v *vm.DebuggerVM) error {
	// If the block difficulty is not set, return 0
	if v.Context.Block == nil || v.Context.Block.Difficulty == nil {
		return v.Push(new(big.Int))
	}

	return v.Push(v.Context.Block.Difficulty)
}
