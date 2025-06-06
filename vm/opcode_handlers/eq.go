package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type EqOpCode struct{}

func (*EqOpCode) Execute(v *vm.DebuggerVM) error {
	// The EQ requires two values on the stack to compare.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Compare the two values and push 1 if they are equal
	if a.Cmp(b) == 0 {
		return v.PushUint64(1)
	}

	// If they are not equal, push 0 onto the stack.
	return v.PushUint64(0)
}
