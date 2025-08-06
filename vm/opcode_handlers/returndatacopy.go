package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type ReturnDataCopyOpCode struct{}

func (*ReturnDataCopyOpCode) Execute(v *vm.DebuggerVM) error {
	// The RETURNDATACOPY opcode requires at least 3 items on the stack:
	if err := v.RequireStack(3); err != nil {
		return err
	}

	// Pop the top three items from the stack.
	memOffset, offset, size, err := v.Pop3()
	if err != nil {
		return err
	}

	returnData := v.ReturnData()
	start := offset.Uint64()
	end := start + size.Uint64()

	// Ensure the start and end indices are within bounds of the return data.
	if end > uint64(len(returnData)) {
		return fmt.Errorf("RETURNDATACOPY out of bounds: %d > %d", end, len(returnData))
	}

	// Write the specified portion of the return data to memory.
	data := returnData[start:end]
	v.Memory().Write(int(memOffset.Uint64()), data)

	return nil
}
