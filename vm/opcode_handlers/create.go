package opcode_handlers

import (
	"fmt"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

type CreateOpCode struct{}

func (*CreateOpCode) Execute(v *vm.DebuggerVM) error {
	err := v.RequireContext()
	if err != nil {
		return fmt.Errorf("create op code requires the execution context to be set")
	}

	if v.StateProvider == nil {
		return fmt.Errorf("create op code requires state provider to be set")
	}

	// CREATE requires three values on the stack (value, offset, size)
	if err := v.RequireStack(3); err != nil {
		return err
	}

	// Pop value, offset, and size from stack
	value, offset, size, err := v.Pop3()
	if err != nil {
		return err
	}

	// Check for static call context - CREATE not allowed in static calls
	frame := v.CurrentFrame()
	if frame != nil && frame.IsStatic {
		return vm.ErrStaticCallStateChange
	}

	// Get initialization code from memory
	offsetUint64 := offset.Uint64()
	sizeUint64 := size.Uint64()

	initCode := v.Memory().Read(int(offsetUint64), int(sizeUint64))

	// Calculate new contract address using CREATE formula: keccak256(rlp([sender, nonce]))
	senderAddr := v.Context.Address
	nonce := v.StateProvider.GetNonce(senderAddr)

	// Simple RLP encoding for [address, nonce]
	rlpEncoded := encodeRLPCreateData(senderAddr, nonce)

	// Calculate keccak256 hash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(rlpEncoded)
	hashResult := hasher.Sum(nil)

	var newAddr [20]byte
	copy(newAddr[:], hashResult[12:32]) // Take last 20 bytes

	// Increment sender's nonce
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

// encodeRLPCreateData encodes [sender, nonce] for CREATE address calculation
// RLP encoding rules:
// - Single bytes < 0x80: encoded as themselves
// - Byte arrays 0-55 bytes: 0x80 + length, then data
// - Integers: minimal big-endian representation (no leading zeros)
// - Lists 0-55 total bytes: 0xc0 + length, then concatenated encoded elements
func encodeRLPCreateData(sender [20]byte, nonce uint64) []byte {
	// Encode sender address (20 bytes)
	senderEncoded := make([]byte, 21) // 0x80 + 20 + 20 bytes
	senderEncoded[0] = 0x80 + 20      // 20-byte string prefix
	copy(senderEncoded[1:], sender[:])

	// Encode nonce (minimal representation)
	var nonceEncoded []byte
	if nonce == 0 {
		nonceEncoded = []byte{0x80} // empty byte string for zero
	} else if nonce < 0x80 {
		nonceEncoded = []byte{byte(nonce)} // single byte
	} else {
		// Multi-byte integer - encode as minimal big-endian
		nonceBytes := []byte{}
		temp := nonce
		for temp > 0 {
			nonceBytes = append([]byte{byte(temp)}, nonceBytes...)
			temp >>= 8
		}
		nonceEncoded = make([]byte, 1+len(nonceBytes))
		nonceEncoded[0] = 0x80 + byte(len(nonceBytes))
		copy(nonceEncoded[1:], nonceBytes)
	}

	// Encode as list: 0xc0 + total_length + sender_encoded + nonce_encoded
	totalContentLength := len(senderEncoded) + len(nonceEncoded)
	result := make([]byte, 1+totalContentLength)
	result[0] = 0xc0 + byte(totalContentLength) // List prefix
	copy(result[1:], senderEncoded)
	copy(result[1+len(senderEncoded):], nonceEncoded)

	return result
}
