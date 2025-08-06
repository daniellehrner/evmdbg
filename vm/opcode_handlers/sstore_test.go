package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSstoreBasic(t *testing.T) {
	// Test: SSTORE(key, value) - store value at storage key
	code := []byte{
		vm.PUSH1, 0x42, // PUSH1 0x42 (value)
		vm.PUSH1, 0x01, // PUSH1 0x01 (key)
		vm.SSTORE, // SSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int) // Initialize storage

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check that storage was written correctly
	key := "0000000000000000000000000000000000000000000000000000000000000001"
	storedValue, exists := d.Storage[key]
	if !exists {
		t.Fatalf("storage key %s not found", key)
	}

	expected := uint256.NewInt(0x42)
	if storedValue.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, storedValue)
	}
}

func TestSstoreZero(t *testing.T) {
	// Test: SSTORE(key, 0) - store zero value
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (value)
		vm.PUSH1, 0x05, // PUSH1 5 (key)
		vm.SSTORE, // SSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check that zero was stored
	key := "0000000000000000000000000000000000000000000000000000000000000005"
	storedValue, exists := d.Storage[key]
	if !exists {
		t.Fatalf("storage key %s not found", key)
	}

	if !storedValue.IsZero() {
		t.Fatalf("expected 0, got %s", storedValue)
	}
}

func TestSstoreOverwrite(t *testing.T) {
	// Test: SSTORE can overwrite existing values
	code := []byte{
		// Store first value
		vm.PUSH1, 0x11, // PUSH1 0x11 (value)
		vm.PUSH1, 0x01, // PUSH1 0x01 (key)
		vm.SSTORE, // SSTORE

		// Store second value at same key
		vm.PUSH1, 0x22, // PUSH1 0x22 (new value)
		vm.PUSH1, 0x01, // PUSH1 0x01 (same key)
		vm.SSTORE, // SSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check that the value was overwritten
	key := "0000000000000000000000000000000000000000000000000000000000000001"
	storedValue, exists := d.Storage[key]
	if !exists {
		t.Fatalf("storage key %s not found", key)
	}

	expected := uint256.NewInt(0x22) // Should be the second value
	if storedValue.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, storedValue)
	}
}

func TestSstoreLargeKey(t *testing.T) {
	// Test: SSTORE with large key
	largeKey := new(uint256.Int).SetAllOne() // Max uint256

	code := []byte{
		vm.PUSH1, 0x99, // PUSH1 0x99 (value)
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeKey)...)
	code = append(code, vm.SSTORE)

	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check storage with large key
	key := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	storedValue, exists := d.Storage[key]
	if !exists {
		t.Fatalf("storage key %s not found", key)
	}

	expected := uint256.NewInt(0x99)
	if storedValue.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, storedValue)
	}
}

func TestSstoreLargeValue(t *testing.T) {
	// Test: SSTORE with large value
	largeValue := new(uint256.Int).SetBytes([]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeValue)...)
	code = append(code, vm.PUSH1, 0x10) // PUSH1 16 (key)
	code = append(code, vm.SSTORE)

	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check that large value was stored correctly
	key := "0000000000000000000000000000000000000000000000000000000000000010"
	storedValue, exists := d.Storage[key]
	if !exists {
		t.Fatalf("storage key %s not found", key)
	}

	if storedValue.Cmp(largeValue) != 0 {
		t.Fatalf("expected %s, got %s", largeValue, storedValue)
	}
}

