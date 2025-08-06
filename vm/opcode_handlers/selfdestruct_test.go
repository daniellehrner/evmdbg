package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

// MockStateProviderForSelfDestruct for testing SELFDESTRUCT opcode
type MockStateProviderForSelfDestruct struct {
	accounts map[[20]byte]*MockSelfDestructAccount
	nonces   map[[20]byte]uint64
	deleted  map[[20]byte]bool
}

type MockSelfDestructAccount struct {
	balance *uint256.Int
	code    []byte
	storage map[string]*uint256.Int
	exists  bool
}

func NewMockStateProviderForSelfDestruct() *MockStateProviderForSelfDestruct {
	return &MockStateProviderForSelfDestruct{
		accounts: make(map[[20]byte]*MockSelfDestructAccount),
		nonces:   make(map[[20]byte]uint64),
		deleted:  make(map[[20]byte]bool),
	}
}

func (m *MockStateProviderForSelfDestruct) GetBalance(addr [20]byte) *uint256.Int {
	if m.deleted[addr] {
		return uint256.NewInt(0)
	}
	if acc, exists := m.accounts[addr]; exists {
		return new(uint256.Int).Set(acc.balance)
	}
	return uint256.NewInt(0)
}

func (m *MockStateProviderForSelfDestruct) SetBalance(addr [20]byte, balance *uint256.Int) {
	if m.deleted[addr] {
		return
	}
	if acc, exists := m.accounts[addr]; exists {
		acc.balance = new(uint256.Int).Set(balance)
	} else {
		m.accounts[addr] = &MockSelfDestructAccount{
			balance: new(uint256.Int).Set(balance),
			storage: make(map[string]*uint256.Int),
			exists:  true,
		}
	}
}

func (m *MockStateProviderForSelfDestruct) GetCode(addr [20]byte) []byte {
	if m.deleted[addr] {
		return nil
	}
	if acc, exists := m.accounts[addr]; exists {
		return acc.code
	}
	return nil
}

func (m *MockStateProviderForSelfDestruct) GetStorage(addr [20]byte, key *uint256.Int) *uint256.Int {
	if m.deleted[addr] {
		return uint256.NewInt(0)
	}
	if acc, exists := m.accounts[addr]; exists {
		keyStr := key.Hex()
		if val, ok := acc.storage[keyStr]; ok {
			return new(uint256.Int).Set(val)
		}
	}
	return uint256.NewInt(0)
}

func (m *MockStateProviderForSelfDestruct) SetStorage(addr [20]byte, key *uint256.Int, value *uint256.Int) {
	if m.deleted[addr] {
		return
	}
	if acc, exists := m.accounts[addr]; exists {
		if acc.storage == nil {
			acc.storage = make(map[string]*uint256.Int)
		}
		acc.storage[key.Hex()] = new(uint256.Int).Set(value)
	}
}

func (m *MockStateProviderForSelfDestruct) AccountExists(addr [20]byte) bool {
	if m.deleted[addr] {
		return false
	}
	if acc, exists := m.accounts[addr]; exists {
		return acc.exists
	}
	return false
}

func (m *MockStateProviderForSelfDestruct) GetBlockHash(blockNumber uint64) [32]byte {
	return [32]byte{}
}

func (m *MockStateProviderForSelfDestruct) CreateAccount(addr [20]byte, code []byte, balance *uint256.Int) error {
	m.accounts[addr] = &MockSelfDestructAccount{
		balance: new(uint256.Int).Set(balance),
		code:    code,
		storage: make(map[string]*uint256.Int),
		exists:  true,
	}
	return nil
}

func (m *MockStateProviderForSelfDestruct) GetNonce(addr [20]byte) uint64 {
	if m.deleted[addr] {
		return 0
	}
	return m.nonces[addr]
}

func (m *MockStateProviderForSelfDestruct) SetNonce(addr [20]byte, nonce uint64) {
	if !m.deleted[addr] {
		m.nonces[addr] = nonce
	}
}

func (m *MockStateProviderForSelfDestruct) DeleteAccount(addr [20]byte) error {
	m.deleted[addr] = true
	delete(m.accounts, addr)
	delete(m.nonces, addr)
	return nil
}

func (m *MockStateProviderForSelfDestruct) AddAccount(addr [20]byte, code []byte, balance *uint256.Int) {
	m.accounts[addr] = &MockSelfDestructAccount{
		balance: new(uint256.Int).Set(balance),
		code:    code,
		storage: make(map[string]*uint256.Int),
		exists:  true,
	}
}

