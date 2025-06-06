package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SignExtendOpCode struct{}

func (*SignExtendOpCode) Execute(vm *vm.DebuggerVM) error {
	if err := vm.RequireStack(2); err != nil {
		return err
	}
	k, val, err := vm.Pop2()
	if err != nil {
		return err
	}

	if k.Cmp(big.NewInt(31)) >= 0 {
		return vm.Push(val)
	}

	kInt := int(k.Int64())
	bitIndex := (kInt+1)*8 - 1
	if val.Bit(bitIndex) == 1 {
		mask := new(big.Int).Lsh(big.NewInt(1), uint(bitIndex))
		mask.Sub(mask, big.NewInt(1))
		val.Or(val, new(big.Int).Not(mask))
	} else {
		mask := new(big.Int).Lsh(big.NewInt(1), uint(bitIndex))
		mask.Sub(mask, big.NewInt(1))
		val.And(val, mask)
	}

	return vm.Push(val)
}
