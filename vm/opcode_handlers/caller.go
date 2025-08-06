package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type CallerOpCode struct{}

func (*CallerOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("caller op code requires the execution context to be set")
	}

	return v.PushBytes(v.Context.Caller[:])
}
