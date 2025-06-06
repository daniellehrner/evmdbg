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

	start := offset.Uint64()
	end := start + size.Uint64()

	// Ensure the start and end indices are within bounds of the return value.
	if end > uint64(len(v.ReturnValue)) {
		return fmt.Errorf("RETURNDATACOPY out of bounds: %d > %d", end, len(v.ReturnValue))
	}

	// Write the specified portion of the return value to memory.
	data := v.ReturnValue[start:end]
	v.Memory.Write(int(memOffset.Uint64()), data)

	return nil
}
