package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type GasPriceOpCode struct{}

func (*GasPriceOpCode) Execute(v *vm.DebuggerVM) error {
	if v.Context.GasPrice == nil {
		return v.Push(new(big.Int))
	}
	return v.Push(new(big.Int).Set(v.Context.GasPrice))
}
