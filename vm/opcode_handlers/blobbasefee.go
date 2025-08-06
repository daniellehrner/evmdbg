package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type BlobBaseFeeOpCode struct{}

func (*BlobBaseFeeOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("blobbasefee op code requires the execution context to be set")
	}

	if v.Context.Block == nil {
		return fmt.Errorf("blobbasefee op code requires block context to be set")
	}

	// Push the blob base fee onto the stack
	return v.Push(v.Context.Block.BlobBaseFee)
}
