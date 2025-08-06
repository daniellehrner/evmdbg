package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
)

type JumpiOpCode struct{}

func (*JumpiOpCode) Execute(v *vm.DebuggerVM) error {
	// JUMPI requires two values on the stack: a target and a condition.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the target and condition from the stack.
	target, cond, err := v.Pop2()
	if err != nil {
		return err
	}

	// Only jump if the condition is non-zero.
	if cond.Sign() != 0 {
		pc := target.Uint64()

		// The target must be a valid jump destination.
		if !v.IsJumpDest(pc) {
			return fmt.Errorf("invalid jumpi target: 0x%x", pc)
		}
		v.SetPC(pc)
	}
	return nil
}
