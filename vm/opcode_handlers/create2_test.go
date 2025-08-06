package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestCreate2OpCode_Execute(t *testing.T) {
	// Bytecode: PUSH1 0x1234 (salt), PUSH1 0x04 (size), PUSH1 0x00 (offset), PUSH1 0x42 (value), CREATE2
	code := []byte{
		0x61, 0x12, 0x34, // PUSH2 0x1234 (salt)
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value = 66)
		0xf5, // CREATE2
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
	for i := 0; i < 5; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check final state
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack after CREATE2, got %d", v.Stack().Len())
	}

	result, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	// Result should be non-zero (address of created contract)
	if result.IsZero() {
		t.Errorf("Expected non-zero address from successful CREATE2, got zero")
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

func TestCreate2OpCode_DeterministicAddress(t *testing.T) {
	// Test that CREATE2 with same inputs produces same address

	// First CREATE2
	code1 := []byte{
		0x61, 0x12, 0x34, // PUSH2 0x1234 (salt)
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value)
		0xf5, // CREATE2
	}

	v1 := vm.NewDebuggerVM(code1, GetHandler)
	mockState1 := NewMockStateProviderWithCreate()
	v1.StateProvider = mockState1

	creatorAddr := [20]byte{0xaa, 0xbb, 0xcc}
	mockState1.accounts[creatorAddr] = &MockCreateAccount{
		balance: uint256.NewInt(1000),
		exists:  true,
	}

	v1.Context = &vm.ExecutionContext{
		Address: creatorAddr,
		Value:   uint256.NewInt(0),
		Block:   &vm.BlockContext{},
	}

	initCode := []byte{0x60, 0x00, 0xf3, 0x00}
	v1.Memory().Write(0, initCode)

	// Execute first CREATE2
	for i := 0; i < 5; i++ {
		err := v1.Step()
		if err != nil {
			t.Fatalf("Unexpected error in first CREATE2 step %d: %v", i, err)
		}
	}

	addr1, err := v1.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error getting first CREATE2 result: %v", err)
	}

	// Second CREATE2 with same inputs but fresh VM
	code2 := []byte{
		0x61, 0x12, 0x34, // PUSH2 0x1234 (same salt)
		0x60, 0x04, // PUSH1 4 (same size)
		0x60, 0x00, // PUSH1 0 (same offset)
		0x60, 0x42, // PUSH1 0x42 (same value)
		0xf5, // CREATE2
	}

	v2 := vm.NewDebuggerVM(code2, GetHandler)
	mockState2 := NewMockStateProviderWithCreate()
	v2.StateProvider = mockState2

	mockState2.accounts[creatorAddr] = &MockCreateAccount{
		balance: uint256.NewInt(1000),
		exists:  true,
	}

	v2.Context = &vm.ExecutionContext{
		Address: creatorAddr,
		Value:   uint256.NewInt(0),
		Block:   &vm.BlockContext{},
	}

	v2.Memory().Write(0, initCode) // Same init code

	// Execute second CREATE2
	for i := 0; i < 5; i++ {
		err := v2.Step()
		if err != nil {
			t.Fatalf("Unexpected error in second CREATE2 step %d: %v", i, err)
		}
	}

	addr2, err := v2.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error getting second CREATE2 result: %v", err)
	}

	// Addresses should be the same (deterministic)
	if addr1.Cmp(addr2) != 0 {
		t.Errorf("CREATE2 addresses should be deterministic: got %s and %s", addr1.Hex(), addr2.Hex())
	}
}

