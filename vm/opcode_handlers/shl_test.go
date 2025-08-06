package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestShlBasic(t *testing.T) {
	// Test: 5 << 1 = 10 (basic left shift)
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (value)
		vm.PUSH1, 0x01, // PUSH1 1 (shift amount)
		vm.SHL, // SHL
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(10)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShlPowerOfTwo(t *testing.T) {
	// Test: 1 << 8 = 256 (shifting 1 by 8 positions)
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1 (value)
		vm.PUSH1, 0x08, // PUSH1 8 (shift amount)
		vm.SHL, // SHL
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(256)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShlZeroShift(t *testing.T) {
	// Test: 42 << 0 = 42 (no shift)
	code := []byte{
		vm.PUSH1, 0x2A, // PUSH1 42 (value)
		vm.PUSH1, 0x00, // PUSH1 0 (shift amount)
		vm.SHL, // SHL
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

func TestShlZeroValue(t *testing.T) {
	// Test: 0 << 5 = 0 (shifting zero)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (value)
		vm.PUSH1, 0x05, // PUSH1 5 (shift amount)
		vm.SHL, // SHL
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

func TestShlLargeShift(t *testing.T) {
	// Test: 5 << 256 = 0 (shift amount >= 256 returns 0)
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (value)
		vm.PUSH2, 0x01, 0x00, // PUSH2 256 (shift amount)
		vm.SHL, // SHL
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
		t.Fatalf("expected 0 for large shift, got %s", result)
	}
}

func TestShlOverflowShift(t *testing.T) {
	// Test: 5 << 300 = 0 (shift amount > 256 returns 0)
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (value)
		vm.PUSH2, 0x01, 0x2C, // PUSH2 300 (shift amount)
		vm.SHL, // SHL
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
		t.Fatalf("expected 0 for overflow shift, got %s", result)
	}
}

func TestShlMaxValue(t *testing.T) {
	// Test: maximum value << 1 (should wrap around due to overflow)
	maxVal := new(uint256.Int).SetAllOne()

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(maxVal)...) // PUSH32 max value
	code = append(code, vm.PUSH1, 0x01)              // PUSH1 1 (shift amount)
	code = append(code, vm.SHL)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// Max value << 1 should wrap around due to 256-bit overflow
	// 0xFFFF...FFFF << 1 = 0xFFFF...FFFE (in 256-bit arithmetic)
	maxForCalc := new(uint256.Int).SetAllOne()
	expected := new(uint256.Int).Lsh(maxForCalc, 1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShlDoubling(t *testing.T) {
	// Test: left shift by 1 is equivalent to doubling
	code := []byte{
		vm.PUSH1, 0x0F, // PUSH1 15 (value)
		vm.PUSH1, 0x01, // PUSH1 1 (shift amount)
		vm.SHL, // SHL
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(30) // 15 * 2 = 30
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShlBitPattern(t *testing.T) {
	// Test: 0x01 << 4 = 0x10 (bit pattern verification)
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1 (value)
		vm.PUSH1, 0x04, // PUSH1 4 (shift amount)
		vm.SHL, // SHL
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(0x10)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShlMultipleShifts(t *testing.T) {
	// Test: 3 << 3 = 24 (3 * 2^3 = 3 * 8 = 24)
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (value)
		vm.PUSH1, 0x03, // PUSH1 3 (shift amount)
		vm.SHL, // SHL
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(24)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShlLargeValue(t *testing.T) {
	// Test: large value shifted left
	largeVal := new(uint256.Int).Lsh(uint256.NewInt(1), 100) // 2^100

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeVal)...) // PUSH32 large value
	code = append(code, vm.PUSH1, 0x05)                // PUSH1 5 (shift amount)
	code = append(code, vm.SHL)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// 2^100 << 5 = 2^105
	expected := new(uint256.Int).Lsh(uint256.NewInt(1), 105)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShlOverflow(t *testing.T) {
	// Test: value that causes overflow when shifted
	// 2^255 << 1 should wrap to 0 (overflow in 256-bit arithmetic)
	val := new(uint256.Int).Lsh(uint256.NewInt(1), 255) // 2^255

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(val)...) // PUSH32 2^255
	code = append(code, vm.PUSH1, 0x01)           // PUSH1 1 (shift amount)
	code = append(code, vm.SHL)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// 2^255 << 1 = 2^256, which overflows to 0 in 256-bit arithmetic
	if !result.IsZero() {
		t.Fatalf("expected 0 for overflow, got %s", result)
	}
}

func TestShlEdgeCase255(t *testing.T) {
	// Test: shift by exactly 255 (maximum valid shift)
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1 (value)
		vm.PUSH1, 0xFF, // PUSH1 255 (shift amount)
		vm.SHL, // SHL
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// 1 << 255 = 2^255 (high bit set)
	expected := new(uint256.Int).Lsh(uint256.NewInt(1), 255)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestShlVsSar(t *testing.T) {
	// Test: demonstrate difference between SHL and SAR
	// SHL shifts left (multiplies by 2^n), SAR shifts right with sign extension

	// Test SHL: 8 << 2 = 32
	code := []byte{
		vm.PUSH1, 0x08, // PUSH1 8 (value)
		vm.PUSH1, 0x02, // PUSH1 2 (shift amount)
		vm.SHL, // SHL
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(32) // 8 * 4 = 32
	if result.Cmp(expected) != 0 {
		t.Fatalf("SHL: expected %s, got %s", expected, result)
	}
}
