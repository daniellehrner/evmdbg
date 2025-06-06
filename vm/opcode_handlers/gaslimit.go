package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type GasLimitOpCode struct{}

func (*GasLimitOpCode) Execute(v *vm.DebuggerVM) error {
	if v.Context.Block == nil {
		return v.Push(new(big.Int))
	}

	return v.PushUint64(v.Context.Block.GasLimit)
}
