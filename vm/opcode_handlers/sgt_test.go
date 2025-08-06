package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSgtPositiveNumbers(t *testing.T) {
	// Test: 10 > 5 should return 1 (signed comparison)
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.SGT, // SGT
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

func TestSgtPositiveNumbersFalse(t *testing.T) {
	// Test: 5 > 10 should return 0
	code := []byte{
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.PUSH1, 0x05, // PUSH1 5
		vm.SGT, // SGT
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

func TestSgtEqual(t *testing.T) {
	// Test: 5 > 5 should return 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5
		vm.PUSH1, 0x05, // PUSH1 5
		vm.SGT, // SGT
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

func TestSgtNegativeNumbers(t *testing.T) {
	// Test: -5 > -10 should return 1 (true)
	negFive := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(5))
	negTen := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(10))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negTen)...) // b = -10
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negFive)...) // a = -5
	code = append(code, vm.SGT)

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

func TestSgtNegativeNumbersFalse(t *testing.T) {
	// Test: -10 > -5 should return 0 (false)
	negFive := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(5))
	negTen := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(10))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negFive)...) // b = -5
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negTen)...) // a = -10
	code = append(code, vm.SGT)

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

func TestSgtMixedPositiveNegative(t *testing.T) {
	// Test: 1 > -1 should return 1 (true)
	negOne := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1))
	posOne := uint256.NewInt(1)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negOne)...) // b = -1
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(posOne)...) // a = 1
	code = append(code, vm.SGT)

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

func TestSgtMixedNegativePositive(t *testing.T) {
	// Test: -1 > 1 should return 0 (false)
	negOne := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1))
	posOne := uint256.NewInt(1)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(posOne)...) // b = 1
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(negOne)...) // a = -1
	code = append(code, vm.SGT)

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

func TestSgtZero(t *testing.T) {
	// Test: 1 > 0 should return 1 (true)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0x01, // PUSH1 1
		vm.SGT, // SGT
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

func TestSgtZeroNegative(t *testing.T) {
	// Test: 0 > -1 should return 1 (true)
	negOne := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negOne)...) // b = -1
	code = append(code, vm.PUSH1, 0x00)              // a = 0
	code = append(code, vm.SGT)

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

func TestSgtLargeNumbers(t *testing.T) {
	// Test with large numbers near the boundary
	// 2^255 - 1 (largest positive) vs 2^255 (smallest negative)
	largestPositive := new(uint256.Int).Sub(new(uint256.Int).Lsh(uint256.NewInt(1), 255), uint256.NewInt(1))
	smallestNegative := new(uint256.Int).Lsh(uint256.NewInt(1), 255) // -2^255

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(smallestNegative)...) // b = -2^255
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largestPositive)...) // a = 2^255 - 1
	code = append(code, vm.SGT)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	// 2^255-1 > -2^255 should be true
	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSgtVsGt(t *testing.T) {
	// Test to show difference between SGT (signed) and GT (unsigned)
	// Using values that differ in signed vs unsigned comparison

	// 0x8000000000000000000000000000000000000000000000000000000000000001 (negative in signed)
	// vs 0x1 (positive)
	largeUnsigned := new(uint256.Int).Add(new(uint256.Int).Lsh(uint256.NewInt(1), 255), uint256.NewInt(1))
	small := uint256.NewInt(1)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(small)...) // b = 1
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeUnsigned)...) // a = large negative number
	code = append(code, vm.SGT)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	// In signed comparison: large negative < 1 should be false (0)
	if !result.IsZero() {
		t.Fatalf("SGT: expected 0, got %s", result)
	}

	// Note: In unsigned comparison (GT), this would be true since the large number > 1
}
