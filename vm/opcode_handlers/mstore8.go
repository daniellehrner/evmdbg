package opcode_handlers

import "github.com/daniellehrner/evmdbg/vm"

type MStore8OpCode struct{}

func (*MStore8OpCode) Execute(v *vm.DebuggerVM) error {
	// MSTORE8 requires two values on the stack.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	addr, val, err := v.Pop2()
	if err != nil {
		return err
	}

	// Write only the least significant byte to memory at the given address.
	v.Memory().Write(int(addr.Uint64()), []byte{byte(val.Uint64())})

	return nil
}
