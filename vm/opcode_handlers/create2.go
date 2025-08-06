package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

type Create2OpCode struct{}

func (*Create2OpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("create2 op code requires the execution context to be set")
	}

	if v.StateProvider == nil {
		return fmt.Errorf("create2 op code requires state provider to be set")
	}

	// CREATE2 requires four values on the stack (value, offset, size, salt)
	if err := v.RequireStack(4); err != nil {
		return err
	}

	// Pop value, offset, size, and salt from stack
	value, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	offset, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	size, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	salt, err := v.Stack().Pop()
	if err != nil {
		return err
	}

	// Check for static call context - CREATE2 not allowed in static calls
	frame := v.CurrentFrame()
	if frame != nil && frame.IsStatic {
		return vm.ErrStaticCallStateChange
	}

	// Get initialization code from memory
	offsetUint64 := offset.Uint64()
	sizeUint64 := size.Uint64()

	initCode := v.Memory().Read(int(offsetUint64), int(sizeUint64))

	// Calculate new contract address using CREATE2 formula:
	// keccak256(0xff || sender || salt || keccak256(initCode))
	senderAddr := v.Context.Address

	// First hash the init code
	codeHasher := sha3.NewLegacyKeccak256()
	codeHasher.Write(initCode)
	codeHash := codeHasher.Sum(nil)

	// Then calculate the final address
	addrHasher := sha3.NewLegacyKeccak256()
	addrHasher.Write([]byte{0xff})  // 0xff prefix
	addrHasher.Write(senderAddr[:]) // sender address (20 bytes)

	// Salt as 32 bytes
	saltBytes := make([]byte, 32)
	salt.WriteToSlice(saltBytes)
	addrHasher.Write(saltBytes) // salt (32 bytes)

	addrHasher.Write(codeHash) // keccak256(initCode) (32 bytes)

	hashResult := addrHasher.Sum(nil)
	var newAddr [20]byte
	copy(newAddr[:], hashResult[12:32]) // Take last 20 bytes

	// Increment sender's nonce (CREATE2 also increments nonce)
	nonce := v.StateProvider.GetNonce(senderAddr)
	v.StateProvider.SetNonce(senderAddr, nonce+1)

	// Check if account already exists
	if v.StateProvider.AccountExists(newAddr) {
		// Push 0 to indicate failure
		return v.Push(uint256.NewInt(0))
	}

	// Create the new contract account
	err = v.StateProvider.CreateAccount(newAddr, initCode, value)
	if err != nil {
		// Push 0 to indicate failure
		return v.Push(uint256.NewInt(0))
	}

	// Mark account as created in current transaction (EIP-6780)
	v.MarkAccountCreatedInTransaction(newAddr)

	// Transfer value from caller to new contract
	if !value.IsZero() {
		callerBalance := v.StateProvider.GetBalance(senderAddr)
		if callerBalance.Cmp(value) < 0 {
			// Insufficient balance - push 0 to indicate failure
			return v.Push(uint256.NewInt(0))
		}

		// Note: In a full implementation, this would require extending StateProvider
		// to support setting balance. For now, we assume CreateAccount handles this.
	}

	// Push the new contract address onto the stack
	addrInt := new(uint256.Int).SetBytes(newAddr[:])
	return v.Push(addrInt)
}
