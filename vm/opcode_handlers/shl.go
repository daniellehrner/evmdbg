package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type ShlOpCode struct{}

func (*ShlOpCode) Execute(v *vm.DebuggerVM) error {
	// The SHL opcode requires two values on the stack
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
		// Handle edge case for large shifts
		// For SHL: always return 0 when shift >= 256
		return v.Push(new(uint256.Int).Clear())
	}

	// Convert to uint for use with Lsh
	result := new(uint256.Int).Lsh(value, uint(shiftAmount))
	return v.Push(result)
}
