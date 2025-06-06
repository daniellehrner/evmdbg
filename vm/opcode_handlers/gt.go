package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type GtOpCode struct{}

func (*GtOpCode) Execute(v *vm.DebuggerVM) error {
	// The GT opcode compares two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Compare the two values and push 1 if a > b, otherwise push 0
	if a.Cmp(b) > 0 {
		return v.PushUint64(1)
	}

	return v.PushUint64(0)
}
