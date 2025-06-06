package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type EqOpCode struct{}

func (*EqOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}
	if a.Cmp(b) == 0 {
		return v.PushUint64(1)
	}
	return v.PushUint64(0)
}
