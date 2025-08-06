package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type TimestampOpCode struct{}

func (*TimestampOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("timestamp op code requires the execution context to be set")
	}

	if v.Context.Block == nil {
		return v.Push(new(uint256.Int))
	}

	return v.PushUint64(v.Context.Block.Timestamp)
}
