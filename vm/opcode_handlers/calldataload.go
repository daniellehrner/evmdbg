package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type CallDataLoadOpCode struct{}

func (*CallDataLoadOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(1); err != nil {
		return err
	}
	offset, err := v.Stack.Pop()
	if err != nil {
		return err
	}
	data := make([]byte, 32)
	start := offset.Uint64()

	if start < uint64(len(v.Context.CallData)) {
		end := start + 32
		if end > uint64(len(v.Context.CallData)) {
			end = uint64(len(v.Context.CallData))
		}
		copy(data, v.Context.CallData[start:end])
	}

	return v.Push(new(big.Int).SetBytes(data))
}
