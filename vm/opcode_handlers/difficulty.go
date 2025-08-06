package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type DifficultyOpCode struct{}

func (*DifficultyOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("address op code requires the execution context to be set")
	}

	// If the block difficulty is not set, return 0
	if v.Context.Block == nil || v.Context.Block.Difficulty == nil {
		return v.Push(new(uint256.Int))
	}

	return v.Push(v.Context.Block.Difficulty)
}
