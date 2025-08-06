package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type GasPriceOpCode struct{}

func (*GasPriceOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("gas price op code requires the execution context to be set")
	}

	// If the gas price is not set, return 0
	if v.Context.GasPrice == nil {
		return v.Push(new(uint256.Int))
	}

	// Otherwise, push the gas price onto the stack
	return v.Push(new(uint256.Int).Set(v.Context.GasPrice))
}
