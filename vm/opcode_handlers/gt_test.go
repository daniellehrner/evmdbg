package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestGtTrue(t *testing.T) {
	// Test: 10 > 5 should return 1
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.GT, // GT
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestGtFalse(t *testing.T) {
	// Test: 5 > 10 should return 0
	code := []byte{
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.PUSH1, 0x05, // PUSH1 5
		vm.GT, // GT
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

func TestGtEqual(t *testing.T) {
	// Test: 5 > 5 should return 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5
		vm.PUSH1, 0x05, // PUSH1 5
		vm.GT, // GT
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

func TestGtMaxValue(t *testing.T) {
	// Test: max_uint256 > 0 should return 1
	maxValue := new(uint256.Int).SetAllOne()

	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(maxValue)...)
	code = append(code, vm.GT)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
