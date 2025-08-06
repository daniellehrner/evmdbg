package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type AddressOpCode struct{}

func (*AddressOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("address op code requires the execution context to be set")
	}

	return v.PushBytes(v.Context.Address[:])
}
