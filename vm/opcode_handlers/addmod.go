package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type AddModOpCode struct{}

func (*AddModOpCode) Execute(vm *vm.DebuggerVM) error {
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

	sum := new(big.Int).Add(a, b)
	sum.Mod(sum, m)
	return vm.Push(sum)
}
