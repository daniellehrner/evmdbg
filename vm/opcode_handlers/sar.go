package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SarOpCode struct{}

func (*SarOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	shift, value, err := v.Pop2()
	if err != nil {
		return err
	}
	if shift.BitLen() > 256 {
		// special case: value < 0 → -1, else → 0
		if value.Bit(255) == 1 {
			return v.Push(new(big.Int).Sub(new(big.Int), big.NewInt(1))) // -1
		}
		return v.Push(new(big.Int)) // 0
	}
	n := shift.Uint64()
	if n >= 256 {
		if value.Bit(255) == 1 {
			return v.Push(new(big.Int).Sub(new(big.Int), big.NewInt(1)))
		}
		return v.Push(new(big.Int))
	}
	// signed interpretation
	tmp := new(big.Int).Set(value)
	if tmp.Bit(255) == 1 {
		tmp.Sub(tmp, new(big.Int).Lsh(big.NewInt(1), 256)) // convert to negative
	}
	tmp.Rsh(tmp, uint(n))
	// wrap back into 256-bit space
	if tmp.Sign() < 0 {
		tmp.Add(tmp, new(big.Int).Lsh(big.NewInt(1), 256))
	}
	return v.Push(tmp)
}
