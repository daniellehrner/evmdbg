package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
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

	// Get shift amount as uint64, handling overflow
	shiftAmount, overflow := shift.Uint64WithOverflow()
	if overflow || shiftAmount >= 256 {
		// If shift >= 256, result depends on sign of value
		if value.Sign() < 0 { // Check if MSB is set (negative in two's complement)
			// Return -1 (all bits set)
			result := new(uint256.Int).SetAllOne()
			return v.Push(result)
		} else {
			// Return 0
			result := new(uint256.Int).Clear()
			return v.Push(result)
		}
	}

	// Perform signed arithmetic right shift using uint256's SRsh method
	result := new(uint256.Int).SRsh(value, uint(shiftAmount))
	return v.Push(result)
}
