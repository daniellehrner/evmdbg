package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type TLoadOpCode struct{}

func (*TLoadOpCode) Execute(v *vm.DebuggerVM) error {
	// TLOAD requires one value on the stack (the storage slot)
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the slot from the stack
	slot, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// Read the transient storage at the specified slot and push it onto the stack
	value := v.ReadTransientStorage(slot)
	return v.Stack().Push(value)
}
