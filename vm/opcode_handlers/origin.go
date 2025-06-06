package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type OriginOpCode struct{}

func (*OriginOpCode) Execute(v *vm.DebuggerVM) error {
	return v.PushBytes(v.Context.Origin[:])
}
