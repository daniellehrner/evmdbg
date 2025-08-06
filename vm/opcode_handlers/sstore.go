package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type SStoreOpCode struct{}

func (*SStoreOpCode) Execute(v *vm.DebuggerVM) error {
	// Check if we're in a static call context
	currentFrame := v.CurrentFrame()
	if currentFrame != nil && currentFrame.IsStatic {
		return vm.ErrStaticCallStateChange
	}

	// The SSTORE opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	slot, val, err := v.Pop2()
	if err != nil {
		return err
	}

	// write the value to the storage at the specified slot
	v.WriteStorage(slot, val)

	return nil
}
