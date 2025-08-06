package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type NotOpCode struct{}

func (*NotOpCode) Execute(v *vm.DebuggerVM) error {
	// NOT requires one value on the stack.
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the top item from the stack.
	x, err := v.Stack.Pop()
	if err != nil {
		return err
	}

	return v.Push(new(uint256.Int).Not(x))
}
