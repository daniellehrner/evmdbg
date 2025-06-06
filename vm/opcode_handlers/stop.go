package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type StopOpCode struct{}

func (*StopOpCode) Execute(vm *vm.DebuggerVM) error {
	vm.Stopped = true
	return nil
}
