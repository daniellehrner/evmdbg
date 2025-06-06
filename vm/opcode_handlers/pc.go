package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type PCOpCode struct{}

func (*PCOpCode) Execute(v *vm.DebuggerVM) error {
	// Push the current program counter (PC) onto the stack.
	return v.PushUint64(v.PC - 1)
}
