package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestShrBasic(t *testing.T) {
	// Test: 8 >> 1 = 4 (logical right shift)
	code := []byte{
		vm.PUSH1, 0x08, // PUSH1 8 (value to shift)
		vm.PUSH1, 0x01, // PUSH1 1 (shift amount)
		vm.SHR, // SHR
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack().Len() != 1 {
		t.Fatalf("expected 1 item on the stack, got %d", d.Stack().Len())
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(4)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShrLargeValue(t *testing.T) {
	// Test: 0xFF >> 4 = 0x0F (logical right shift)
	code := []byte{
		vm.PUSH1, 0xFF, // PUSH1 255
		vm.PUSH1, 0x04, // PUSH1 4 (shift amount)
		vm.SHR, // SHR
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(0x0F) // 255 >> 4 = 15
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShrNegativeValue(t *testing.T) {
	// Test: -1 >> 1 (logical right shift treats as unsigned)
	// -1 in 256-bit = 0xFFFF...FFFF, so -1 >> 1 should be 0x7FFF...FFFF
	negOne := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negOne)...)
	code = append(code, vm.PUSH1, 0x01) // shift by 1
	code = append(code, vm.SHR)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()

	// Expected: 0x7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
	// This is the largest positive 256-bit number (MSB = 0)
	expected := new(uint256.Int).Rsh(new(uint256.Int).SetAllOne(), 1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShrZero(t *testing.T) {
	// Test: 0 >> 5 = 0
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0x05, // PUSH1 5
		vm.SHR, // SHR
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

func TestShrLargeShift(t *testing.T) {
	// Test: any value >> 256 = 0
	code := []byte{
		vm.PUSH1, 0xFF, // PUSH1 255 (any non-zero value)
		vm.PUSH2, 0x01, 0x00, // PUSH2 256 (shift amount >= 256)
		vm.SHR, // SHR
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

func TestShrOverflowShift(t *testing.T) {
	// Test with very large shift value that overflows uint64
	largeShift := new(uint256.Int).SetAllOne() // Max uint256 value

	code := []byte{vm.PUSH1, 0x42} // PUSH1 66 (any value)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeShift)...)
	code = append(code, vm.SHR)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0 for overflow shift, got %s", result)
	}
}

func TestShrShiftBy255(t *testing.T) {
	// Test: MSB >> 255 should give 1
	msb := new(uint256.Int).Lsh(uint256.NewInt(1), 255) // 2^255

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(msb)...)
	code = append(code, vm.PUSH1, 255) // shift by 255
	code = append(code, vm.SHR)

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

func TestShrVsSar(t *testing.T) {
	// Test to show difference between SHR (logical) and SAR (arithmetic)
	// Using -8 (0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF8)
	negEight := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(8))

	// Test SHR: -8 >> 1 (logical)
	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negEight)...)
	code = append(code, vm.PUSH1, 0x01) // shift by 1
	code = append(code, vm.SHR)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	shrResult, _ := d.Stack().Pop()

	// SHR result should be a large positive number (MSB = 0)
	if shrResult.Sign() < 0 {
		t.Fatalf("SHR of negative number should produce positive result, got %s", shrResult)
	}

	// The result should be 0x7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFC
	// which is (2^256 - 8) >> 1 = (2^256 - 8) / 2 = 2^255 - 4
	expected := new(uint256.Int).Sub(new(uint256.Int).Lsh(uint256.NewInt(1), 255), uint256.NewInt(4))
	if shrResult.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, shrResult)
	}
}
