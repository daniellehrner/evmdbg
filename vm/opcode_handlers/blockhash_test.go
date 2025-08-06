package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

// mockStateProviderForBlockHash extends the mock with GetBlockHash
type mockStateProviderForBlockHash struct {
	blockHashes map[uint64][32]byte
}

func (m *mockStateProviderForBlockHash) GetBalance(addr [20]byte) *uint256.Int {
	return uint256.NewInt(0)
}

func (m *mockStateProviderForBlockHash) GetCode(addr [20]byte) []byte {
	return []byte{}
}

func (m *mockStateProviderForBlockHash) GetStorage(addr [20]byte, key *uint256.Int) *uint256.Int {
	return uint256.NewInt(0)
}

func (m *mockStateProviderForBlockHash) SetStorage(addr [20]byte, key *uint256.Int, value *uint256.Int) {
}

func (m *mockStateProviderForBlockHash) AccountExists(addr [20]byte) bool {
	return false
}

func (m *mockStateProviderForBlockHash) GetBlockHash(blockNumber uint64) [32]byte {
	if hash, exists := m.blockHashes[blockNumber]; exists {
		return hash
	}
	return [32]byte{} // Return zero hash if not found
}

func (m *mockStateProviderForBlockHash) CreateAccount(addr [20]byte, code []byte, balance *uint256.Int) error {
	return nil // No-op for block hash tests
}

func (m *mockStateProviderForBlockHash) GetNonce(addr [20]byte) uint64 {
	return 0
}

func (m *mockStateProviderForBlockHash) SetNonce(addr [20]byte, nonce uint64) {
}

func (m *mockStateProviderForBlockHash) SetBalance(addr [20]byte, balance *uint256.Int) {
}

func (m *mockStateProviderForBlockHash) DeleteAccount(addr [20]byte) error {
	return nil
}

func TestBlockHashOpCode_Execute(t *testing.T) {
	// Test getting hash for block 100 when current block is 150
	code := []byte{
		0x60, 0x64, // PUSH1 0x64 (block 100)
		0x40, // BLOCKHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			Number: 150, // Current block number
		},
	}

	// Set up mock state provider with known block hash
	expectedHash := [32]byte{0xaa, 0xbb, 0xcc, 0xdd}
	for i := 4; i < 32; i++ {
		expectedHash[i] = 0x00
	}

	mock := &mockStateProviderForBlockHash{
		blockHashes: map[uint64][32]byte{
			100: expectedHash,
		},
	}
	v.StateProvider = mock

	// Execute PUSH1 and BLOCKHASH
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

	expectedInt := new(uint256.Int).SetBytes(expectedHash[:])
	if hash.Cmp(expectedInt) != 0 {
		t.Errorf("Expected block hash %s, got %s", expectedInt.Hex(), hash.Hex())
	}
}

func TestBlockHashOpCode_FutureBlock(t *testing.T) {
	// Test getting hash for future block (should return 0)
	code := []byte{
		0x60, 0x96, // PUSH1 0x96 (block 150)
		0x40, // BLOCKHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context - requesting current block
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			Number: 150, // Current block number
		},
	}

	// Set up mock state provider
	mock := &mockStateProviderForBlockHash{
		blockHashes: map[uint64][32]byte{},
	}
	v.StateProvider = mock

	// Execute PUSH1 and BLOCKHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should be 0 for future/current block
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !hash.IsZero() {
		t.Errorf("Expected block hash 0 for future block, got %s", hash.Hex())
	}
}

func TestBlockHashOpCode_TooOldBlock(t *testing.T) {
	// Test getting hash for block too far in the past (should return 0)
	code := []byte{
		0x60, 0x01, // PUSH1 0x01 (block 1)
		0x40, // BLOCKHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context - current block is way ahead (more than 256 blocks)
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			Number: 300, // Current block number (300 - 1 > 256)
		},
	}

	// Set up mock state provider
	mock := &mockStateProviderForBlockHash{
		blockHashes: map[uint64][32]byte{
			1: {0xaa, 0xbb, 0xcc, 0xdd}, // This should not be returned
		},
	}
	v.StateProvider = mock

	// Execute PUSH1 and BLOCKHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should be 0 for block too far in the past
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !hash.IsZero() {
		t.Errorf("Expected block hash 0 for block too far in past, got %s", hash.Hex())
	}
}

func TestBlockHashOpCode_NoContext(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01
		0x40, // BLOCKHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set context

	// Execute PUSH1
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH1: %v", err)
	}

	// BLOCKHASH should fail without context
	err = v.Step()
	if err == nil {
		t.Fatal("Expected error when context is not set, got nil")
	}

	if !strings.Contains(err.Error(), "execution context") {
		t.Errorf("Expected error to mention 'execution context', got: %v", err)
	}
}
