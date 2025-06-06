package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type DivOpCode struct{}

func (*DivOpCode) Execute(v *vm.DebuggerVM) error {
	// The DIV opcode divides the top two stack items.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// If the divisor is zero, push zero onto the stack.
	if b.Sign() == 0 {
		return v.Push(new(big.Int)) // EVM: push 0 on divide by zero
	}

	// Perform the division and push the result back onto the stack.
	return v.Push(new(big.Int).Div(a, b))
}
