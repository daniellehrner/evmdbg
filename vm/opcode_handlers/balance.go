package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type BalanceOpCode struct{}

func (*BalanceOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("balance op code requires the execution context to be set")
	}

	// The BALANCE opcode requires at least one item on the stack.
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the address from the stack.
	addrBytes, err := v.Stack.Pop()
	if err != nil {
		return err
	}

	// The address must be 20 bytes (160 bits) long.
	if addrBytes.BitLen() > 160 {
		return fmt.Errorf("invalid address length")
	}

	return v.Push(v.Context.Balance)
}
