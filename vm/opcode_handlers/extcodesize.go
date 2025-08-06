package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type ExtCodeSizeOpCode struct{}

func (*ExtCodeSizeOpCode) Execute(v *vm.DebuggerVM) error {
	// EXTCODESIZE requires one value on the stack (the address)
	if err := v.RequireStack(1); err != nil {
		return err
	}

	// Pop the address from the stack
	addrInt, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// Convert to 20-byte address
	var addr [20]byte
	addrBytes := addrInt.Bytes()
	if len(addrBytes) > 20 {
		// Take only the last 20 bytes if longer
		copy(addr[:], addrBytes[len(addrBytes)-20:])
	} else {
		// Right-align if shorter
		copy(addr[20-len(addrBytes):], addrBytes)
	}

	var codeSize *uint256.Int
	if v.StateProvider != nil {
		// Get code from state provider
		code := v.StateProvider.GetCode(addr)
		codeSize = uint256.NewInt(uint64(len(code)))
	} else {
		// If no state provider, return 0 (account doesn't exist)
		codeSize = uint256.NewInt(0)
	}

	// Push the code size onto the stack
	return v.Push(codeSize)
}
