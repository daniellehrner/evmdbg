package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type MulOpCode struct{}

func (*MulOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}
	return v.Push(new(big.Int).Mul(a, b))
}
