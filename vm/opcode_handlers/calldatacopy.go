package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type CallDataCopyOpCode struct{}

func (*CallDataCopyOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(3); err != nil {
		return err
	}
	memOffset, dataOffset, length, err := v.Pop3()
	if err != nil {
		return err
	}

	start := dataOffset.Uint64()
	end := start + length.Uint64()

	var data []byte
	if start >= uint64(len(v.Context.CallData)) {
		data = make([]byte, length.Uint64())
	} else {
		if end > uint64(len(v.Context.CallData)) {
			end = uint64(len(v.Context.CallData))
		}
		data = make([]byte, length.Uint64())
		copy(data, v.Context.CallData[start:end])
	}

	v.Memory.Write(int(memOffset.Uint64()), data)
	return nil
}
