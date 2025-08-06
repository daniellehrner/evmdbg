package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type SDivOpCode struct{}

func (*SDivOpCode) Execute(v *vm.DebuggerVM) error {
	// The SDIV opcode requires two values on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// If the divisor is zero, return 0 as per EVM specification.
	if b.IsZero() {
		return v.Push(uint256.NewInt(0))
	}

	// Perform signed division using uint256's SDiv method
	result := new(uint256.Int).SDiv(a, b)
	return v.Push(result)
}
