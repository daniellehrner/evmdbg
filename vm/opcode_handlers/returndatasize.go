package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type ReturnDataSizeOpCode struct{}

func (*ReturnDataSizeOpCode) Execute(v *vm.DebuggerVM) error {
	size := uint64(len(v.ReturnValue))
	return v.Push(new(uint256.Int).SetUint64(size))
}
