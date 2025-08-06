package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type MStoreOpCode struct{}

func (*MStoreOpCode) Execute(v *vm.DebuggerVM) error {
	// MSTORE requires two values on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	addr, val, err := v.Pop2()
	if err != nil {
		return err
	}

	// Write the value to memory at the given address.
	v.Memory().WriteWord(addr.Uint64(), val)

	return nil
}
