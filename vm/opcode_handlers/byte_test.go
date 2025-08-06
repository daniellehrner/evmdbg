package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestByteBasic(t *testing.T) {
	// Test: BYTE(31, 0x1234) = 0x34 (get last byte)
	code := []byte{
		vm.PUSH2, 0x12, 0x34, // PUSH2 0x1234
		vm.PUSH1, 31, // PUSH1 31 (get byte at index 31 - rightmost byte)
		vm.BYTE, // BYTE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x34)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestByteFirstByte(t *testing.T) {
	// Test: BYTE(0, 0xABCDEF) = 0x00 (get first byte - most significant)
	code := []byte{
		vm.PUSH4, 0x00, 0xAB, 0xCD, 0xEF, // PUSH4 0x00ABCDEF
		vm.PUSH1, 0, // PUSH1 0 (get byte at index 0 - leftmost byte)
		vm.BYTE, // BYTE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x00)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestByteMiddleByte(t *testing.T) {
	// Test: BYTE(30, 0x1234) = 0x12 (get second-to-last byte)
	code := []byte{
		vm.PUSH2, 0x12, 0x34, // PUSH2 0x1234
		vm.PUSH1, 30, // PUSH1 30 (get byte at index 30)
		vm.BYTE, // BYTE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x12)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestByteOutOfBounds(t *testing.T) {
	// Test: BYTE(32, 0x1234) = 0x00 (index >= 32 returns 0)
	code := []byte{
		vm.PUSH2, 0x12, 0x34, // PUSH2 0x1234
		vm.PUSH1, 32, // PUSH1 32 (out of bounds)
		vm.BYTE, // BYTE
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
		t.Fatalf("expected 0 for out of bounds, got %s", result)
	}
}

func TestByteLargeIndex(t *testing.T) {
	// Test: BYTE(1000, 0x1234) = 0x00 (large index returns 0)
	code := []byte{
		vm.PUSH2, 0x12, 0x34, // PUSH2 0x1234
		vm.PUSH2, 0x03, 0xE8, // PUSH2 1000 (very large index)
		vm.BYTE, // BYTE
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
		t.Fatalf("expected 0 for large index, got %s", result)
	}
}

func TestByteFullWord(t *testing.T) {
	// Test with a full 32-byte word
	fullWord := new(uint256.Int).SetBytes([]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
	})

	// Test BYTE(0, fullWord) = 0x00
	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(fullWord)...)
	code = append(code, vm.PUSH1, 0) // PUSH1 0
	code = append(code, vm.BYTE)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x00)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestByteLastByteFullWord(t *testing.T) {
	// Test BYTE(31, fullWord) = 0x1F (last byte)
	fullWord := new(uint256.Int).SetBytes([]byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(fullWord)...)
	code = append(code, vm.PUSH1, 31) // PUSH1 31
	code = append(code, vm.BYTE)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x1F)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
