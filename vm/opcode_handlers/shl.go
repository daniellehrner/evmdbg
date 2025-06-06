package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ShlOpCode struct{}

func (*ShlOpCode) Execute(vm *vm.DebuggerVM) error {
	if err := vm.RequireStack(2); err != nil {
		return err
	}

	shift, value, err := vm.Pop2()
	if err != nil {
		return err
	}

	if shift.BitLen() > 256 {
		return vm.Push(new(big.Int)) // result is zero
	}

	n := shift.Uint64()
	if n >= 256 {
		return vm.Push(new(big.Int)) // shift too far, zero
	}

	res := new(big.Int).Lsh(value, uint(n))
	res.And(res, uint256Mask) // mask to 256 bits

	return vm.Push(res)
}
