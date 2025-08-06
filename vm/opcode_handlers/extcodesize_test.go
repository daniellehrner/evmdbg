package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

// Mock state provider for testing
type mockStateProvider struct {
	codeMap map[[20]byte][]byte
}

func (m *mockStateProvider) GetBalance(addr [20]byte) *uint256.Int {
	return uint256.NewInt(0)
}

func (m *mockStateProvider) GetCode(addr [20]byte) []byte {
	if code, exists := m.codeMap[addr]; exists {
		return code
	}
	return []byte{}
}

func (m *mockStateProvider) GetStorage(addr [20]byte, key *uint256.Int) *uint256.Int {
	return uint256.NewInt(0)
}

func (m *mockStateProvider) SetStorage(addr [20]byte, key *uint256.Int, value *uint256.Int) {}

func (m *mockStateProvider) AccountExists(addr [20]byte) bool {
	_, exists := m.codeMap[addr]
	return exists
}

func (m *mockStateProvider) GetBlockHash(blockNumber uint64) [32]byte {
	// Return a simple hash for testing
	var hash [32]byte
	hash[0] = byte(blockNumber)      // Low byte of block number
	hash[1] = byte(blockNumber >> 8) // High byte of block number
	return hash
}

func TestExtCodeSizeOpCode_Execute(t *testing.T) {
	// Test address: 0x1234567890123456789012345678901234567890
	testAddr := [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
	testCode := []byte{0x60, 0x01, 0x60, 0x02, 0x01} // Simple bytecode

	code := []byte{
		0x73,                                                       // PUSH20
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, // address bytes
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90,
		0x3b, // EXTCODESIZE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider
	mock := &mockStateProvider{
		codeMap: map[[20]byte][]byte{
			testAddr: testCode,
		},
	}
	v.StateProvider = mock

	// Execute PUSH20
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH20: %v", err)
	}

	// Execute EXTCODESIZE
	err = v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during EXTCODESIZE: %v", err)
	}

	// Check result
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	size, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	expectedSize := uint64(len(testCode))
	if size.Uint64() != expectedSize {
		t.Errorf("Expected code size %d, got %d", expectedSize, size.Uint64())
	}
}

func TestExtCodeSizeOpCode_NonExistentAccount(t *testing.T) {
	code := []byte{
		0x73,                                                       // PUSH20
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // non-existent address
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		0x3b, // EXTCODESIZE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider with no accounts
	mock := &mockStateProvider{
		codeMap: map[[20]byte][]byte{},
	}
	v.StateProvider = mock

	// Execute PUSH20 and EXTCODESIZE
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

	size, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !size.IsZero() {
		t.Errorf("Expected code size 0 for non-existent account, got %d", size.Uint64())
	}
}

func TestExtCodeSizeOpCode_NoStateProvider(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01 (dummy address)
		0x3b, // EXTCODESIZE
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set StateProvider (remains nil)

	// Execute PUSH1 and EXTCODESIZE
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

	size, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !size.IsZero() {
		t.Errorf("Expected code size 0 with no state provider, got %d", size.Uint64())
	}
}
