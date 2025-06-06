package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type JumpDestOpCode struct{}

func (*JumpDestOpCode) Execute(_ *vm.DebuggerVM) error {
	// No-op. Valid jump target.
	return nil
}