func (m *MockStateProviderForSelfDestruct) IsDeleted(addr [20]byte) bool {
	return m.deleted[addr]
}

func TestSelfDestructOpCode_CreatedInTransaction(t *testing.T) {
	// Test SELFDESTRUCT when contract was created in same transaction (full deletion)
	contractAddr := [20]byte{0xaa, 0xbb, 0xcc}
	beneficiaryAddr := [20]byte{0x11, 0x22, 0x33}

	code := []byte{
		0x73,                                           // PUSH20 beneficiary
		0x11, 0x22, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, // beneficiary address
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xff, // SELFDESTRUCT
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderForSelfDestruct()
	v.StateProvider = mockState

	// Set up contract with some balance
	mockState.AddAccount(contractAddr, []byte{0x60, 0x01}, uint256.NewInt(1000))
	mockState.AddAccount(beneficiaryAddr, []byte{}, uint256.NewInt(500))

	// Mark contract as created in current transaction
	v.MarkAccountCreatedInTransaction(contractAddr)

	v.Context = &vm.ExecutionContext{
		Address: contractAddr,
		Block:   &vm.BlockContext{},
	}

	// Execute PUSH20
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH20: %v", err)
	}

	// Execute SELFDESTRUCT
	err = v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during SELFDESTRUCT: %v", err)
	}

	// Check that execution stopped
	if !v.Stopped {
		t.Error("Expected VM to be stopped after SELFDESTRUCT")
	}

	// Check that contract was deleted
	if !mockState.IsDeleted(contractAddr) {
		t.Error("Expected contract to be deleted when created in same transaction")
	}

	// Check beneficiary balance increased
	beneficiaryBalance := mockState.GetBalance(beneficiaryAddr)
	if beneficiaryBalance.Cmp(uint256.NewInt(1500)) != 0 { // 500 + 1000
		t.Errorf("Expected beneficiary balance 1500, got %s", beneficiaryBalance.String())
	}
}

func TestSelfDestructOpCode_NotCreatedInTransaction(t *testing.T) {
	// Test SELFDESTRUCT when contract was NOT created in same transaction (balance transfer only)
	contractAddr := [20]byte{0xaa, 0xbb, 0xcc}
	beneficiaryAddr := [20]byte{0x11, 0x22, 0x33}

	code := []byte{
		0x73,                                           // PUSH20 beneficiary
		0x11, 0x22, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, // beneficiary address
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xff, // SELFDESTRUCT
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderForSelfDestruct()
	v.StateProvider = mockState

	// Set up contract with some balance (NOT marked as created in transaction)
	mockState.AddAccount(contractAddr, []byte{0x60, 0x01}, uint256.NewInt(1000))
	mockState.AddAccount(beneficiaryAddr, []byte{}, uint256.NewInt(500))

	v.Context = &vm.ExecutionContext{
		Address: contractAddr,
		Block:   &vm.BlockContext{},
	}

	// Execute PUSH20
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH20: %v", err)
	}

	// Execute SELFDESTRUCT
	err = v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during SELFDESTRUCT: %v", err)
	}

	// Check that execution stopped
	if !v.Stopped {
		t.Error("Expected VM to be stopped after SELFDESTRUCT")
	}

	// Check that contract was NOT deleted (EIP-6780)
	if mockState.IsDeleted(contractAddr) {
		t.Error("Expected contract NOT to be deleted when not created in same transaction")
	}

	// Check contract balance is now 0
	contractBalance := mockState.GetBalance(contractAddr)
	if !contractBalance.IsZero() {
		t.Errorf("Expected contract balance to be 0, got %s", contractBalance.String())
	}

	// Check beneficiary balance increased
	beneficiaryBalance := mockState.GetBalance(beneficiaryAddr)
	if beneficiaryBalance.Cmp(uint256.NewInt(1500)) != 0 { // 500 + 1000
		t.Errorf("Expected beneficiary balance 1500, got %s", beneficiaryBalance.String())
	}
}

