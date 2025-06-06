package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type SwapOpCode struct {
	N int // 1 to 16
}

func (s *SwapOpCode) Execute(v *vm.DebuggerVM) error {
	// The SWAP opcode requires N+1 values on the stack
	if err := v.RequireStack(s.N + 1); err != nil {
		return err
	}

	// Swap item at position N with the top of the stack
	return v.Stack.Swap(s.N)
}
