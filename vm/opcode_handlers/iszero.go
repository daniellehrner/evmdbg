package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type IsZeroOpCode struct{}

func (*IsZeroOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(1); err != nil {
		return err
	}
	x, err := v.Stack.Pop()
	if err != nil {
		return err
	}
	if x.Sign() == 0 {
		return v.PushUint64(1)
	}
	return v.PushUint64(0)
}
