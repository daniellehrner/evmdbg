package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ExpOpCode struct{}

func (*ExpOpCode) Execute(vm *vm.DebuggerVM) error {
	if err := vm.RequireStack(2); err != nil {
		return err
	}
	base, exponent, err := vm.Pop2()
	if err != nil {
		return err
	}
	result := new(big.Int).Exp(base, exponent, uint256)
	return vm.Push(result)
}
