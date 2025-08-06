package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type SLoadOpCode struct{}

func (*SLoadOpCode) Execute(v *vm.DebuggerVM) error {
	// If the sign bit is 0 (positive number):
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// If the sign bit is 0 (positive number):
	slot, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// Read the storage at the specified slot and push it onto the stack.
	return v.Stack().Push(v.ReadStorage(slot))
}
