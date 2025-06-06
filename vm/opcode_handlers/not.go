package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type NotOpCode struct{}

func (*NotOpCode) Execute(v *vm.DebuggerVM) error {
	// NOT requires one value on the stack.
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the top item from the stack.
	x, err := v.Stack.Pop()
	if err != nil {
		return err
	}

	// EVM word size is 256 bits: apply ^ over 32-byte mask
	// Perform the NOT operation by XORing with the mask.
	return v.Push(new(big.Int).Xor(x, uint256Mask))
}
