package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type RevertOpCode struct{}

func (*RevertOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	offset, size, err := v.Pop2()
	if err != nil {
		return err
	}
	v.ReturnValue = v.Memory.Read(int(offset.Uint64()), int(size.Uint64()))
	v.Reverted = true
	v.Stopped = true
	return nil
}
