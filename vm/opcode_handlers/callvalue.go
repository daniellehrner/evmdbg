package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type CallValueOpCode struct{}

func (*CallValueOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("call value code requires the execution context to be set")
	}

	// If the call value is not set, return 0
	if v.Context.Value == nil {
		return v.Push(new(uint256.Int))
	}

	// Otherwise, push the call value onto the stack
	return v.Push(new(uint256.Int).Set(v.Context.Value))
}
