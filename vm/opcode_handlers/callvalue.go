package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type CallValueOpCode struct{}

func (*CallValueOpCode) Execute(v *vm.DebuggerVM) error {
	if v.Context.Value == nil {
		return v.Push(new(big.Int))
	}
	return v.Push(new(big.Int).Set(v.Context.Value))
}
