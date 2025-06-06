package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type SwapOpCode struct {
	N int // 1 to 16
}

func (s *SwapOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(s.N + 1); err != nil {
		return err
	}
	return v.Stack.Swap(s.N)
}
