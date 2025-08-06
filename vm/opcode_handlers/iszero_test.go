package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestIsZeroTrue(t *testing.T) {
	// Test: ISZERO(0) should return 1
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.ISZERO, // ISZERO
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

func TestIsZeroFalse(t *testing.T) {
	// Test: ISZERO(1) should return 0
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1
		vm.ISZERO, // ISZERO
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

func TestIsZeroLargeNumber(t *testing.T) {
	// Test: ISZERO(large_number) should return 0
	largeNum := new(uint256.Int).SetAllOne()

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeNum)...)
	code = append(code, vm.ISZERO)

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

func TestIsZeroNegativeOne(t *testing.T) {
	// Test: ISZERO(-1) should return 0 (since -1 is 0xFFFF...FFFF)
	negOne := new(uint256.Int).Sub(new(uint256.Int), uint256.NewInt(1))

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(negOne)...)
	code = append(code, vm.ISZERO)

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
