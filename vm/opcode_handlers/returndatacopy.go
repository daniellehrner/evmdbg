package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type ReturnDataCopyOpCode struct{}

func (*ReturnDataCopyOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(3); err != nil {
		return err
	}
	memOffset, offset, size, err := v.Pop3()
	if err != nil {
		return err
	}

	start := offset.Uint64()
	end := start + size.Uint64()
	if end > uint64(len(v.ReturnValue)) {
		return fmt.Errorf("RETURNDATACOPY out of bounds: %d > %d", end, len(v.ReturnValue))
	}

	data := v.ReturnValue[start:end]
	v.Memory.Write(int(memOffset.Uint64()), data)
	return nil
}
