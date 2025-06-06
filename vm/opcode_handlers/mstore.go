package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type MStoreOpCode struct{}

func (*MStoreOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	addr, val, err := v.Pop2()
	if err != nil {
		return err
	}
	v.Memory.WriteWord(addr.Uint64(), val)
	return nil
}