func TestSelfDestructOpCode_SameBeneficiary_CreatedInTransaction(t *testing.T) {
	// Test SELFDESTRUCT with same address as beneficiary, created in transaction (burns ether)
	contractAddr := [20]byte{0xaa, 0xbb, 0xcc}

	code := []byte{
		0x73,                                           // PUSH20 beneficiary (same as contract)
		0xaa, 0xbb, 0xcc, 0x00, 0x00, 0x00, 0x00, 0x00, // contract address
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xff, // SELFDESTRUCT
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderForSelfDestruct()
	v.StateProvider = mockState

	// Set up contract with balance
	mockState.AddAccount(contractAddr, []byte{0x60, 0x01}, uint256.NewInt(1000))

	// Mark contract as created in current transaction
	v.MarkAccountCreatedInTransaction(contractAddr)

	v.Context = &vm.ExecutionContext{
		Address: contractAddr,
		Block:   &vm.BlockContext{},
	}

	// Execute instructions
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check that contract was deleted (ether burned)
	if !mockState.IsDeleted(contractAddr) {
		t.Error("Expected contract to be deleted, burning ether")
	}
}

func TestSelfDestructOpCode_SameBeneficiary_NotCreatedInTransaction(t *testing.T) {
	// Test SELFDESTRUCT with same address as beneficiary, NOT created in transaction (no change)
	contractAddr := [20]byte{0xaa, 0xbb, 0xcc}

	code := []byte{
		0x73,                                           // PUSH20 beneficiary (same as contract)
		0xaa, 0xbb, 0xcc, 0x00, 0x00, 0x00, 0x00, 0x00, // contract address
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xff, // SELFDESTRUCT
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderForSelfDestruct()
	v.StateProvider = mockState

	// Set up contract with balance (NOT created in transaction)
	mockState.AddAccount(contractAddr, []byte{0x60, 0x01}, uint256.NewInt(1000))

	v.Context = &vm.ExecutionContext{
		Address: contractAddr,
		Block:   &vm.BlockContext{},
	}

	// Execute instructions
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check that contract was NOT deleted and balance unchanged
	if mockState.IsDeleted(contractAddr) {
		t.Error("Expected contract NOT to be deleted")
	}

	contractBalance := mockState.GetBalance(contractAddr)
	if contractBalance.Cmp(uint256.NewInt(1000)) != 0 {
		t.Errorf("Expected contract balance unchanged at 1000, got %s", contractBalance.String())
	}
}

func TestSelfDestructOpCode_StaticCallContext(t *testing.T) {
	contractAddr := [20]byte{0xaa, 0xbb, 0xcc}

	code := []byte{
		0x73,                                           // PUSH20
		0x11, 0x22, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00, // beneficiary address
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xff, // SELFDESTRUCT
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderForSelfDestruct()
	v.StateProvider = mockState

	v.Context = &vm.ExecutionContext{
		Address: contractAddr,
		Block:   &vm.BlockContext{},
	}

	// Set static call context
	frame := v.CurrentFrame()
	frame.IsStatic = true

	// Execute PUSH20
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH20: %v", err)
	}

	// SELFDESTRUCT should fail in static context
	err = v.Step()
	if err == nil {
		t.Fatal("Expected error when SELFDESTRUCT is called in static context, got nil")
	}

	if err != vm.ErrStaticCallStateChange {
		t.Errorf("Expected ErrStaticCallStateChange, got: %v", err)
	}
}

func TestSelfDestructOpCode_NoStateProvider(t *testing.T) {
	code := []byte{
		0x73, // PUSH20
		0x11, 0x22, 0x33, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xff, // SELFDESTRUCT
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set StateProvider

	v.Context = &vm.ExecutionContext{
		Address: [20]byte{0xaa, 0xbb, 0xcc},
		Block:   &vm.BlockContext{},
	}

	// Execute PUSH20
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH20: %v", err)
	}

	// SELFDESTRUCT should fail without state provider
	err = v.Step()
	if err == nil {
		t.Fatal("Expected error when state provider is not set, got nil")
	}

	if !strings.Contains(err.Error(), "state provider") {
		t.Errorf("Expected error to mention 'state provider', got: %v", err)
	}
}

func TestSelfDestructOpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0xff, // SELFDESTRUCT (no beneficiary on stack)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	mockState := NewMockStateProviderForSelfDestruct()
	v.StateProvider = mockState

	v.Context = &vm.ExecutionContext{
		Address: [20]byte{0xaa, 0xbb, 0xcc},
		Block:   &vm.BlockContext{},
	}

	// SELFDESTRUCT should fail with stack underflow
	err := v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}
