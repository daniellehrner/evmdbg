package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

type ExtCodeHashOpCode struct{}

func (*ExtCodeHashOpCode) Execute(v *vm.DebuggerVM) error {
	// EXTCODEHASH requires one value on the stack (the address)
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

	var codeHash *uint256.Int

	if v.StateProvider != nil {
		// Get code from state provider
		code := v.StateProvider.GetCode(addr)

		if len(code) == 0 && !v.StateProvider.AccountExists(addr) {
			// Non-existent account returns 0
			codeHash = uint256.NewInt(0)
		} else {
			// Compute Keccak-256 hash of the code
			hasher := sha3.NewLegacyKeccak256()
			hasher.Write(code)
			hashBytes := hasher.Sum(nil)
			codeHash = new(uint256.Int).SetBytes(hashBytes)
		}
	} else {
		// If no state provider, return 0 (account doesn't exist)
		codeHash = uint256.NewInt(0)
	}

	// Push the code hash onto the stack
	return v.Push(codeHash)
}
