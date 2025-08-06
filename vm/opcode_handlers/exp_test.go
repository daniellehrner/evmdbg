package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestExpBasic(t *testing.T) {
	// Test: 2^3 = 8
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (exponent)
		vm.PUSH1, 0x02, // PUSH1 2 (base)
		vm.EXP, // EXP
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(8)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestExpPowerOfZero(t *testing.T) {
	// Test: a^0 = 1 (for any a != 0)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (exponent)
		vm.PUSH1, 0x05, // PUSH1 5 (base)
		vm.EXP, // EXP
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

func TestExpZeroBase(t *testing.T) {
	// Test: 0^n = 0 (for n > 0)
	code := []byte{
		vm.PUSH1, 0x05, // PUSH1 5 (exponent)
		vm.PUSH1, 0x00, // PUSH1 0 (base)
		vm.EXP, // EXP
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

func TestExpOne(t *testing.T) {
	// Test: 1^n = 1 (for any n)
	code := []byte{
		vm.PUSH1, 0x64, // PUSH1 100 (exponent)
		vm.PUSH1, 0x01, // PUSH1 1 (base)
		vm.EXP, // EXP
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

func TestExpPowerOfTwo(t *testing.T) {
	// Test: 2^8 = 256
	code := []byte{
		vm.PUSH1, 0x08, // PUSH1 8 (exponent)
		vm.PUSH1, 0x02, // PUSH1 2 (base)
		vm.EXP, // EXP
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

func TestExpLargeBase(t *testing.T) {
	// Test: 10^3 = 1000
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3 (exponent)
		vm.PUSH1, 0x0A, // PUSH1 10 (base)
		vm.EXP, // EXP
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(1000)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestExpOverflow(t *testing.T) {
	// Test: Large exponentiation that causes overflow
	// 2^256 should wrap to 0 (mod 2^256)
	code := []byte{
		vm.PUSH2, 0x01, 0x00, // PUSH2 256 (exponent)
		vm.PUSH1, 0x02, // PUSH1 2 (base)
		vm.EXP, // EXP
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// 2^256 mod 2^256 = 0
	if !result.IsZero() {
		t.Fatalf("expected 0 for 2^256 overflow, got %s", result)
	}
}

func TestExpModularArithmetic(t *testing.T) {
	// Test: EXP with values that test modular arithmetic
	// 3^4 = 81
	code := []byte{
		vm.PUSH1, 0x04, // PUSH1 4 (exponent)
		vm.PUSH1, 0x03, // PUSH1 3 (base)
		vm.EXP, // EXP
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(81)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestExpZeroToZero(t *testing.T) {
	// Test: 0^0 (edge case - typically returns 1 in most implementations)
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (exponent)
		vm.PUSH1, 0x00, // PUSH1 0 (base)
		vm.EXP, // EXP
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// Most implementations define 0^0 = 1
	expected := uint256.NewInt(1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s for 0^0, got %s", expected, result)
	}
}
