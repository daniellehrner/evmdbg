package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type OrOpCode struct{}

func (*OrOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}
	return v.Push(new(big.Int).Or(a, b))
}
