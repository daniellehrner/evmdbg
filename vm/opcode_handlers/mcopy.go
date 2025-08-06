package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type MCopyOpCode struct{}

func (*MCopyOpCode) Execute(v *vm.DebuggerVM) error {
	// MCOPY requires three values on the stack: destOffset, offset, size
	if err := v.RequireStack(3); err != nil {
		return err
	}

	// Pop the top three items from the stack.
	destOffset, offset, size, err := v.Pop3()
	if err != nil {
		return err
	}

	// If size is 0, do nothing
	if size.IsZero() {
		return nil
	}

	// Read from source memory location
	data := v.Memory().Read(int(offset.Uint64()), int(size.Uint64()))

	// Write to destination memory location
	v.Memory().Write(int(destOffset.Uint64()), data)

	return nil
}
