package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type CallDataSizeOpCode struct{}

func (*CallDataSizeOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("call data size op code requires the execution context to be set")
	}

	size := uint64(len(v.Context.CallData))

	return v.Push(new(uint256.Int).SetUint64(size))
}
