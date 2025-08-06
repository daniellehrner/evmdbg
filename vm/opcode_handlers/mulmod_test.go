package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestMulmodBasic(t *testing.T) {
	// Test: (6 * 7) % 5 = 42 % 5 = 2
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (modulus)
		vm.PUSH1, 0x07, // PUSH1 7 (b)
		vm.PUSH1, 0x06, // PUSH1 6 (a)
		vm.MULMOD, // MULMOD
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

func TestMulmodZeroModulus(t *testing.T) {
	// Test: (a * b) % 0 = 0 (EVM returns 0 for division by zero)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (modulus)
		vm.PUSH1, 0x07, // PUSH1 7 (b)
		vm.PUSH1, 0x06, // PUSH1 6 (a)
		vm.MULMOD, // MULMOD
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
		t.Fatalf("expected 0 for zero modulus, got %s", result)
	}
}

func TestMulmodZeroOperands(t *testing.T) {
	// Test: (0 * b) % m = 0
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (modulus)
		vm.PUSH1, 0x07, // PUSH1 7 (b)
		vm.PUSH1, 0x00, // PUSH1 0 (a)
		vm.MULMOD, // MULMOD
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

func TestMulmodOne(t *testing.T) {
	// Test: (a * 1) % m = a % m
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (modulus)
		vm.PUSH1, 0x01, // PUSH1 1 (b)
		vm.PUSH1, 0x07, // PUSH1 7 (a)
		vm.MULMOD, // MULMOD
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(2) // 7 % 5 = 2
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestMulmodLargeNumbers(t *testing.T) {
	// Test: MULMOD with large numbers that would overflow in regular multiplication
	// Use values where a * b > 2^256 but (a * b) % m fits in 256 bits
	a := new(uint256.Int).Lsh(uint256.NewInt(1), 200) // 2^200
	b := new(uint256.Int).Lsh(uint256.NewInt(1), 100) // 2^100
	m := new(uint256.Int).Lsh(uint256.NewInt(1), 64)  // 2^64

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(m)...) // modulus
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(b)...) // b
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(a)...) // a
	code = append(code, vm.MULMOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	// 2^200 * 2^100 = 2^300
	// 2^300 % 2^64 = 0 (since 2^300 is divisible by 2^64)
	if !result.IsZero() {
		t.Fatalf("expected 0, got %s", result)
	}
}

func TestMulmodModulusOne(t *testing.T) {
	// Test: (a * b) % 1 = 0 (anything mod 1 is 0)
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1 (modulus)
		vm.PUSH1, 0x07, // PUSH1 7 (b)
		vm.PUSH1, 0x06, // PUSH1 6 (a)
		vm.MULMOD, // MULMOD
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
		t.Fatalf("expected 0 for modulus 1, got %s", result)
	}
}

func TestMulmodExactDivision(t *testing.T) {
	// Test: (a * b) % (a * b) = 0
	code := []byte{
		vm.PUSH1, 0x0C, // PUSH1 12 (modulus = 3 * 4)
		vm.PUSH1, 0x04, // PUSH1 4 (b)
		vm.PUSH1, 0x03, // PUSH1 3 (a)
		vm.MULMOD, // MULMOD: (3 * 4) % 12 = 12 % 12 = 0
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
		t.Fatalf("expected 0 for exact division, got %s", result)
	}
}

func TestMulmodOverflowPrevention(t *testing.T) {
	// Test: MULMOD should prevent overflow that would occur with MUL followed by MOD
	// Use max uint256 values
	maxVal := new(uint256.Int).SetAllOne()
	modulus := uint256.NewInt(1000)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(modulus)...) // modulus
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(maxVal)...) // b
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(maxVal)...) // a
	code = append(code, vm.MULMOD)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// max * max = (2^256 - 1)^2 = 2^512 - 2^257 + 1
	// This mod 1000 should give some non-zero result
	// The exact value is complex to calculate, but it should not be zero
	// and should be less than 1000
	if result.GtUint64(999) {
		t.Fatalf("result should be < 1000, got %s", result)
	}
}

func TestMulmodCommutative(t *testing.T) {
	// Test: (a * b) % m = (b * a) % m (commutative property)

	// First: (3 * 7) % 5
	code1 := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (modulus)
		vm.PUSH1, 0x07, // PUSH1 7 (b)
		vm.PUSH1, 0x03, // PUSH1 3 (a)
		vm.MULMOD, // MULMOD
	}
	d1 := vm.NewDebuggerVM(code1, GetHandler)

	for !d1.Stopped {
		err := d1.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result1, _ := d1.Stack.Pop()

	// Second: (7 * 3) % 5
	code2 := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (modulus)
		vm.PUSH1, 0x03, // PUSH1 3 (b)
		vm.PUSH1, 0x07, // PUSH1 7 (a)
		vm.MULMOD, // MULMOD
	}
	d2 := vm.NewDebuggerVM(code2, GetHandler)

	for !d2.Stopped {
		err := d2.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result2, _ := d2.Stack.Pop()

	// Results should be equal
	if result1.Cmp(result2) != 0 {
		t.Fatalf("commutative property failed: %s != %s", result1, result2)
	}
}
