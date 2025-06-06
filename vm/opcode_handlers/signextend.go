package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SignExtendOpCode struct{}

func (*SignExtendOpCode) Execute(v *vm.DebuggerVM) error {
	// The SIGNEXTEND opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	k, val, err := v.Pop2()
	if err != nil {
		return err
	}

	// If k is greater than or equal to 31, we push the value as is
	if k.Cmp(big.NewInt(31)) >= 0 {
		return v.Push(val)
	}

	kInt := int(k.Int64())

	// The bit to check is (k + 1) * 8 - 1
	bitIndex := (kInt+1)*8 - 1

	// If the sign bit is 1 (negative number in two's complement):
	if val.Bit(bitIndex) == 1 {
		// Create a mask with 1 << bitIndex (i.e., the sign bit position)
		mask := new(big.Int).Lsh(big.NewInt(1), uint(bitIndex))

		// Subtract 1 to get a bitmask with all bits below the sign bit set (e.g., 0x7F for bitIndex=7)
		mask.Sub(mask, big.NewInt(1))

		// Invert the mask so that all bits above the sign bit are 1 (e.g., 0xFFFFFF80)
		// OR the value with the inverted mask to extend the sign bit across higher bits
		val.Or(val, new(big.Int).Not(mask))
	} else {
		// Same: mask = 1 << bitIndex
		mask := new(big.Int).Lsh(big.NewInt(1), uint(bitIndex))

		// mask = (1 << bitIndex) - 1, so all bits below bitIndex are set
		mask.Sub(mask, big.NewInt(1))

		// AND the value with the mask to clear bits above the sign bit
		val.And(val, mask)
	}

	return v.Push(val)
}
