package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestPopBasic(t *testing.T) {
	// Test: Push a value then pop it - stack should be empty
	code := []byte{
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.POP, // POP
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack().Len() != 0 {
		t.Fatalf("expected empty stack, got %d items", d.Stack().Len())
	}
}

func TestPopMultiple(t *testing.T) {
	// Test: Push two values, pop one - one should remain
	code := []byte{
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.PUSH1, 0x1E, // PUSH1 30
		vm.POP, // POP (removes 30)
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack().Len() != 1 {
		t.Fatalf("expected 1 item on stack, got %d", d.Stack().Len())
	}

	remaining, _ := d.Stack().Pop()
	expected := uint256.NewInt(42)
	if remaining.Cmp(expected) != 0 {
		t.Fatalf("expected %s to remain, got %s", expected, remaining)
	}
}

func TestPopOrder(t *testing.T) {
	// Test: Push A, Push B, Pop -> B is popped, A remains
	code := []byte{
		vm.PUSH1, 0x11, // PUSH1 0x11 (A)
		vm.PUSH1, 0x22, // PUSH1 0x22 (B)
		vm.PUSH1, 0x33, // PUSH1 0x33 (C)
		vm.POP, // POP (removes C=0x33)
		vm.POP, // POP (removes B=0x22)
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack().Len() != 1 {
		t.Fatalf("expected 1 item on stack, got %d", d.Stack().Len())
	}

	remaining, _ := d.Stack().Pop()
	expected := uint256.NewInt(0x11) // A should remain
	if remaining.Cmp(expected) != 0 {
		t.Fatalf("expected %s to remain, got %s", expected, remaining)
	}
}

func TestPopLargeValue(t *testing.T) {
	// Test: Pop works with large 256-bit values
	largeValue := new(uint256.Int).SetAllOne()

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeValue)...)
	code = append(code, vm.POP) // POP

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack().Len() != 0 {
		t.Fatalf("expected empty stack, got %d items", d.Stack().Len())
	}
}

func TestPopEmptyStackError(t *testing.T) {
	// Test: POP on empty stack should cause an error
	code := []byte{
		vm.POP, // POP with empty stack
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	var executionError error
	for !d.Stopped {
		err := d.Step()
		if err != nil {
			executionError = err
			break
		}
	}

	if executionError == nil {
		t.Fatalf("expected error when popping empty stack, but got none")
	}
}
