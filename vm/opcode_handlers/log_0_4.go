package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
)

type LogNOpCode struct {
	N int // number of topics: 0..4
}

func (op *LogNOpCode) Execute(v *vm.DebuggerVM) error {
	// LogN requires 2 + N values on the stack: offset, size, and N topics.
	if err := v.RequireStack(2 + op.N); err != nil {
		return err
	}

	// Pop in reverse order: topics, then size, then offset
	topics := make([][]byte, op.N)
	for i := op.N - 1; i >= 0; i-- {
		t, err := v.Stack.Pop()
		if err != nil {
			return err
		}
		topics[i] = padTo256Bytes(t.Bytes())
	}

	size, err := v.Stack.Pop()
	if err != nil {
		return err
	}

	offset, err := v.Stack.Pop()
	if err != nil {
		return err
	}

	// read data from memory at the specified offset and size
	data := v.Memory.Read(int(offset.Uint64()), int(size.Uint64()))

	v.Logs = append(v.Logs, vm.LogEntry{
		Address: v.Context.Address,
		Topics:  topics,
		Data:    data,
	})

	return nil
}
