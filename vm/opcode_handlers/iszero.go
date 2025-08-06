package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type IsZeroOpCode struct{}

func (*IsZeroOpCode) Execute(v *vm.DebuggerVM) error {
	// ISZERO requires one value on the stack.
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the top value from the stack.
	x, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// Check if the value is zero and push 1 if true, otherwise push 0.
	if x.Sign() == 0 {
		return v.PushUint64(1)
	}

	return v.PushUint64(0)
}
