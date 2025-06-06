package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type CallerOpCode struct{}

func (*CallerOpCode) Execute(v *vm.DebuggerVM) error {
	return v.PushBytes(v.Context.Caller[:])
}
