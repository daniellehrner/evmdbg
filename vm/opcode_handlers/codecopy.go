package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type CodeCopyOpCode struct{}

func (*CodeCopyOpCode) Execute(v *vm.DebuggerVM) error {
	// The CODECOPY opcode requires at least three items on the stack.
	if err := v.RequireStack(3); err != nil {
		return err
	}

	// Pop the top three items from the stack.
	memOffset, codeOffset, length, err := v.Pop3()
	if err != nil {
		return err
	}

	start := codeOffset.Uint64()
	end := start + length.Uint64()
	data := make([]byte, length.Uint64())

	// If the start offset is beyond the length of the code, we write zeroes.
	if start < uint64(len(v.Code)) {
		copyEnd := end

		// If the end offset exceeds the length of code, we adjust it.
		if copyEnd > uint64(len(v.Code)) {
			copyEnd = uint64(len(v.Code))
		}
		copy(data, v.Code[start:copyEnd])
	}

	v.Memory.Write(int(memOffset.Uint64()), data)
	return nil
}
