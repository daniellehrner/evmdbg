package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
)

type ReturnDataSizeOpCode struct{}

func (*ReturnDataSizeOpCode) Execute(v *vm.DebuggerVM) error {
	// RETURNDATASIZE pushes the size of return data from the last call onto the stack
	returnDataSize := v.ReturnDataSize()
	return v.Push(returnDataSize)
}
