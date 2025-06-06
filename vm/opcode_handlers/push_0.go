package opcode_handlers

import (
	"math/big"

	"github.com/daniellehrner/evmdbg/vm"
)

type Push0OpCode struct{}

func (*Push0OpCode) Execute(v *vm.DebuggerVM) error {
	// Push the value 0 onto the stack.
	return v.Push(big.NewInt(0))
}
