package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type DupOpCode struct {
	N int // 1 to 16
}

func (d *DupOpCode) Execute(v *vm.DebuggerVM) error {
	// The DUP opcodes require the stack to have at least N elements.
	if err := v.RequireStack(d.N); err != nil {
		return err
	}

	// Peek the N-th element from the top of the stack (0-indexed).
	val, err := v.Stack.Peek(d.N - 1)
	if err != nil {
		return err
	}

	// Push a copy of the N-th element onto the stack.
	return v.Stack.Push(new(uint256.Int).Set(val))
}
