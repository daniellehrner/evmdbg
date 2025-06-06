package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type XorOpCode struct{}

func (*XorOpCode) Execute(v *vm.DebuggerVM) error {
	// The XOR opcode requires two values on the stack
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}

	// Perform the bitwise XOR operation
	return v.Push(new(big.Int).Xor(a, b))
}
