package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type InvalidOpCode struct{}

func (*InvalidOpCode) Execute(v *vm.DebuggerVM) error {
	// INVALID opcode always causes execution to halt with an error
	return fmt.Errorf("invalid opcode 0xFE encountered")
}
