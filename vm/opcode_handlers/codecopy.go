package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type CodeCopyOpCode struct{}

func (*CodeCopyOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(3); err != nil {
		return err
	}
	memOffset, codeOffset, length, err := v.Pop3()
	if err != nil {
		return err
	}

	start := codeOffset.Uint64()
	end := start + length.Uint64()

	data := make([]byte, length.Uint64())
	if start < uint64(len(v.Code)) {
		copyEnd := end
		if copyEnd > uint64(len(v.Code)) {
			copyEnd = uint64(len(v.Code))
		}
		copy(data, v.Code[start:copyEnd])
	}

	v.Memory.Write(int(memOffset.Uint64()), data)
	return nil
}
