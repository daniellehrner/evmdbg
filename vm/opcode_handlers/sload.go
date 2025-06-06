package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type SLoadOpCode struct{}

func (*SLoadOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(1); err != nil {
		return err
	}
	slot, err := v.Stack.Pop()
	if err != nil {
		return err
	}
	return v.Stack.Push(v.ReadStorage(slot))
}
