package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ModOpCode struct{}

func (*ModOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}
	if b.Sign() == 0 {
		return v.Push(new(big.Int)) // EVM: mod by zero = 0
	}
	return v.Push(new(big.Int).Mod(a, b))
}
