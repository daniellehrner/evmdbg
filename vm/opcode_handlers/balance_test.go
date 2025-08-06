package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestBalanceOpCode_Execute(t *testing.T) {
	// Test address: 0x1234567890123456789012345678901234567890
	testAddr := [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
	testBalance := uint256.NewInt(123456789)

	code := []byte{
		0x73,                                                       // PUSH20
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, // address bytes
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90,
		0x31, // BALANCE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider with the test account
	mock := NewMockStateProvider()
	mock.AddAccount(testAddr, []byte{}, testBalance)
	v.StateProvider = mock

	// Execute PUSH20
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH20: %v", err)
	}

	// Execute BALANCE
	err = v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during BALANCE: %v", err)
	}

	// Check result
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	balance, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if balance.Cmp(testBalance) != 0 {
		t.Errorf("Expected balance %s, got %s", testBalance.String(), balance.String())
	}
}

func TestBalanceOpCode_NonExistentAccount(t *testing.T) {
	code := []byte{
		0x73,                                                       // PUSH20
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, // non-existent address
		0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		0x31, // BALANCE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider with no accounts
	mock := NewMockStateProvider()
	v.StateProvider = mock

	// Execute PUSH20 and BALANCE
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

	balance, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !balance.IsZero() {
		t.Errorf("Expected balance 0 for non-existent account, got %s", balance.String())
	}
}

func TestBalanceOpCode_ShortAddress(t *testing.T) {
	// Test with address shorter than 20 bytes (should be right-aligned)
	code := []byte{
		0x60, 0x01, // PUSH1 0x01 (short address)
		0x31, // BALANCE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Expected address: 0x0000000000000000000000000000000000000001
	expectedAddr := [20]byte{}
	expectedAddr[19] = 0x01
	testBalance := uint256.NewInt(999)

	// Set up mock state provider
	mock := NewMockStateProvider()
	mock.AddAccount(expectedAddr, []byte{}, testBalance)
	v.StateProvider = mock

	// Execute PUSH1 and BALANCE
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

	balance, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if balance.Cmp(testBalance) != 0 {
		t.Errorf("Expected balance %s for short address, got %s", testBalance.String(), balance.String())
	}
}

func TestBalanceOpCode_LongAddress(t *testing.T) {
	// Test with address longer than 20 bytes (should take last 20 bytes)
	code := []byte{
		0x7f,                                                                   // PUSH32
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55, // first 12 bytes (ignored)
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, // last 20 bytes (used)
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90,
		0x31, // BALANCE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Expected address: last 20 bytes
	expectedAddr := [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
	testBalance := uint256.NewInt(777)

	// Set up mock state provider
	mock := NewMockStateProvider()
	mock.AddAccount(expectedAddr, []byte{}, testBalance)
	v.StateProvider = mock

	// Execute PUSH32 and BALANCE
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

	balance, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if balance.Cmp(testBalance) != 0 {
		t.Errorf("Expected balance %s for long address, got %s", testBalance.String(), balance.String())
	}
}

func TestBalanceOpCode_NoStateProvider(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01 (dummy address)
		0x31, // BALANCE
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set StateProvider (remains nil)

	// Execute PUSH1 and BALANCE
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

	balance, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !balance.IsZero() {
		t.Errorf("Expected balance 0 with no state provider, got %s", balance.String())
	}
}

func TestBalanceOpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0x31, // BALANCE (no address on stack)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// BALANCE should fail with stack underflow
	err := v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}

func TestBalanceOpCode_ZeroBalance(t *testing.T) {
	// Test account with zero balance
	testAddr := [20]byte{0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44}
	zeroBalance := uint256.NewInt(0)

	code := []byte{
		0x73, // PUSH20
		0x99, 0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00,
		0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44,
		0x31, // BALANCE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider with zero balance account
	mock := NewMockStateProvider()
	mock.AddAccount(testAddr, []byte{}, zeroBalance)
	v.StateProvider = mock

	// Execute PUSH20 and BALANCE
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

	balance, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !balance.IsZero() {
		t.Errorf("Expected zero balance, got %s", balance.String())
	}
}
