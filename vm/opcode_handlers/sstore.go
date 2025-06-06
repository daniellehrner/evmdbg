package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type SStoreOpCode struct{}

func (*SStoreOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	slot, val, err := v.Pop2()
	if err != nil {
		return err
	}
	v.WriteStorage(slot, val)
	return nil
}
