package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type CodeSizeOpCode struct{}

func (*CodeSizeOpCode) Execute(v *vm.DebuggerVM) error {
	return v.PushUint64(uint64(len(v.Code)))
}
