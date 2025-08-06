package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type GasOpCode struct{}

func (*GasOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("gas op code requires the execution context to be set")
	}

	return v.PushUint64(v.Context.Gas)
}
