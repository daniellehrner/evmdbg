package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

type BalanceOpCode struct{}

func (*BalanceOpCode) Execute(v *vm.DebuggerVM) error {
	// BALANCE requires one value on the stack (the address)
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

	var balance *uint256.Int
	if v.StateProvider != nil {
		// Get balance from state provider
		balance = v.StateProvider.GetBalance(addr)
	} else {
		// If no state provider, return 0
		balance = uint256.NewInt(0)
	}

	// Push the balance onto the stack
	return v.Push(balance)
}
