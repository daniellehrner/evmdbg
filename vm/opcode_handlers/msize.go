package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type MSizeOpCode struct{}

func (*MSizeOpCode) Execute(v *vm.DebuggerVM) error {
	return v.PushUint64(uint64(v.Memory.Size()))
}
