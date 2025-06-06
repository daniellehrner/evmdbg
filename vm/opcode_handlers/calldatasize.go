package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type CallDataSizeOpCode struct{}

func (*CallDataSizeOpCode) Execute(v *vm.DebuggerVM) error {
	size := uint64(len(v.Context.CallData))
	return v.Push(new(big.Int).SetUint64(size))
}
