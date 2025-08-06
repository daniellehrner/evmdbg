package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSubBasic(t *testing.T) {
	// Test: 10 - 3 = 7
	code := []byte{
		vm.PUSH1, 0x03, // PUSH1 3
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.SUB, // SUB
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

func TestSubZero(t *testing.T) {
	// Test: 42 - 0 = 42
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.SUB, // SUB
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

func TestSubSelf(t *testing.T) {
	// Test: a - a = 0
	code := []byte{
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.SUB, // SUB
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

func TestSubUnderflow(t *testing.T) {
	// Test: 3 - 10 = 2^256 - 7 (underflow wraps around)
	code := []byte{
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.PUSH1, 0x03, // PUSH1 3
		vm.SUB, // SUB
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// 3 - 10 = -7, which wraps to 2^256 - 7
	expected := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(7))
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSubZeroFromZero(t *testing.T) {
	// Test: 0 - 0 = 0
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0x00, // PUSH1 0
		vm.SUB, // SUB
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

func TestSubLargeNumbers(t *testing.T) {
	// Test: 2^128 - 1 = 2^128 - 1
	largeNum := new(uint256.Int).Lsh(uint256.NewInt(1), 128) // 2^128
	one := uint256.NewInt(1)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(one)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeNum)...)
	code = append(code, vm.SUB)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := new(uint256.Int).Sub(largeNum, one) // 2^128 - 1
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSubMaxFromOne(t *testing.T) {
	// Test: 1 - max_uint256 = 2 (wrap around)
	maxVal := new(uint256.Int).SetAllOne()

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(maxVal)...)
	code = append(code, vm.PUSH1, 0x01) // PUSH1 1
	code = append(code, vm.SUB)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(2) // 1 - (-1) = 2
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
