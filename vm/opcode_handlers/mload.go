package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type MLoadOpCode struct{}

func (*MLoadOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(1); err != nil {
		return err
	}
	addr, err := v.Stack.Pop()
	if err != nil {
		return err
	}
	return v.Stack.Push(v.Memory.ReadWord(addr.Uint64()))
}
