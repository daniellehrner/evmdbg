package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type ShrOpCode struct{}

func (*ShrOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	shift, value, err := v.Pop2()
	if err != nil {
		return err
	}
	if shift.BitLen() > 256 {
		return v.Push(new(big.Int)) // result is zero
	}
	n := shift.Uint64()
	if n >= 256 {
		return v.Push(new(big.Int)) // shift too far, zero
	}
	res := new(big.Int).Rsh(value, uint(n))
	return v.Push(res)
}
