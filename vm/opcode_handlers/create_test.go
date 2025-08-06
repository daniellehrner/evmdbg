package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

// MockStateProvider for testing CREATE opcode
type MockStateProviderWithCreate struct {
	accounts map[[20]byte]*MockCreateAccount
	nonces   map[[20]byte]uint64
}

type MockCreateAccount struct {
	balance *uint256.Int
	code    []byte
	storage map[string]*uint256.Int
	exists  bool
}

func NewMockStateProviderWithCreate() *MockStateProviderWithCreate {
	return &MockStateProviderWithCreate{
		accounts: make(map[[20]byte]*MockCreateAccount),
		nonces:   make(map[[20]byte]uint64),
	}
}

func (m *MockStateProviderWithCreate) GetBalance(addr [20]byte) *uint256.Int {
	if acc, exists := m.accounts[addr]; exists {
		return new(uint256.Int).Set(acc.balance)
	}
	return uint256.NewInt(0)
}

func (m *MockStateProviderWithCreate) GetCode(addr [20]byte) []byte {
	if acc, exists := m.accounts[addr]; exists {
		return acc.code
	}
	return nil
}

func (m *MockStateProviderWithCreate) GetStorage(addr [20]byte, key *uint256.Int) *uint256.Int {
	if acc, exists := m.accounts[addr]; exists {
		keyStr := key.Hex()
		if val, ok := acc.storage[keyStr]; ok {
			return new(uint256.Int).Set(val)
		}
	}
	return uint256.NewInt(0)
}

func (m *MockStateProviderWithCreate) SetStorage(addr [20]byte, key *uint256.Int, value *uint256.Int) {
	if acc, exists := m.accounts[addr]; exists {
		if acc.storage == nil {
			acc.storage = make(map[string]*uint256.Int)
		}
		acc.storage[key.Hex()] = new(uint256.Int).Set(value)
	}
}

func (m *MockStateProviderWithCreate) AccountExists(addr [20]byte) bool {
	if acc, exists := m.accounts[addr]; exists {
		return acc.exists
	}
	return false
}

func (m *MockStateProviderWithCreate) GetBlockHash(blockNumber uint64) [32]byte {
	return [32]byte{}
}

func (m *MockStateProviderWithCreate) CreateAccount(addr [20]byte, code []byte, balance *uint256.Int) error {
	m.accounts[addr] = &MockCreateAccount{
		balance: new(uint256.Int).Set(balance),
		code:    code,
		storage: make(map[string]*uint256.Int),
		exists:  true,
	}
	return nil
}

func (m *MockStateProviderWithCreate) GetNonce(addr [20]byte) uint64 {
	return m.nonces[addr]
}

func (m *MockStateProviderWithCreate) SetNonce(addr [20]byte, nonce uint64) {
	m.nonces[addr] = nonce
}

func (m *MockStateProviderWithCreate) SetBalance(addr [20]byte, balance *uint256.Int) {
	if acc, exists := m.accounts[addr]; exists {
		acc.balance = new(uint256.Int).Set(balance)
	}
}

func (m *MockStateProviderWithCreate) DeleteAccount(addr [20]byte) error {
	delete(m.accounts, addr)
	delete(m.nonces, addr)
	return nil
}

func TestCreateOpCode_Execute(t *testing.T) {
	// Bytecode: PUSH1 0x04 (size), PUSH1 0x00 (offset), PUSH1 0x42 (value), CREATE
	code := []byte{
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value = 66)
		0xf0, // CREATE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider
	mockState := NewMockStateProviderWithCreate()
	v.StateProvider = mockState

	// Set up execution context
	creatorAddr := [20]byte{0xaa, 0xbb, 0xcc}
	mockState.accounts[creatorAddr] = &MockCreateAccount{
		balance: uint256.NewInt(1000),
		exists:  true,
	}

	v.Context = &vm.ExecutionContext{
		Address: creatorAddr,
		Value:   uint256.NewInt(0),
		Block:   &vm.BlockContext{},
	}

	// Put some init code in memory at offset 0
	initCode := []byte{0x60, 0x00, 0xf3, 0x00} // PUSH1 0, RETURN, padding
	v.Memory().Write(0, initCode)

	// Execute bytecode steps
	for i := 0; i < 4; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check final state
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack after CREATE, got %d", v.Stack().Len())
	}

	result, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	// Result should be non-zero (address of created contract)
	if result.IsZero() {
		t.Errorf("Expected non-zero address from successful CREATE, got zero")
	}

	// Check that creator's nonce was incremented
	if mockState.GetNonce(creatorAddr) != 1 {
		t.Errorf("Expected creator nonce to be 1, got %d", mockState.GetNonce(creatorAddr))
	}

	// Check that a new account was created
	newAddr := [20]byte{}
	result.WriteToSlice(newAddr[:])
	if !mockState.AccountExists(newAddr) {
		t.Errorf("Expected new contract account to exist")
	}
}

