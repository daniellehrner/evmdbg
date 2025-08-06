package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestAndBasic(t *testing.T) {
	// Test: 0xFF & 0x0F = 0x0F
	code := []byte{
		vm.PUSH1, 0x0F, // PUSH1 0x0F
		vm.PUSH1, 0xFF, // PUSH1 0xFF
		vm.AND, // AND
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(0x0F)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestAndZero(t *testing.T) {
	// Test: anything & 0 = 0
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH1, 0xFF, // PUSH1 0xFF
		vm.AND, // AND
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

func TestAndAllOnes(t *testing.T) {
	// Test: anything & 0xFFFFFFFF = anything (for 32-bit values)
	testValue := uint256.NewInt(0x12345678)
	allOnes := uint256.NewInt(0xFFFFFFFF)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(allOnes)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(testValue)...)
	code = append(code, vm.AND)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if result.Cmp(testValue) != 0 {
		t.Fatalf("expected %s, got %s", testValue, result)
	}
}

func TestAndLargeNumbers(t *testing.T) {
	// Test with large 256-bit numbers
	a := new(uint256.Int).SetBytes([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00,
		0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00,
	})
	b := new(uint256.Int).SetBytes([]byte{
		0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
		0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
		0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
		0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(b)...)
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(a)...)
	code = append(code, vm.AND)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := new(uint256.Int).SetBytes([]byte{
		0xF0, 0xF0, 0xF0, 0xF0, 0x00, 0x00, 0x00, 0x00,
		0xF0, 0xF0, 0xF0, 0xF0, 0x00, 0x00, 0x00, 0x00,
		0xF0, 0xF0, 0xF0, 0xF0, 0x00, 0x00, 0x00, 0x00,
		0xF0, 0xF0, 0xF0, 0xF0, 0x00, 0x00, 0x00, 0x00,
	})

	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
