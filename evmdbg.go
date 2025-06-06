package main

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/daniellehrner/evmdbg/vm/opcode_handlers"
)

func NewDebuggerVM(code []byte) *vm.DebuggerVM {
	d := vm.NewDebuggerVM(code, opcode_handlers.GetHandler)
	return d
}
