package opcode_handlers

import (
	"math/big"

	"github.com/daniellehrner/evmdbg/vm"
)

type Push0OpCode struct{}

func (*Push0OpCode) Execute(vm *vm.DebuggerVM) error {
	return vm.Push(big.NewInt(0))
}
