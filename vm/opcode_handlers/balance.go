package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type BalanceOpCode struct{}

func (*BalanceOpCode) Execute(vm *vm.DebuggerVM) error {
	if err := vm.RequireStack(1); err != nil {
		return err
	}
	addrBytes, err := vm.Stack.Pop()
	if err != nil {
		return err
	}
	// The debugger does not simulate other accounts, so we assume balance is 0 for now
	// You could later plug in external state if needed
	if addrBytes.BitLen() > 160 {
		return fmt.Errorf("invalid address length")
	}
	return vm.Push(vm.Context.Balance)
}
