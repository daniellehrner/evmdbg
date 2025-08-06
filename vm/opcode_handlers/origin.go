package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type OriginOpCode struct{}

func (*OriginOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("origin op code requires the execution context to be set")
	}

	return v.PushBytes(v.Context.Origin[:])
}
