package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestEqEqual(t *testing.T) {
	// Test: 42 == 42 should return 1
	code := []byte{
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.EQ, // EQ
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

func TestEqNotEqual(t *testing.T) {
	// Test: 42 == 43 should return 0
	code := []byte{
		vm.PUSH1, 0x2B, // PUSH1 43
		vm.PUSH1, 0x2A, // PUSH1 42
		vm.EQ, // EQ
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

func TestEqZero(t *testing.T) {
	// Test: 0 == 0 should return 1
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0x00, // PUSH1 0
		vm.EQ, // EQ
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

func TestEqLargeNumbers(t *testing.T) {
	// Test equality with large 256-bit numbers
	largeNum := new(uint256.Int).SetBytes([]byte{
		0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11,
		0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeNum)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeNum)...)
	code = append(code, vm.EQ)

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

func TestEqLargeNumbersDifferent(t *testing.T) {
	// Test inequality with large numbers differing by 1
	largeNum1 := new(uint256.Int).SetBytes([]byte{
		0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11,
		0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
	})
	largeNum2 := new(uint256.Int).SetBytes([]byte{
		0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE, 0xF0,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11,
		0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x9A, // Different last byte
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeNum2)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeNum1)...)
	code = append(code, vm.EQ)

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