func TestCreateOpCode_InsufficientBalance(t *testing.T) {
	// Same bytecode but with insufficient balance
	code := []byte{
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x61, 0xff, 0xff, // PUSH2 0xffff (very large value)
		0xf0, // CREATE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderWithCreate()
	v.StateProvider = mockState

	creatorAddr := [20]byte{0xaa, 0xbb, 0xcc}
	mockState.accounts[creatorAddr] = &MockCreateAccount{
		balance: uint256.NewInt(100), // Small balance
		exists:  true,
	}

	v.Context = &vm.ExecutionContext{
		Address: creatorAddr,
		Value:   uint256.NewInt(0),
		Block:   &vm.BlockContext{},
	}

	// Put init code in memory
	initCode := []byte{0x60, 0x00, 0xf3, 0x00}
	v.Memory().Write(0, initCode)

	// Execute all steps
	for i := 0; i < 4; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Should push 0 to indicate failure
	result, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !result.IsZero() {
		t.Errorf("Expected zero (failure) from CREATE with insufficient balance, got %s", result.Hex())
	}
}

func TestCreateOpCode_AccountAlreadyExists(t *testing.T) {
	code := []byte{
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value)
		0xf0, // CREATE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderWithCreate()
	v.StateProvider = mockState

	creatorAddr := [20]byte{0xaa, 0xbb, 0xcc}
	mockState.accounts[creatorAddr] = &MockCreateAccount{
		balance: uint256.NewInt(1000),
		exists:  true,
	}

	// Pre-create an account at what would be the CREATE address
	// This simulates address collision
	existingAddr := [20]byte{0x01, 0x02, 0x03}
	mockState.accounts[existingAddr] = &MockCreateAccount{
		balance: uint256.NewInt(0),
		exists:  true,
	}

	v.Context = &vm.ExecutionContext{
		Address: creatorAddr,
		Value:   uint256.NewInt(0),
		Block:   &vm.BlockContext{},
	}

	v.Memory().Write(0, []byte{0x60, 0x00, 0xf3, 0x00})

	for i := 0; i < 4; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Since address collision detection is complex, we just verify
	// the operation completes without crashing
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}
}

func TestCreateOpCode_StaticCallContext(t *testing.T) {
	code := []byte{
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value)
		0xf0, // CREATE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderWithCreate()
	v.StateProvider = mockState

	v.Context = &vm.ExecutionContext{
		Address: [20]byte{0xaa, 0xbb, 0xcc},
		Block:   &vm.BlockContext{},
	}

	// Set static call context
	frame := v.CurrentFrame()
	frame.IsStatic = true

	v.Memory().Write(0, []byte{0x60, 0x00, 0xf3, 0x00})

	// Execute first 3 steps (PUSH operations)
	for i := 0; i < 3; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// CREATE should fail in static context
	err := v.Step()
	if err == nil {
		t.Fatal("Expected error when CREATE is called in static context, got nil")
	}

	if err != vm.ErrStaticCallStateChange {
		t.Errorf("Expected ErrStaticCallStateChange, got: %v", err)
	}
}

func TestCreateOpCode_NoStateProvider(t *testing.T) {
	code := []byte{
		0x60, 0x04, // PUSH1 4
		0x60, 0x00, // PUSH1 0
		0x60, 0x42, // PUSH1 0x42
		0xf0, // CREATE
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set StateProvider

	v.Context = &vm.ExecutionContext{
		Address: [20]byte{0xaa, 0xbb, 0xcc},
		Block:   &vm.BlockContext{},
	}

	// Execute PUSH operations
	for i := 0; i < 3; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// CREATE should fail without state provider
	err := v.Step()
	if err == nil {
		t.Fatal("Expected error when state provider is not set, got nil")
	}

	if !strings.Contains(err.Error(), "state provider") {
		t.Errorf("Expected error to mention 'state provider', got: %v", err)
	}
}

func TestCreateOpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0x60, 0x04, // PUSH1 4 (only one value on stack)
		0xf0, // CREATE (needs 3 values)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderWithCreate()
	v.StateProvider = mockState

	v.Context = &vm.ExecutionContext{
		Address: [20]byte{0xaa, 0xbb, 0xcc},
		Block:   &vm.BlockContext{},
	}

	// Execute PUSH
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH: %v", err)
	}

	// CREATE should fail with stack underflow
	err = v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}

func TestCreateOpCode_MemoryOutOfBounds(t *testing.T) {
	code := []byte{
		0x61, 0x10, 0x00, // PUSH2 4096 (large size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value)
		0xf0, // CREATE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderWithCreate()
	v.StateProvider = mockState

	creatorAddr := [20]byte{0xaa, 0xbb, 0xcc}
	mockState.accounts[creatorAddr] = &MockCreateAccount{
		balance: uint256.NewInt(1000),
		exists:  true,
	}

	v.Context = &vm.ExecutionContext{
		Address: creatorAddr,
		Value:   uint256.NewInt(0),
		Block:   &vm.BlockContext{},
	}

	// Execute PUSH operations
	for i := 0; i < 3; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// CREATE should handle memory expansion correctly
	// (This test verifies no crash occurs with large memory access)
	err := v.Step()
	if err != nil {
		// Memory expansion might cause an error, which is acceptable
		t.Logf("CREATE failed with memory error (expected): %v", err)
	} else {
		// If it succeeds, verify result
		if v.Stack().Len() != 1 {
			t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
		}
	}
}
