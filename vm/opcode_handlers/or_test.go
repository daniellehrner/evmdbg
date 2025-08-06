package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestOrBasic(t *testing.T) {
	// Test: 0xF0 | 0x0F = 0xFF
	code := []byte{
		vm.PUSH1, 0x0F, // PUSH1 0x0F
		vm.PUSH1, 0xF0, // PUSH1 0xF0
		vm.OR, // OR
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0xFF)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestOrZero(t *testing.T) {
	// Test: anything | 0 = anything
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0xAB, // PUSH1 0xAB
		vm.OR, // OR
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0xAB)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestOrSelf(t *testing.T) {
	// Test: a | a = a (idempotent)
	code := []byte{
		vm.PUSH1, 0x55, // PUSH1 0x55
		vm.PUSH1, 0x55, // PUSH1 0x55
		vm.OR, // OR
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x55)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