func TestSstoreMultipleKeys(t *testing.T) {
	// Test: SSTORE with multiple different keys
	code := []byte{
		// Store at key 1
		vm.PUSH1, 0xAA, // PUSH1 0xAA
		vm.PUSH1, 0x01, // PUSH1 1
		vm.SSTORE, // SSTORE

		// Store at key 2
		vm.PUSH1, 0xBB, // PUSH1 0xBB
		vm.PUSH1, 0x02, // PUSH1 2
		vm.SSTORE, // SSTORE

		// Store at key 3
		vm.PUSH1, 0xCC, // PUSH1 0xCC
		vm.PUSH1, 0x03, // PUSH1 3
		vm.SSTORE, // SSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check all three storage slots
	testCases := []struct {
		key      string
		expected uint8
	}{
		{"0000000000000000000000000000000000000000000000000000000000000001", 0xAA},
		{"0000000000000000000000000000000000000000000000000000000000000002", 0xBB},
		{"0000000000000000000000000000000000000000000000000000000000000003", 0xCC},
	}

	for _, tc := range testCases {
		storedValue, exists := d.Storage[tc.key]
		if !exists {
			t.Fatalf("storage key %s not found", tc.key)
		}

		expected := uint256.NewInt(uint64(tc.expected))
		if storedValue.Cmp(expected) != 0 {
			t.Fatalf("key %s: expected %s, got %s", tc.key, expected, storedValue)
		}
	}
}

func TestSstoreZeroKey(t *testing.T) {
	// Test: SSTORE with key = 0
	code := []byte{
		vm.PUSH1, 0x77, // PUSH1 0x77 (value)
		vm.PUSH1, 0x00, // PUSH1 0 (key = 0)
		vm.SSTORE, // SSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check storage at key 0
	key := "0000000000000000000000000000000000000000000000000000000000000000"
	storedValue, exists := d.Storage[key]
	if !exists {
		t.Fatalf("storage key %s not found", key)
	}

	expected := uint256.NewInt(0x77)
	if storedValue.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, storedValue)
	}
}

func TestSstoreInStaticCall(t *testing.T) {
	// Test that SSTORE fails in a static call context
	code := []byte{
		vm.PUSH1, 0x42, // value
		vm.PUSH1, 0x00, // slot
		vm.SSTORE, // SSTORE (should fail in static context)
	}

	d := vm.NewDebuggerVM(code, GetHandler)

	// Create a static frame (simulating STATICCALL context)
	staticFrame := vm.MessageFrame{
		Code:         code,
		PC:           0,
		Stack:        vm.NewStack(),
		Memory:       vm.NewMemory(),
		ReturnData:   nil,
		Gas:          1000,
		CallType:     vm.CallTypeStaticCall,
		IsStatic:     true, // This is the key flag
		CodeMetadata: vm.ScanCodeMetadata(code),
	}

	// Push the static frame
	err := d.PushFrame(staticFrame)
	if err != nil {
		t.Fatalf("failed to push static frame: %v", err)
	}

	// Execute until we hit SSTORE
	for !d.Stopped && d.PC() < uint64(len(d.Code())) {
		op := d.Code()[d.PC()]
		if op == vm.SSTORE {
			// This should fail with static call error
			err := d.Step()
			if err != vm.ErrStaticCallStateChange {
				t.Fatalf("expected ErrStaticCallStateChange, got: %v", err)
			}
			return
		}
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error before SSTORE: %v", err)
		}
	}

	t.Fatalf("SSTORE was not executed")
}

func TestSstoreAfterStaticCallReturns(t *testing.T) {
	// Test that SSTORE works after returning from a static call
	code := []byte{
		vm.PUSH1, 0x99, // value
		vm.PUSH1, 0x01, // slot
		vm.SSTORE, // SSTORE (should work after static call returns)
	}

	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	// Simulate being in a static call and then returning
	staticFrame := vm.MessageFrame{
		Code:         []byte{vm.RETURN}, // Simple return
		PC:           0,
		Stack:        vm.NewStack(),
		Memory:       vm.NewMemory(),
		ReturnData:   nil,
		Gas:          1000,
		CallType:     vm.CallTypeStaticCall,
		IsStatic:     true,
		CodeMetadata: vm.ScanCodeMetadata([]byte{vm.RETURN}),
	}

	// Push static frame and then pop it (simulating return from static call)
	err := d.PushFrame(staticFrame)
	if err != nil {
		t.Fatalf("failed to push static frame: %v", err)
	}

	err = d.PopFrame()
	if err != nil {
		t.Fatalf("failed to pop static frame: %v", err)
	}

	// Now we should be back in normal context, SSTORE should work
	for !d.Stopped && d.PC() < uint64(len(d.Code())) {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Verify that the storage was written
	key := "0000000000000000000000000000000000000000000000000000000000000001"
	storedValue, exists := d.Storage[key]
	if !exists {
		t.Fatalf("storage key %s not found", key)
	}

	expected := uint256.NewInt(0x99)
	if storedValue.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, storedValue)
	}
}
