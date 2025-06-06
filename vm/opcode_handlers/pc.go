package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type PCOpCode struct{}

func (*PCOpCode) Execute(v *vm.DebuggerVM) error {
	// PC was already incremented before Execute
	return v.PushUint64(v.PC - 1)
}
