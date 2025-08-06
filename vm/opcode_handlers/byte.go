package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
)

type ByteOpCode struct{}

func (*ByteOpCode) Execute(v *vm.DebuggerVM) error {
	// The BYTE opcode requires at least two items on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	shift, word, err := v.Pop2()
	if err != nil {
		return err
	}

	return v.Push(word.Byte(shift))
}
