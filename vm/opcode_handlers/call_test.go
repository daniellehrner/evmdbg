package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

// MockStateProvider implements vm.StateProvider for testing
type MockStateProvider struct {
	accounts map[[20]byte]MockAccount
}

type MockAccount struct {
	code    []byte
	balance *uint256.Int
	storage map[string]*uint256.Int
	exists  bool
}

func NewMockStateProvider() *MockStateProvider {
	return &MockStateProvider{
		accounts: make(map[[20]byte]MockAccount),
	}
}

func (m *MockStateProvider) GetBalance(addr [20]byte) *uint256.Int {
	if acc, exists := m.accounts[addr]; exists {
		return new(uint256.Int).Set(acc.balance)
	}
	return uint256.NewInt(0)
}

func (m *MockStateProvider) GetCode(addr [20]byte) []byte {
	if acc, exists := m.accounts[addr]; exists {
		return acc.code
	}
	return nil
}

func (m *MockStateProvider) GetStorage(addr [20]byte, key *uint256.Int) *uint256.Int {
	if acc, exists := m.accounts[addr]; exists {
		keyStr := key.Hex()
		if val, exists := acc.storage[keyStr]; exists {
			return new(uint256.Int).Set(val)
		}
	}
	return uint256.NewInt(0)
}

func (m *MockStateProvider) SetStorage(addr [20]byte, key *uint256.Int, value *uint256.Int) {
	if acc, exists := m.accounts[addr]; exists {
		if acc.storage == nil {
			acc.storage = make(map[string]*uint256.Int)
		}
		keyStr := key.Hex()
		acc.storage[keyStr] = new(uint256.Int).Set(value)
		m.accounts[addr] = acc
	}
}

func (m *MockStateProvider) AccountExists(addr [20]byte) bool {
	acc, exists := m.accounts[addr]
	return exists && acc.exists
}

func (m *MockStateProvider) GetBlockHash(blockNumber uint64) [32]byte {
	// Simple deterministic hash for testing
	var hash [32]byte
	hash[0] = byte(blockNumber)
	hash[1] = byte(blockNumber >> 8)
	return hash
}

func (m *MockStateProvider) CreateAccount(addr [20]byte, code []byte, balance *uint256.Int) error {
	m.accounts[addr] = MockAccount{
		code:    code,
		balance: new(uint256.Int).Set(balance),
		storage: make(map[string]*uint256.Int),
		exists:  true,
	}
	return nil
}

func (m *MockStateProvider) GetNonce(addr [20]byte) uint64 {
	// For simplicity, return 0 for all addresses in tests
	return 0
}

func (m *MockStateProvider) SetNonce(addr [20]byte, nonce uint64) {
	// For simplicity, ignore nonce updates in these tests
}

func (m *MockStateProvider) AddAccount(addr [20]byte, code []byte, balance *uint256.Int) {
	m.accounts[addr] = MockAccount{
		code:    code,
		balance: new(uint256.Int).Set(balance),
		storage: make(map[string]*uint256.Int),
		exists:  true,
	}
}

