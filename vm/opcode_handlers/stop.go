package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type StopOpCode struct{}

func (*StopOpCode) Execute(vm *vm.DebuggerVM) error {
	// Set the Stopped flag to true to indicate that the VM has stopped executing.
	vm.Stopped = true

	return nil
}
