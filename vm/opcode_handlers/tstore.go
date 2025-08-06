package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type TStoreOpCode struct{}

func (*TStoreOpCode) Execute(v *vm.DebuggerVM) error {
	// Check if we're in a static call context
	currentFrame := v.CurrentFrame()
	if currentFrame != nil && currentFrame.IsStatic {
		return vm.ErrStaticCallStateChange
	}

	// TSTORE requires two values on the stack (slot and value)
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	slot, val, err := v.Pop2()
	if err != nil {
		return err
	}

	// Write the value to the transient storage at the specified slot
	v.WriteTransientStorage(slot, val)

	return nil
}
