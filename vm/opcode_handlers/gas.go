package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type GasOpCode struct{}

func (*GasOpCode) Execute(v *vm.DebuggerVM) error {
	return v.PushUint64(v.Context.Gas)
}
