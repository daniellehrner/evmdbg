package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SLTOpCode struct{}

func (*SLTOpCode) Execute(v *vm.DebuggerVM) error {
	// The SLT opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// if the signed value of a is less than the signed value of b
	if toSigned(a, uint256Half).Cmp(toSigned(b, uint256Half)) < 0 {
		return v.Push(big.NewInt(1))
	}

	// If a is not less than b, push 0 onto the stack
	return v.Push(big.NewInt(0))
}
