package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type Push0OpCode struct{}

func (*Push0OpCode) Execute(v *vm.DebuggerVM) error {
	// Push the value 0 onto the stack.
	return v.Push(uint256.NewInt(0))
}
