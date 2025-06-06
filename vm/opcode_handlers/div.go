package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type DivOpCode struct{}

func (*DivOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}
	if b.Sign() == 0 {
		return v.Push(new(big.Int)) // EVM: push 0 on divide by zero
	}
	return v.Push(new(big.Int).Div(a, b))
}
