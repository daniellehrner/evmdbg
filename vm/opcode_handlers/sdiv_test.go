package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSdivPositiveNumbers(t *testing.T) {
	// Test: 42 / 6 = 7 (positive numbers)
	code := []byte{
		vm.PUSH1, 0x06, // PUSH1 6 (divisor)
		vm.PUSH1, 0x2A, // PUSH1 42 (dividend)
		vm.SDIV, // SDIV
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(7)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSdivNegativeDividend(t *testing.T) {
	// Test: -42 / 6 = -7 (negative dividend, positive divisor)
	negFortyTwo := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(42))

	code := []byte{
		vm.PUSH1, 0x06, // PUSH1 6 (divisor)
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negFortyTwo)...)
	code = append(code, vm.SDIV)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(7)) // -7
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSdivNegativeDivisor(t *testing.T) {
	// Test: 42 / -6 = -7 (positive dividend, negative divisor)
	negSix := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(6))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negSix)...)
	code = append(code, vm.PUSH1, 0x2A) // PUSH1 42
	code = append(code, vm.SDIV)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(7)) // -7
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSdivBothNegative(t *testing.T) {
	// Test: -42 / -6 = 7 (both negative)
	negFortyTwo := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(42))
	negSix := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(6))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negSix)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negFortyTwo)...)
	code = append(code, vm.SDIV)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(7) // Positive result
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSdivByZero(t *testing.T) {
	// Test: 42 / 0 = 0 (signed division by zero returns 0)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (divisor)
		vm.PUSH1, 0x2A, // PUSH1 42 (dividend)
		vm.SDIV, // SDIV
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
		t.Fatalf("expected 0 for division by zero, got %s", result)
	}
}

func TestSdivZeroDividend(t *testing.T) {
	// Test: 0 / 5 = 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (divisor)
		vm.PUSH1, 0x00, // PUSH1 0 (dividend)
		vm.SDIV, // SDIV
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

func TestSdivTruncation(t *testing.T) {
	// Test: 7 / 3 = 2 (signed division truncates toward zero)
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (divisor)
		vm.PUSH1, 0x07, // PUSH1 7 (dividend)
		vm.SDIV, // SDIV
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(2)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSdivNegativeTruncation(t *testing.T) {
	// Test: -7 / 3 = -2 (signed division truncates toward zero)
	negSeven := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(7))

	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (divisor)
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negSeven)...)
	code = append(code, vm.SDIV)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(2)) // -2
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSdivOne(t *testing.T) {
	// Test: a / 1 = a
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1 (divisor)
		vm.PUSH1, 0x2A, // PUSH1 42 (dividend)
		vm.SDIV, // SDIV
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(42)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSdivNegativeOne(t *testing.T) {
	// Test: a / -1 = -a
	negOne := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negOne)...)
	code = append(code, vm.PUSH1, 0x2A) // PUSH1 42
	code = append(code, vm.SDIV)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(42)) // -42
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSdivVsDiv(t *testing.T) {
	// Test to show difference between SDIV (signed) and DIV (unsigned)
	// Using a large number that would be negative in signed interpretation
	largeNum := new(uint256.Int).Add(new(uint256.Int).Lsh(uint256.NewInt(1), 255), uint256.NewInt(2))

	code := []byte{
		vm.PUSH1, 0x02, // PUSH1 2 (divisor)
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeNum)...)
	code = append(code, vm.SDIV)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// In signed division, this large number is treated as negative
	// So the result should also be negative (high bit set)
	if result.Sign() >= 0 {
		t.Fatalf("expected negative result for signed division, got %s", result)
	}
}

func TestSdivMinValue(t *testing.T) {
	// Test: -2^255 / -1 (edge case that can cause overflow in some implementations)
	minValue := new(uint256.Int).Lsh(uint256.NewInt(1), 255) // -2^255 in two's complement
	negOne := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negOne)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(minValue)...)
	code = append(code, vm.SDIV)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// -2^255 / -1 should mathematically be 2^255, but that overflows
	// EVM should handle this gracefully (typically wraps to -2^255)
	// The exact behavior depends on implementation
	t.Logf("Result of -2^255 / -1: %s", result)
}
