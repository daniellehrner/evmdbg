package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
)

func TestInvalidOpCode_Execute(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01
		0xfe,       // INVALID
		0x60, 0x02, // PUSH1 0x02 (should never execute)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// First step should succeed (PUSH1)
	err := v.Step()
	if err != nil {
		t.Fatalf("Expected PUSH1 to succeed, got error: %v", err)
	}

	// Verify PUSH1 worked
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack after PUSH1, got %d", v.Stack().Len())
	}

	// Second step should fail on INVALID opcode
	err = v.Step()
	if err == nil {
		t.Fatal("Expected INVALID opcode to cause error, got nil")
	}

	// Verify error message
	if !strings.Contains(err.Error(), "invalid opcode") {
		t.Errorf("Expected error to mention 'invalid opcode', got: %v", err)
	}

	// VM should not be in stopped state (it errored before setting stopped)
	if v.Stopped {
		t.Error("Expected VM not to be in stopped state after error")
	}

	// PC should be at the instruction after INVALID (PC advances before execution)
	if v.PC() != 3 {
		t.Errorf("Expected PC to be at 3 (after INVALID instruction), got %d", v.PC())
	}
}

func TestInvalidOpCode_DoesNotConsumeStack(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01
		0x60, 0x02, // PUSH1 0x02
		0xfe, // INVALID
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Execute the two PUSH instructions
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during PUSH: %v", err)
		}
	}

	// Verify stack has 2 items
	if v.Stack().Len() != 2 {
		t.Fatalf("Expected 2 items on stack, got %d", v.Stack().Len())
	}

	// INVALID should fail without consuming stack
	err := v.Step()
	if err == nil {
		t.Fatal("Expected INVALID opcode to cause error, got nil")
	}

	// Stack should still have 2 items
	if v.Stack().Len() != 2 {
		t.Errorf("Expected stack to remain unchanged after INVALID, got %d items", v.Stack().Len())
	}
}
