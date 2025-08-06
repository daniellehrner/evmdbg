package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSmodPositiveNumbers(t *testing.T) {
	// Test: 7 % 3 = 1 (positive numbers)
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (b - divisor)
		vm.PUSH1, 0x07, // PUSH1 7 (a - dividend)
		vm.SMOD, // SMOD
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

func TestSmodNegativeDividend(t *testing.T) {
	// Test: -7 % 3 = -1 (negative dividend, positive divisor)
	// Result should have same sign as dividend
	negSeven := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(7))

	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (b - divisor)
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negSeven)...) // a = -7
	code = append(code, vm.SMOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1)) // -1
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSmodPositiveDividendNegativeDivisor(t *testing.T) {
	// Test: 7 % -3 = 1 (positive dividend, negative divisor)
	// Result should have same sign as dividend (positive)
	negThree := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(3))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negThree)...) // b = -3
	code = append(code, vm.PUSH1, 0x07)                // a = 7
	code = append(code, vm.SMOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(1) // Positive result (sign follows dividend)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSmodNegativeBoth(t *testing.T) {
	// Test: -7 % -3 = -1 (both negative)
	// Result should have same sign as dividend (negative)
	negSeven := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(7))
	negThree := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(3))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negThree)...) // b = -3
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negSeven)...) // a = -7
	code = append(code, vm.SMOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1)) // -1
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSmodDivisionByZero(t *testing.T) {
	// Test: 7 % 0 = 0 (division by zero should return 0)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (b - divisor)
		vm.PUSH1, 0x07, // PUSH1 7 (a - dividend)
		vm.SMOD, // SMOD
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

func TestSmodZeroDividend(t *testing.T) {
	// Test: 0 % 5 = 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (b - divisor)
		vm.PUSH1, 0x00, // PUSH1 0 (a - dividend)
		vm.SMOD, // SMOD
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

func TestSmodExactDivision(t *testing.T) {
	// Test: 9 % 3 = 0 (exact division)
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (b - divisor)
		vm.PUSH1, 0x09, // PUSH1 9 (a - dividend)
		vm.SMOD, // SMOD
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

func TestSmodNegativeExactDivision(t *testing.T) {
	// Test: -9 % 3 = 0 (exact division with negative dividend)
	negNine := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(9))

	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (b - divisor)
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negNine)...) // a = -9
	code = append(code, vm.SMOD)

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

func TestSmodLargeNumbers(t *testing.T) {
	// Test with large numbers
	// (2^255 - 1) % 2^254 = 2^254 - 1
	dividend := new(uint256.Int).Sub(new(uint256.Int).Lsh(uint256.NewInt(1), 255), uint256.NewInt(1)) // 2^255 - 1
	divisor := new(uint256.Int).Lsh(uint256.NewInt(1), 254)                                           // 2^254

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(divisor)...) // b = 2^254
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(dividend)...) // a = 2^255 - 1
	code = append(code, vm.SMOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int).Lsh(uint256.NewInt(1), 254), uint256.NewInt(1)) // 2^254 - 1
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSmodSignBehavior(t *testing.T) {
	// Test that result sign follows dividend (not divisor)
	// This is the key difference from regular MOD

	// Case 1: -10 % 3 = -1 (dividend negative → result negative)
	negTen := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(10))

	code1 := []byte{
		vm.PUSH1, 0x03, // PUSH1 3
	}
	code1 = append(code1, vm.PUSH32)
	code1 = append(code1, bytes32WithValue(negTen)...) // a = -10
	code1 = append(code1, vm.SMOD)

	d1 := vm.NewDebuggerVM(code1, GetHandler)

	for !d1.Stopped {
		err := d1.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result1, _ := d1.Stack.Pop()
	expected1 := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1)) // -1
	if result1.Cmp(expected1) != 0 {
		t.Fatalf("Case 1: expected %s, got %s", expected1, result1)
	}

	// Case 2: 10 % -3 = 1 (dividend positive → result positive)
	negThree := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(3))

	code2 := []byte{vm.PUSH32}
	code2 = append(code2, bytes32WithValue(negThree)...) // b = -3
	code2 = append(code2, vm.PUSH1, 0x0A)                // a = 10
	code2 = append(code2, vm.SMOD)

	d2 := vm.NewDebuggerVM(code2, GetHandler)

	for !d2.Stopped {
		err := d2.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result2, _ := d2.Stack.Pop()
	expected2 := uint256.NewInt(1) // 1 (positive)
	if result2.Cmp(expected2) != 0 {
		t.Fatalf("Case 2: expected %s, got %s", expected2, result2)
	}
}

func TestSmodVsMod(t *testing.T) {
	// Test to show difference between SMOD (signed) and MOD (unsigned)
	// Using a case where they would differ

	// For unsigned: large_number % small_number
	// For signed: negative_number % small_number

	// Using 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFE (-2 in signed)
	negTwo := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(2))

	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (divisor)
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negTwo)...) // a = -2
	code = append(code, vm.SMOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(2)) // -2 (since -2 % 5 = -2)
	if result.Cmp(expected) != 0 {
		t.Fatalf("SMOD: expected %s, got %s", expected, result)
	}

	// Note: Regular MOD would treat this as a huge positive number % 5
}
