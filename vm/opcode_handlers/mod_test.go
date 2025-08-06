package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestModBasic(t *testing.T) {
	// Test: 7 % 3 = 1
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (divisor)
		vm.PUSH1, 0x07, // PUSH1 7 (dividend)
		vm.MOD, // MOD
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestModByZero(t *testing.T) {
	// Test: 7 % 0 = 0 (EVM returns 0 for modulo by zero)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (divisor)
		vm.PUSH1, 0x07, // PUSH1 7 (dividend)
		vm.MOD, // MOD
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0 for modulo by zero, got %s", result)
	}
}

func TestModZeroDividend(t *testing.T) {
	// Test: 0 % 5 = 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (divisor)
		vm.PUSH1, 0x00, // PUSH1 0 (dividend)
		vm.MOD, // MOD
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0, got %s", result)
	}
}

func TestModExactDivision(t *testing.T) {
	// Test: 9 % 3 = 0 (exact division)
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (divisor)
		vm.PUSH1, 0x09, // PUSH1 9 (dividend)
		vm.MOD, // MOD
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0, got %s", result)
	}
}

func TestModLargerDivisor(t *testing.T) {
	// Test: 5 % 10 = 5 (dividend smaller than divisor)
	code := []byte{
		vm.PUSH1, 0x0A, // PUSH1 10 (divisor)
		vm.PUSH1, 0x05, // PUSH1 5 (dividend)
		vm.MOD, // MOD
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(5)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestModLargeNumbers(t *testing.T) {
	// Test: (2^128 + 5) % 2^64 = 5
	dividend := new(uint256.Int).Add(new(uint256.Int).Lsh(uint256.NewInt(1), 128), uint256.NewInt(5))
	divisor := new(uint256.Int).Lsh(uint256.NewInt(1), 64)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(divisor)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(dividend)...)
	code = append(code, vm.MOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(5)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestModSelfIsZero(t *testing.T) {
	// Test: a % a = 0 (for any non-zero a)
	testValue := uint256.NewInt(12345)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(testValue)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(testValue)...)
	code = append(code, vm.MOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0, got %s", result)
	}
}
