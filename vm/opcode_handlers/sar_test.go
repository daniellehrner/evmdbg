package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSarBasic(t *testing.T) {
	// Test: 8 >> 1 = 4
	code := []byte{
		vm.PUSH1, 0x08, // PUSH1 8 (value to shift)
		vm.PUSH1, 0x01, // PUSH1 1 (shift amount)
		vm.SAR, // SAR
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

func TestSarNegative(t *testing.T) {
	// Test: -8 >> 1 = -4
	// -8 in two's complement 256-bit is 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF8
	negEight := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(8))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negEight)...)
	code = append(code, vm.PUSH1, 0x01) // shift by 1
	code = append(code, vm.SAR)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(4)) // -4
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSarZero(t *testing.T) {
	// Test: 0 >> 5 = 0
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0x05, // PUSH1 5
		vm.SAR, // SAR
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

func TestSarLargeShift(t *testing.T) {
	// Test: positive value >> 256 = 0
	code := []byte{
		vm.PUSH1, 0x08, // PUSH1 8 (positive value)
		vm.PUSH2, 0x01, 0x00, // PUSH2 256 (shift amount >= 256)
		vm.SAR, // SAR
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

func TestSarLargeShiftNegative(t *testing.T) {
	// Test: negative value >> 256 = -1 (all bits set)
	negOne := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negOne)...)
	code = append(code, vm.PUSH2, 0x01, 0x00) // shift by 256
	code = append(code, vm.SAR)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := new(uint256.Int).SetAllOne() // -1 (all bits set)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSarShiftBy31(t *testing.T) {
	// Test: positive number >> 31 (0x80000000 is positive in 256-bit)
	// 0x80000000 >> 31 = 1 (2^31 >> 31 = 1)
	value := uint256.NewInt(0x80000000)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(value)...)
	code = append(code, vm.PUSH1, 31) // shift by 31
	code = append(code, vm.SAR)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(1) // 2^31 >> 31 = 1
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSarMostSignificantBit(t *testing.T) {
	// Test: MSB set (0x8000...0000) >> 1 should preserve sign
	msb := new(uint256.Int).Lsh(uint256.NewInt(1), 255) // 2^255

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(msb)...)
	code = append(code, vm.PUSH1, 0x01) // shift by 1
	code = append(code, vm.SAR)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()

	// MSB should still be set after right shift
	if result.Sign() >= 0 {
		t.Fatalf("expected negative result (MSB set), got %s", result)
	}
}

func TestSarOverflowShift(t *testing.T) {
	// Test with very large shift value that overflows uint64
	largeShift := new(uint256.Int).SetAllOne() // Max uint256 value

	code := []byte{vm.PUSH1, 0x01} // PUSH1 1 (positive value)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeShift)...)
	code = append(code, vm.SAR)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0 for positive value with overflow shift, got %s", result)
	}
}
