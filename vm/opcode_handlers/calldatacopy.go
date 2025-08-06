package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type CallDataCopyOpCode struct{}

func (*CallDataCopyOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("call data copy op code requires the execution context to be set")
	}

	// The CALLDATACOPY opcode requires at least three items on the stack:
	if err := v.RequireStack(3); err != nil {
		return err
	}

	// Pop the top three items from the stack.
	memOffset, dataOffset, length, err := v.Pop3()
	if err != nil {
		return err
	}

	start := dataOffset.Uint64()
	end := start + length.Uint64()
	var data []byte

	// If the start offset is beyond the length of call data, we write zeroes.
	if start >= uint64(len(v.Context.CallData)) {
		data = make([]byte, length.Uint64())
	} else {
		// If the end offset exceeds the length of call data, we adjust it.
		if end > uint64(len(v.Context.CallData)) {
			end = uint64(len(v.Context.CallData))
		}
		// Copy the relevant slice of call data.
		data = make([]byte, length.Uint64())
		copy(data, v.Context.CallData[start:end])
	}

	// Write the data to memory at the specified memory offset.
	v.Memory.Write(int(memOffset.Uint64()), data)
	return nil
}
