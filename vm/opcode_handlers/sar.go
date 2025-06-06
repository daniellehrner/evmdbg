package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SarOpCode struct{}

func (*SarOpCode) Execute(v *vm.DebuggerVM) error {
	// SAR requires two values on the stack: shift and value.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	shift, value, err := v.Pop2()
	if err != nil {
		return err
	}

	// EVM word size is 256 bits: apply arithmetic right shift
	if shift.BitLen() > 256 {
		// special case: value < 0 → -1, else → 0
		if value.Bit(255) == 1 {
			return v.Push(new(big.Int).Sub(new(big.Int), big.NewInt(1))) // -1
		}
		// value >= 0 → 0
		return v.Push(new(big.Int))
	}

	// n is the number of bits to shift
	n := shift.Uint64()

	// If n is greater than or equal to 256, we can directly return the value based on its sign.
	if n >= 256 {
		if value.Bit(255) == 1 {
			return v.Push(new(big.Int).Sub(new(big.Int), big.NewInt(1)))
		}

		// If the value is non-negative, we return 0.
		return v.Push(new(big.Int))
	}

	// signed interpretation
	tmp := new(big.Int).Set(value)
	if tmp.Bit(255) == 1 {
		tmp.Sub(tmp, new(big.Int).Lsh(big.NewInt(1), 256)) // convert to negative
	}

	// perform arithmetic right shift
	tmp.Rsh(tmp, uint(n))
	// wrap back into 256-bit space
	if tmp.Sign() < 0 {
		tmp.Add(tmp, new(big.Int).Lsh(big.NewInt(1), 256))
	}

	return v.Push(tmp)
}
