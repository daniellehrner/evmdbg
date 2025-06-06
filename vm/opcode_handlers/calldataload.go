package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type CallDataLoadOpCode struct{}

func (*CallDataLoadOpCode) Execute(v *vm.DebuggerVM) error {
	// The CALLDATALOAD opcode requires at least one item on the stack.
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the offset from the stack.
	offset, err := v.Stack.Pop()
	if err != nil {
		return err
	}

	data := make([]byte, 32)
	start := offset.Uint64()

	// If the start offset is beyond the length of call data, we write zeroes.
	if start < uint64(len(v.Context.CallData)) {
		end := start + 32

		// If the end offset exceeds the length of call data, we adjust it.
		if end > uint64(len(v.Context.CallData)) {
			end = uint64(len(v.Context.CallData))
		}

		// Copy the relevant slice of call data.
		copy(data, v.Context.CallData[start:end])
	}

	return v.Push(new(big.Int).SetBytes(data))
}
