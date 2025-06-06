package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type MulModOpCode struct{}

func (*MulModOpCode) Execute(vm *vm.DebuggerVM) error {
	if err := vm.RequireStack(3); err != nil {
		return err
	}
	a, b, m, err := vm.Pop3()
	if err != nil {
		return err
	}

	if m.Sign() == 0 {
		return vm.Push(big.NewInt(0))
	}

	prod := new(big.Int).Mul(a, b)
	prod.Mod(prod, m)
	return vm.Push(prod)
}
