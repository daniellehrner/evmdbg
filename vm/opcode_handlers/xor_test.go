package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestXorBasic(t *testing.T) {
	// Test: 0xFF ^ 0x0F = 0xF0
	code := []byte{
		vm.PUSH1, 0x0F, // PUSH1 0x0F
		vm.PUSH1, 0xFF, // PUSH1 0xFF
		vm.XOR, // XOR
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0xF0)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestXorZero(t *testing.T) {
	// Test: anything ^ 0 = anything
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0xAB, // PUSH1 0xAB
		vm.XOR, // XOR
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

func TestXorSelf(t *testing.T) {
	// Test: a ^ a = 0
	code := []byte{
		vm.PUSH1, 0x55, // PUSH1 0x55
		vm.PUSH1, 0x55, // PUSH1 0x55
		vm.XOR, // XOR
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0, got %s", result)
	}
}

func TestXorInverse(t *testing.T) {
	// Test: a ^ 0xFF = ~a (for single byte)
	code := []byte{
		vm.PUSH1, 0xFF, // PUSH1 0xFF
		vm.PUSH1, 0xAA, // PUSH1 0xAA
		vm.XOR, // XOR
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x55) // ~0xAA = 0x55
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
