package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type ReturnOpCode struct{}

func (*ReturnOpCode) Execute(v *vm.DebuggerVM) error {
	// The RETURN opcode requires two values on the stack:
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	offset, size, err := v.Pop2()
	if err != nil {
		return err
	}

	// The return value is the memory content from the specified offset and size.
	v.ReturnValue = v.Memory.Read(int(offset.Uint64()), int(size.Uint64()))

	// Set the stopped flag to true to indicate that the execution should stop.
	v.Stopped = true

	return nil
}