func TestCreate2OpCode_DifferentSalt(t *testing.T) {
	// Test that CREATE2 with different salts produces different addresses
	creatorAddr := [20]byte{0xaa, 0xbb, 0xcc}
	initCode := []byte{0x60, 0x00, 0xf3, 0x00}

	// First CREATE2 with salt 0x1234
	code1 := []byte{
		0x61, 0x12, 0x34, // PUSH2 0x1234 (salt)
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value)
		0xf5, // CREATE2
	}

	v1 := vm.NewDebuggerVM(code1, GetHandler)
	mockState1 := NewMockStateProviderWithCreate()
	v1.StateProvider = mockState1
	mockState1.accounts[creatorAddr] = &MockCreateAccount{balance: uint256.NewInt(1000), exists: true}
	v1.Context = &vm.ExecutionContext{Address: creatorAddr, Value: uint256.NewInt(0), Block: &vm.BlockContext{}}
	v1.Memory().Write(0, initCode)

	for i := 0; i < 5; i++ {
		v1.Step()
	}
	addr1, _ := v1.Stack().Peek(0)

	// Second CREATE2 with salt 0x5678
	code2 := []byte{
		0x61, 0x56, 0x78, // PUSH2 0x5678 (different salt)
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value)
		0xf5, // CREATE2
	}

	v2 := vm.NewDebuggerVM(code2, GetHandler)
	mockState2 := NewMockStateProviderWithCreate()
	v2.StateProvider = mockState2
	mockState2.accounts[creatorAddr] = &MockCreateAccount{balance: uint256.NewInt(1000), exists: true}
	v2.Context = &vm.ExecutionContext{Address: creatorAddr, Value: uint256.NewInt(0), Block: &vm.BlockContext{}}
	v2.Memory().Write(0, initCode)

	for i := 0; i < 5; i++ {
		v2.Step()
	}
	addr2, _ := v2.Stack().Peek(0)

	// Addresses should be different
	if addr1.Cmp(addr2) == 0 {
		t.Errorf("CREATE2 with different salts should produce different addresses: both got %s", addr1.Hex())
	}
}

func TestCreate2OpCode_InsufficientBalance(t *testing.T) {
	code := []byte{
		0x61, 0x12, 0x34, // PUSH2 0x1234 (salt)
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x62, 0xff, 0xff, 0xff, // PUSH3 0xffffff (very large value)
		0xf5, // CREATE2
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

	v.Memory().Write(0, []byte{0x60, 0x00, 0xf3, 0x00})

	// Execute all steps
	for i := 0; i < 5; i++ {
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
		t.Errorf("Expected zero (failure) from CREATE2 with insufficient balance, got %s", result.Hex())
	}
}

func TestCreate2OpCode_StaticCallContext(t *testing.T) {
	code := []byte{
		0x61, 0x12, 0x34, // PUSH2 0x1234 (salt)
		0x60, 0x04, // PUSH1 4 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0x60, 0x42, // PUSH1 0x42 (value)
		0xf5, // CREATE2
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

	// Execute first 4 steps (PUSH operations)
	for i := 0; i < 4; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// CREATE2 should fail in static context
	err := v.Step()
	if err == nil {
		t.Fatal("Expected error when CREATE2 is called in static context, got nil")
	}

	if err != vm.ErrStaticCallStateChange {
		t.Errorf("Expected ErrStaticCallStateChange, got: %v", err)
	}
}

func TestCreate2OpCode_NoStateProvider(t *testing.T) {
	code := []byte{
		0x61, 0x12, 0x34, // PUSH2 0x1234
		0x60, 0x04, // PUSH1 4
		0x60, 0x00, // PUSH1 0
		0x60, 0x42, // PUSH1 0x42
		0xf5, // CREATE2
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set StateProvider

	v.Context = &vm.ExecutionContext{
		Address: [20]byte{0xaa, 0xbb, 0xcc},
		Block:   &vm.BlockContext{},
	}

	// Execute PUSH operations
	for i := 0; i < 4; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// CREATE2 should fail without state provider
	err := v.Step()
	if err == nil {
		t.Fatal("Expected error when state provider is not set, got nil")
	}

	if !strings.Contains(err.Error(), "state provider") {
		t.Errorf("Expected error to mention 'state provider', got: %v", err)
	}
}

func TestCreate2OpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0x60, 0x04, // PUSH1 4 (only one value on stack)
		0xf5, // CREATE2 (needs 4 values)
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

	// CREATE2 should fail with stack underflow
	err = v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}
