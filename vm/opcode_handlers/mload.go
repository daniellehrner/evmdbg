package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type MLoadOpCode struct{}

func (*MLoadOpCode) Execute(v *vm.DebuggerVM) error {
	// MLOAD requires one value on the stack.
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the address from the stack.
	addr, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// Read the word from memory at the given address and push it onto the stack.
	return v.Stack().Push(v.Memory().ReadWord(addr.Uint64()))
}
