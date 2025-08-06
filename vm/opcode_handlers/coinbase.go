package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type CoinbaseOpCode struct{}

func (*CoinbaseOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("coinbase op code requires the execution context to be set")
	}

	if v.Context.Block == nil {
		return v.Push(new(uint256.Int))
	}

	return v.PushBytes(v.Context.Block.Coinbase[:])
}
