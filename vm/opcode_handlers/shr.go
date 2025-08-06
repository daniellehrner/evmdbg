package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type ShrOpCode struct{}

func (*ShrOpCode) Execute(v *vm.DebuggerVM) error {
	// The SHR opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	shift, value, err := v.Pop2()
	if err != nil {
		return err
	}

	// Get shift amount as uint64, handling overflow
	shiftAmount, overflow := shift.Uint64WithOverflow()
	if overflow || shiftAmount >= 256 {
		// If shift >= 256, result is always 0 for logical right shift
		return v.Push(new(uint256.Int).Clear())
	}

	// Perform logical right shift using uint256's Rsh method
	result := new(uint256.Int).Rsh(value, uint(shiftAmount))
	return v.Push(result)
}
