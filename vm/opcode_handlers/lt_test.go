package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestLtTrue(t *testing.T) {
	// Test: 5 < 10 should return 1 (unsigned comparison)
	code := []byte{
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.PUSH1, 0x05, // PUSH1 5
		vm.LT, // LT
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

func TestLtFalse(t *testing.T) {
	// Test: 10 < 5 should return 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.LT, // LT
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

func TestLtEqual(t *testing.T) {
	// Test: 5 < 5 should return 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5
		vm.PUSH1, 0x05, // PUSH1 5
		vm.LT, // LT
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

func TestLtUnsignedBehavior(t *testing.T) {
	// Test unsigned behavior: large number (would be negative if signed) vs small number
	// 0x8000000000000000000000000000000000000000000000000000000000000001 vs 1
	largeNum := new(uint256.Int).Add(new(uint256.Int).Lsh(uint256.NewInt(1), 255), uint256.NewInt(1))
	smallNum := uint256.NewInt(1)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(smallNum)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeNum)...)
	code = append(code, vm.LT)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	// In unsigned comparison: large number > small number, so result should be 0
	if !result.IsZero() {
		t.Fatalf("expected 0 (large > small in unsigned), got %s", result)
	}
}

func TestLtZero(t *testing.T) {
	// Test: 0 < 1 should return 1
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1
		vm.PUSH1, 0x00, // PUSH1 0
		vm.LT, // LT
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
