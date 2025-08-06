package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type BlobHashOpCode struct{}

func (*BlobHashOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("blobhash op code requires the execution context to be set")
	}

	if v.Context.Block == nil {
		return fmt.Errorf("blobhash op code requires block context to be set")
	}

	// BLOBHASH requires one value on the stack (the index)
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the index from the stack
	indexInt, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	index := indexInt.Uint64()
	var hash *uint256.Int

	// Check if index is within bounds of available blob hashes
	if index < uint64(len(v.Context.Block.BlobHashes)) {
		// Return the blob hash at the specified index
		hashBytes := v.Context.Block.BlobHashes[index]
		hash = new(uint256.Int).SetBytes(hashBytes[:])
	} else {
		// Return 0 if index is out of bounds
		hash = uint256.NewInt(0)
	}

	// Push the blob hash onto the stack
	return v.Push(hash)
}
