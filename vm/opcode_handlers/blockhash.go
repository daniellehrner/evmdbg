package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type BlockHashOpCode struct{}

func (*BlockHashOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("blockhash op code requires the execution context to be set")
	}

	if v.Context.Block == nil {
		return fmt.Errorf("blockhash op code requires block context to be set")
	}

	// BLOCKHASH requires one value on the stack (the block number)
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the block number from the stack
	blockNumber, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	var blockHash *uint256.Int

	// BLOCKHASH only returns hashes for the most recent 256 blocks
	// and only for blocks that are strictly less than the current block number
	currentBlockNumber := v.Context.Block.Number
	requestedBlock := blockNumber.Uint64()

	if requestedBlock >= currentBlockNumber ||
		(currentBlockNumber > 256 && requestedBlock <= currentBlockNumber-256) {
		// Return 0 for:
		// - Future blocks (requestedBlock >= currentBlockNumber)
		// - Blocks too far in the past (more than 256 blocks ago)
		blockHash = uint256.NewInt(0)
	} else if v.StateProvider != nil {
		// Get block hash from state provider
		hashBytes := v.StateProvider.GetBlockHash(requestedBlock)
		blockHash = new(uint256.Int).SetBytes(hashBytes[:])
	} else {
		// If no state provider, return 0
		blockHash = uint256.NewInt(0)
	}

	// Push the block hash onto the stack
	return v.Push(blockHash)
}
