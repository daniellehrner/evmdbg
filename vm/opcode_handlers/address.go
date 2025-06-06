package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type AddressOpCode struct{}

func (*AddressOpCode) Execute(v *vm.DebuggerVM) error {
	return v.PushBytes(v.Context.Address[:])
}
