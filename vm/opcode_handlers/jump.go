package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type JumpOpCode struct{}

func (*JumpOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(1); err != nil {
		return err
	}
	target, err := v.Stack.Pop()
	if err != nil {
		return err
	}
	pc := target.Uint64()
	if _, ok := v.JumpDests[pc]; !ok {
		return fmt.Errorf("invalid jump target: 0x%x", pc)
	}
	v.PC = pc
	return nil
}