func TestCallWithoutStateProvider(t *testing.T) {
	// Test CALL without StateProvider (should return success)
	code := []byte{
		vm.PUSH1, 0x00, // retSize
		vm.PUSH1, 0x00, // retOffset
		vm.PUSH1, 0x00, // argsSize
		vm.PUSH1, 0x00, // argsOffset
		vm.PUSH1, 0x00, // value
		vm.PUSH1, 0x01, // address (0x01)
		vm.PUSH1, 0x64, // gas (100)
		vm.CALL, // CALL
	}

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped && d.PC() < uint64(len(d.Code())) {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Should have success result (1) on stack
	result, err := d.Stack().Pop()
	if err != nil {
		t.Fatalf("failed to pop result: %v", err)
	}

	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestCallWithStateProvider(t *testing.T) {
	// Create mock state provider
	stateProvider := NewMockStateProvider()

	// Add a mock account at address 0x01
	var targetAddr [20]byte
	targetAddr[19] = 0x01
	targetCode := []byte{
		vm.PUSH1, 0x01, // size = 1 byte
		vm.PUSH1, 0x00, // offset = 0
		vm.PUSH1, 0x42, // value to store
		vm.PUSH1, 0x00, // offset for MSTORE
		vm.MSTORE,      // store 0x42 in memory at offset 0
		vm.PUSH1, 0x01, // size = 1 byte
		vm.PUSH1, 0x00, // offset = 0
		vm.RETURN, // return 1 byte from offset 0
	} // Code that returns 0x42
	stateProvider.AddAccount(targetAddr, targetCode, uint256.NewInt(1000))

	// Test CALL with StateProvider
	code := []byte{
		vm.PUSH1, 0x00, // retSize
		vm.PUSH1, 0x00, // retOffset
		vm.PUSH1, 0x00, // argsSize
		vm.PUSH1, 0x00, // argsOffset
		vm.PUSH1, 0x00, // value
		vm.PUSH1, 0x01, // address (0x01)
		vm.PUSH1, 0x64, // gas (100)
		vm.CALL, // CALL
	}

	d := vm.NewDebuggerVM(code, GetHandler)
	d.StateProvider = stateProvider

	// Set up execution context
	var currentAddr [20]byte
	d.Context = &vm.ExecutionContext{
		Caller:   [20]byte{},
		Address:  currentAddr,
		Origin:   [20]byte{},
		Value:    uint256.NewInt(0),
		CallData: []byte{},
		GasPrice: uint256.NewInt(1),
		Gas:      1000,
		Balance:  uint256.NewInt(1000),
		Block: &vm.BlockContext{
			Coinbase:   [20]byte{},
			Timestamp:  1,
			Number:     1,
			Difficulty: uint256.NewInt(1),
			GasLimit:   1000000,
			ChainID:    uint256.NewInt(1),
		},
	}

	for !d.Stopped && d.PC() < uint64(len(d.Code())) {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Should have success result (1) on stack
	result, err := d.Stack().Pop()
	if err != nil {
		t.Fatalf("failed to pop result: %v", err)
	}

	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestDelegateCallStackArgs(t *testing.T) {
	// Test DELEGATECALL has correct stack requirements (6 args, not 7)
	code := []byte{
		vm.PUSH1, 0x00, // retSize
		vm.PUSH1, 0x00, // retOffset
		vm.PUSH1, 0x00, // argsSize
		vm.PUSH1, 0x00, // argsOffset
		vm.PUSH1, 0x01, // address (0x01)
		vm.PUSH1, 0x64, // gas (100)
		vm.DELEGATECALL, // DELEGATECALL
	}

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped && d.PC() < uint64(len(d.Code())) {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Should have success result (1) on stack
	result, err := d.Stack().Pop()
	if err != nil {
		t.Fatalf("failed to pop result: %v", err)
	}

	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestStaticCallStackArgs(t *testing.T) {
	// Test STATICCALL has correct stack requirements (6 args, not 7)
	code := []byte{
		vm.PUSH1, 0x00, // retSize
		vm.PUSH1, 0x00, // retOffset
		vm.PUSH1, 0x00, // argsSize
		vm.PUSH1, 0x00, // argsOffset
		vm.PUSH1, 0x01, // address (0x01)
		vm.PUSH1, 0x64, // gas (100)
		vm.STATICCALL, // STATICCALL
	}

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped && d.PC() < uint64(len(d.Code())) {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Should have success result (1) on stack
	result, err := d.Stack().Pop()
	if err != nil {
		t.Fatalf("failed to pop result: %v", err)
	}

	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestCallDepth(t *testing.T) {
	// Test that call depth is tracked correctly
	code := []byte{
		vm.PUSH1, 0x00, // retSize
		vm.PUSH1, 0x00, // retOffset
		vm.PUSH1, 0x00, // argsSize
		vm.PUSH1, 0x00, // argsOffset
		vm.PUSH1, 0x00, // value
		vm.PUSH1, 0x01, // address (0x01)
		vm.PUSH1, 0x64, // gas (100)
		vm.CALL, // CALL
	}

	d := vm.NewDebuggerVM(code, GetHandler)

	// Initial call depth should be 1 (root frame)
	if d.CallDepth() != 1 {
		t.Fatalf("initial call depth should be 1, got %d", d.CallDepth())
	}

	// Execute until CALL
	for !d.Stopped && d.PC() < uint64(len(d.Code())) {
		if d.Code()[d.PC()] == vm.CALL {
			break
		}
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Execute the CALL
	err := d.Step()
	if err != nil {
		t.Fatalf("CALL execution error: %v", err)
	}

	// Call depth should still be 1 after call returns
	if d.CallDepth() != 1 {
		t.Fatalf("call depth after CALL should be 1, got %d", d.CallDepth())
	}
}
