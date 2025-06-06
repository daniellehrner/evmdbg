package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type CallValueOpCode struct{}

func (*CallValueOpCode) Execute(v *vm.DebuggerVM) error {
	// If the call value is not set, return 0
	if v.Context.Value == nil {
		return v.Push(new(big.Int))
	}

	// Otherwise, push the call value onto the stack
	return v.Push(new(big.Int).Set(v.Context.Value))
}
