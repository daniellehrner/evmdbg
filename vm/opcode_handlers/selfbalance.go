package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type SelfBalanceOpCode struct{}

func (*SelfBalanceOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("selfbalance op code requires the execution context to be set")
	}

	// Push the current account's balance onto the stack
	return v.Push(v.Context.Balance)
}
