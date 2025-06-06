package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type DupOpCode struct {
	N int // 1 to 16
}

func (d *DupOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(d.N); err != nil {
		return err
	}
	val, err := v.Stack.Peek(d.N - 1)
	if err != nil {
		return err
	}
	return v.Stack.Push(new(big.Int).Set(val)) // copy
}
