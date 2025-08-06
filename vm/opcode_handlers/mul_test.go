package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestMulBasic(t *testing.T) {
	// Test: 6 * 7 = 42
	code := []byte{
		vm.PUSH1, 0x07, // PUSH1 7
		vm.PUSH1, 0x06, // PUSH1 6
		vm.MUL, // MUL
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

func TestMulZero(t *testing.T) {
	// Test: anything * 0 = 0
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0xFF, // PUSH1 255
		vm.MUL, // MUL
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

func TestMulOne(t *testing.T) {
	// Test: anything * 1 = anything
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.MUL, // MUL
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

func TestMulOverflow(t *testing.T) {
	// Test: multiplication overflow wraps around (mod 2^256)
	// 2^255 * 2 = 0 (overflow)
	largeNum := new(uint256.Int).Lsh(uint256.NewInt(1), 255) // 2^255

	code := []byte{
		vm.PUSH1, 0x02, // PUSH1 2
	}
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeNum)...)
	code = append(code, vm.MUL)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	// 2^255 * 2 = 2^256 = 0 (mod 2^256)
	if !result.IsZero() {
		t.Fatalf("expected 0 (overflow), got %s", result)
	}
}

func TestMulLargeNumbers(t *testing.T) {
	// Test: multiplication of large numbers that don't overflow
	// 2^128 * 2^64 = 2^192
	a := new(uint256.Int).Lsh(uint256.NewInt(1), 128) // 2^128
	b := new(uint256.Int).Lsh(uint256.NewInt(1), 64)  // 2^64

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(b)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(a)...)
	code = append(code, vm.MUL)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).Lsh(uint256.NewInt(1), 192) // 2^192
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestMulMaxValues(t *testing.T) {
	// Test: max_uint256 * max_uint256 (massive overflow)
	maxVal := new(uint256.Int).SetAllOne()

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(maxVal)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(maxVal)...)
	code = append(code, vm.MUL)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	// 0xFFFF...FFFF * 0xFFFF...FFFF = 1 (due to modular arithmetic)
	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
