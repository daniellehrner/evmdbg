package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type RevertOpCode struct{}

func (*RevertOpCode) Execute(v *vm.DebuggerVM) error {
	// The REVERT opcode requires at least two items on the stack:
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	offset, size, err := v.Pop2()
	if err != nil {
		return err
	}

	// The REVERT opcode sets the return value to the memory content from the specified offset and size.
	v.ReturnValue = v.Memory().Read(int(offset.Uint64()), int(size.Uint64()))

	// Set the reverted and stopped flags to true to indicate that the execution has been reverted.
	v.Reverted = true

	// Set the stopped flag to true to indicate that the execution should stop.
	v.Stopped = true

	return nil
}
