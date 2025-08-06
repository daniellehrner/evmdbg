package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestMsizeInitial(t *testing.T) {
	// Test: MSIZE on fresh VM should return 0
	code := []byte{
		vm.MSIZE, // MSIZE
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
		t.Fatalf("expected 0 for initial memory size, got %s", result)
	}
}

func TestMsizeAfterMstore(t *testing.T) {
	// Test: MSIZE after MSTORE should show memory expansion
	code := []byte{
		vm.PUSH1, 0x42, // PUSH1 0x42 (value)
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.MSTORE, // MSTORE
		vm.MSIZE,  // MSIZE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(32) // MSTORE at offset 0 allocates 32 bytes
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestMsizeAfterMload(t *testing.T) {
	// Test: MSIZE after MLOAD should show memory expansion
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.MLOAD, // MLOAD (reading from uninitialized memory)
		vm.MSIZE, // MSIZE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Pop MSIZE first (it was executed last, so it's on top)
	result, _ := d.Stack().Pop()

	// The loaded value is still on the stack
	loadedValue, _ := d.Stack().Pop()
	_ = loadedValue                // Should be 0
	expected := uint256.NewInt(32) // MLOAD at offset 0 allocates 32 bytes
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestMsizeHighOffset(t *testing.T) {
	// Test: MSIZE after accessing high memory offset
	code := []byte{
		vm.PUSH1, 0xFF, // PUSH1 0xFF (value)
		vm.PUSH2, 0x01, 0x00, // PUSH2 256 (offset)
		vm.MSTORE, // MSTORE
		vm.MSIZE,  // MSIZE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// MSTORE at offset 256 needs to store 32 bytes, so memory size = 256 + 32 = 288
	expected := uint256.NewInt(288)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestMsizeProgressive(t *testing.T) {
	// Test: MSIZE increases as we access higher offsets
	code := []byte{
		// Check initial size
		vm.MSIZE, // MSIZE (should be 0)

		// Store at offset 0
		vm.PUSH1, 0x11, // PUSH1 0x11
		vm.PUSH1, 0x00, // PUSH1 0
		vm.MSTORE, // MSTORE
		vm.MSIZE,  // MSIZE (should be 32)

		// Store at offset 64
		vm.PUSH1, 0x22, // PUSH1 0x22
		vm.PUSH1, 64, // PUSH1 64
		vm.MSTORE, // MSTORE
		vm.MSIZE,  // MSIZE (should be 96 = 64 + 32)
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Stack should have: [0, 32, 96] (most recent on top)
	if d.Stack().Len() != 3 {
		t.Fatalf("expected 3 items on stack, got %d", d.Stack().Len())
	}

	// Check final size (96)
	result3, _ := d.Stack().Pop()
	expected3 := uint256.NewInt(96)
	if result3.Cmp(expected3) != 0 {
		t.Fatalf("expected final size %s, got %s", expected3, result3)
	}

	// Check middle size (32)
	result2, _ := d.Stack().Pop()
	expected2 := uint256.NewInt(32)
	if result2.Cmp(expected2) != 0 {
		t.Fatalf("expected middle size %s, got %s", expected2, result2)
	}

	// Check initial size (0)
	result1, _ := d.Stack().Pop()
	if !result1.IsZero() {
		t.Fatalf("expected initial size 0, got %s", result1)
	}
}

func TestMsizeUnalignedAccess(t *testing.T) {
	// Test: MSIZE after unaligned memory access
	code := []byte{
		vm.PUSH1, 0xFF, // PUSH1 0xFF (value)
		vm.PUSH1, 0x05, // PUSH1 5 (unaligned offset)
		vm.MSTORE, // MSTORE
		vm.MSIZE,  // MSIZE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// MSTORE at offset 5 writes 32 bytes (offsets 5-36), memory size = word-aligned(37) = 64
	expected := uint256.NewInt(64)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestMsizeAfterMultipleOperations(t *testing.T) {
	// Test: MSIZE remains consistent after multiple memory operations
	code := []byte{
		// Multiple MSTORE operations
		vm.PUSH1, 0x11, // PUSH1 0x11
		vm.PUSH1, 0x20, // PUSH1 32
		vm.MSTORE, // MSTORE at 32

		vm.PUSH1, 0x22, // PUSH1 0x22
		vm.PUSH1, 0x40, // PUSH1 64
		vm.MSTORE, // MSTORE at 64

		// MLOAD operations (shouldn't change size)
		vm.PUSH1, 0x20, // PUSH1 32
		vm.MLOAD, // MLOAD from 32
		vm.POP,   // POP (discard loaded value)

		vm.PUSH1, 0x40, // PUSH1 64
		vm.MLOAD, // MLOAD from 64
		vm.POP,   // POP (discard loaded value)

		vm.MSIZE, // MSIZE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	// Highest access was MSTORE at offset 64, requiring 32 bytes
	// So memory size = 64 + 32 = 96
	expected := uint256.NewInt(96)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
