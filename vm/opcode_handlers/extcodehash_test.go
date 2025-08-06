package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/sha3"
)

func TestExtCodeHashOpCode_Execute(t *testing.T) {
	// Test address: 0x1234567890123456789012345678901234567890
	testAddr := [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
	testCode := []byte{0x60, 0x01, 0x60, 0x02, 0x01} // Simple bytecode

	code := []byte{
		0x73,                                                       // PUSH20
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, // address bytes
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90,
		0x3f, // EXTCODEHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider
	mock := &mockStateProvider{
		codeMap: map[[20]byte][]byte{
			testAddr: testCode,
		},
	}
	v.StateProvider = mock

	// Compute expected hash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(testCode)
	expectedHashBytes := hasher.Sum(nil)
	expectedHash := new(uint256.Int).SetBytes(expectedHashBytes)

	// Execute PUSH20
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH20: %v", err)
	}

	// Execute EXTCODEHASH
	err = v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during EXTCODEHASH: %v", err)
	}

	// Check result
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if hash.Cmp(expectedHash) != 0 {
		t.Errorf("Expected code hash %s, got %s", expectedHash.Hex(), hash.Hex())
	}
}

func TestExtCodeHashOpCode_EmptyCode(t *testing.T) {
	// Test with account that has empty code
	testAddr := [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
	testCode := []byte{} // Empty code

	code := []byte{
		0x73,                                                       // PUSH20
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, // address bytes
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90,
		0x3f, // EXTCODEHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider - account exists but has empty code
	mock := &mockStateProvider{
		codeMap: map[[20]byte][]byte{
			testAddr: testCode,
		},
	}
	v.StateProvider = mock

	// Compute expected hash of empty code
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(testCode) // Empty code
	expectedHashBytes := hasher.Sum(nil)
	expectedHash := new(uint256.Int).SetBytes(expectedHashBytes)

	// Execute PUSH20 and EXTCODEHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if hash.Cmp(expectedHash) != 0 {
		t.Errorf("Expected empty code hash %s, got %s", expectedHash.Hex(), hash.Hex())
	}
}

func TestExtCodeHashOpCode_NonExistentAccount(t *testing.T) {
	code := []byte{
		0x73,                                                       // PUSH20
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // non-existent address
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		0x3f, // EXTCODEHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider with no accounts
	mock := &mockStateProvider{
		codeMap: map[[20]byte][]byte{},
	}
	v.StateProvider = mock

	// Execute PUSH20 and EXTCODEHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should be 0 for non-existent account
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !hash.IsZero() {
		t.Errorf("Expected code hash 0 for non-existent account, got %s", hash.Hex())
	}
}

func TestExtCodeHashOpCode_NoStateProvider(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01 (dummy address)
		0x3f, // EXTCODEHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set StateProvider (remains nil)

	// Execute PUSH1 and EXTCODEHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should be 0 when no state provider
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !hash.IsZero() {
		t.Errorf("Expected code hash 0 with no state provider, got %s", hash.Hex())
	}
}
