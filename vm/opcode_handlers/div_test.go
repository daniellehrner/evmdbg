package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestDivBasic(t *testing.T) {
	// Test: 42 / 6 = 7
	code := []byte{
		vm.PUSH1, 0x06, // PUSH1 6 (divisor)
		vm.PUSH1, 0x2A, // PUSH1 42 (dividend)
		vm.DIV, // DIV
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(7)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestDivByZero(t *testing.T) {
	// Test: 42 / 0 = 0 (EVM returns 0 for division by zero)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (divisor)
		vm.PUSH1, 0x2A, // PUSH1 42 (dividend)
		vm.DIV, // DIV
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
		t.Fatalf("expected 0 for division by zero, got %s", result)
	}
}

func TestDivZeroDividend(t *testing.T) {
	// Test: 0 / 5 = 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (divisor)
		vm.PUSH1, 0x00, // PUSH1 0 (dividend)
		vm.DIV, // DIV
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

func TestDivOne(t *testing.T) {
	// Test: 42 / 1 = 42
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1 (divisor)
		vm.PUSH1, 0x2A, // PUSH1 42 (dividend)
		vm.DIV, // DIV
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(42)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestDivTruncation(t *testing.T) {
	// Test: 7 / 3 = 2 (integer division truncates)
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (divisor)
		vm.PUSH1, 0x07, // PUSH1 7 (dividend)
		vm.DIV, // DIV
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(2)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestDivLargeNumbers(t *testing.T) {
	// Test: 2^255 / 2^128 = 2^127
	dividend := new(uint256.Int).Lsh(uint256.NewInt(1), 255) // 2^255
	divisor := new(uint256.Int).Lsh(uint256.NewInt(1), 128)  // 2^128

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(divisor)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(dividend)...)
	code = append(code, vm.DIV)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := new(uint256.Int).Lsh(uint256.NewInt(1), 127) // 2^127
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestDivSelfIsOne(t *testing.T) {
	// Test: a / a = 1 (for any non-zero a)
	testValue := uint256.NewInt(12345)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(testValue)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(testValue)...)
	code = append(code, vm.DIV)

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
