package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type JumpiOpCode struct{}

func (*JumpiOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	target, cond, err := v.Pop2()
	if err != nil {
		return err
	}
	if cond.Sign() != 0 {
		pc := target.Uint64()
		if _, ok := v.JumpDests[pc]; !ok {
			return fmt.Errorf("invalid jumpi target: 0x%x", pc)
		}
		v.PC = pc
	}
	return nil
}
