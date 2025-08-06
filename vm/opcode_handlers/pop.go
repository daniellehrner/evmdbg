package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type PopOpCode struct{}

func (*PopOpCode) Execute(v *vm.DebuggerVM) error {
	// The POP opcode requires at least one item on the stack to pop.
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the top item from the stack.
	_, err := v.Stack().Pop()

	return err
}
