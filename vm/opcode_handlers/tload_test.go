package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestTLoadOpCode_Execute(t *testing.T) {
	code := []byte{
		0x60, 0x42, // PUSH1 0x42 (slot)
		0x5c, // TLOAD
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Manually set some value in transient storage for testing
	slot := uint256.NewInt(0x42)
	expectedValue := uint256.NewInt(0xdeadbeef)
	v.WriteTransientStorage(slot, expectedValue)

	// Execute PUSH1
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH1: %v", err)
	}

	// Execute TLOAD
	err = v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during TLOAD: %v", err)
	}

	// Check result
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	value, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if value.Cmp(expectedValue) != 0 {
		t.Errorf("Expected value %s, got %s", expectedValue.Hex(), value.Hex())
	}
}

func TestTLoadOpCode_EmptySlot(t *testing.T) {
	code := []byte{
		0x60, 0x99, // PUSH1 0x99 (unused slot)
		0x5c, // TLOAD
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Execute PUSH1 and TLOAD
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should be 0 for empty slot
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	value, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !value.IsZero() {
		t.Errorf("Expected zero value for empty slot, got %s", value.Hex())
	}
}

func TestTLoadOpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0x5c, // TLOAD (no slot on stack)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// TLOAD should fail with stack underflow
	err := v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}
