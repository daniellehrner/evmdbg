package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestPush0(t *testing.T) {
	// Test: PUSH0 should push 0 onto the stack
	code := []byte{
		vm.PUSH0, // PUSH0
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack.Len() != 1 {
		t.Fatalf("expected 1 item on stack, got %d", d.Stack.Len())
	}

	result, _ := d.Stack.Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0, got %s", result)
	}
}

func TestPush0Multiple(t *testing.T) {
	// Test: Multiple PUSH0 operations
	code := []byte{
		vm.PUSH0, // PUSH0
		vm.PUSH0, // PUSH0
		vm.PUSH0, // PUSH0
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack.Len() != 3 {
		t.Fatalf("expected 3 items on stack, got %d", d.Stack.Len())
	}

	// All should be zero
	for i := 0; i < 3; i++ {
		result, _ := d.Stack.Pop()
		if !result.IsZero() {
			t.Fatalf("expected 0 at position %d, got %s", i, result)
		}
	}
}

func TestPush0WithOtherOps(t *testing.T) {
	// Test: PUSH0 mixed with other operations
	code := []byte{
		vm.PUSH0,       // PUSH0 (0)
		vm.PUSH1, 0x42, // PUSH1 0x42
		vm.ADD, // ADD (0 + 0x42 = 0x42)
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack.Len() != 1 {
		t.Fatalf("expected 1 item on stack, got %d", d.Stack.Len())
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x42)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
