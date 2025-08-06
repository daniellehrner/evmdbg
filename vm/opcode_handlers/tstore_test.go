package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestTStoreOpCode_Execute(t *testing.T) {
	code := []byte{
		0x61, 0xde, 0xad, // PUSH2 0xdead (value)
		0x60, 0x42, // PUSH1 0x42 (slot)
		0x5d, // TSTORE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Execute all instructions
	for !v.Stopped {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during execution: %v", err)
		}
	}

	// Check that value was stored in transient storage
	slot := uint256.NewInt(0x42)
	expectedValue := uint256.NewInt(0xdead)

	storedValue := v.ReadTransientStorage(slot)
	if storedValue.Cmp(expectedValue) != 0 {
		t.Errorf("Expected stored value %s, got %s", expectedValue.Hex(), storedValue.Hex())
	}

	// Stack should be empty
	if v.Stack().Len() != 0 {
		t.Errorf("Expected empty stack, got %d items", v.Stack().Len())
	}
}

func TestTStoreOpCode_TLoadIntegration(t *testing.T) {
	// Test TSTORE followed by TLOAD
	code := []byte{
		0x61, 0xca, 0xfe, // PUSH2 0xcafe (value)
		0x60, 0x10, // PUSH1 0x10 (slot)
		0x5d, // TSTORE

		0x60, 0x10, // PUSH1 0x10 (same slot)
		0x5c, // TLOAD
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Execute all instructions
	for !v.Stopped {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during execution: %v", err)
		}
	}

	// Check that TLOAD retrieved the value stored by TSTORE
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	value, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	expectedValue := uint256.NewInt(0xcafe)
	if value.Cmp(expectedValue) != 0 {
		t.Errorf("Expected loaded value %s, got %s", expectedValue.Hex(), value.Hex())
	}
}

func TestTStoreOpCode_OverwriteValue(t *testing.T) {
	// Test overwriting a transient storage slot
	code := []byte{
		// First store
		0x60, 0x01, // PUSH1 0x01 (first value)
		0x60, 0x20, // PUSH1 0x20 (slot)
		0x5d, // TSTORE

		// Overwrite with new value
		0x60, 0x02, // PUSH1 0x02 (second value)
		0x60, 0x20, // PUSH1 0x20 (same slot)
		0x5d, // TSTORE

		// Load the value
		0x60, 0x20, // PUSH1 0x20 (same slot)
		0x5c, // TLOAD
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Execute all instructions
	for !v.Stopped {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during execution: %v", err)
		}
	}

	// Should have the second value (0x02), not the first (0x01)
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	value, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	expectedValue := uint256.NewInt(0x02)
	if value.Cmp(expectedValue) != 0 {
		t.Errorf("Expected overwritten value %s, got %s", expectedValue.Hex(), value.Hex())
	}
}

func TestTStoreOpCode_StaticCallRestriction(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01 (value)
		0x60, 0x00, // PUSH1 0x00 (slot)
		0x5d, // TSTORE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Simulate static call context by setting IsStatic on current frame
	currentFrame := v.CurrentFrame()
	if currentFrame != nil {
		currentFrame.IsStatic = true
	}

	// Execute PUSH instructions
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during PUSH: %v", err)
		}
	}

	// TSTORE should fail in static context
	err := v.Step()
	if err == nil {
		t.Fatal("Expected error in static call context, got nil")
	}

	if err != vm.ErrStaticCallStateChange {
		t.Errorf("Expected ErrStaticCallStateChange, got: %v", err)
	}
}

func TestTStoreOpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01 (only one item on stack)
		0x5d, // TSTORE (needs two items)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Execute PUSH1
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH1: %v", err)
	}

	// TSTORE should fail with stack underflow
	err = v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}
