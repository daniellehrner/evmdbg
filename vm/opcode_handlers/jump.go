package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
)

type JumpOpCode struct{}

func (*JumpOpCode) Execute(v *vm.DebuggerVM) error {
	// Jump opcode requires at least one item on the stack, which is the target PC.
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the target PC from the stack.
	target, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// The target must be a valid jump destination.
	pc := target.Uint64()
	if !v.IsJumpDest(pc) {
		return vm.ErrInvalidJump
	}

	v.SetPC(pc)
	return nil
}
