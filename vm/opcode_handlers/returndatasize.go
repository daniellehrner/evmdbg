package opcode_handlers

import (
	"math/big"

	"github.com/daniellehrner/evmdbg/vm"
)

type ReturnDataSizeOpCode struct{}

func (*ReturnDataSizeOpCode) Execute(v *vm.DebuggerVM) error {
	size := uint64(len(v.ReturnValue))
	return v.Push(new(big.Int).SetUint64(size))
}
