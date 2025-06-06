package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SGTOpCode struct{}

func (*SGTOpCode) Execute(v *vm.DebuggerVM) error {
	// The SGT opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Define the half of the uint256 value for signed comparison
	if toSigned(a, uint256Half).Cmp(toSigned(b, uint256Half)) > 0 {
		return v.Push(big.NewInt(1))
	}

	// If a is not greater than b, push 0 onto the stack
	return v.Push(big.NewInt(0))
}
