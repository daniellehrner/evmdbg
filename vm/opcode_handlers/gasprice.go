package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type GasPriceOpCode struct{}

func (*GasPriceOpCode) Execute(v *vm.DebuggerVM) error {
	// If the gas price is not set, return 0
	if v.Context.GasPrice == nil {
		return v.Push(new(big.Int))
	}

	// Otherwise, push the gas price onto the stack
	return v.Push(new(big.Int).Set(v.Context.GasPrice))
}
