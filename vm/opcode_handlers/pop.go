package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type PopOpCode struct{}

func (*PopOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(1); err != nil {
		return err
	}
	_, err := v.Stack.Pop()
	return err
}
