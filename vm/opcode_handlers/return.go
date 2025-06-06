package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type ReturnOpCode struct{}

func (*ReturnOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	offset, size, err := v.Pop2()
	if err != nil {
		return err
	}
	v.ReturnValue = v.Memory.Read(int(offset.Uint64()), int(size.Uint64()))
	v.Stopped = true
	return nil
}
