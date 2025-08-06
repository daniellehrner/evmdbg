package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestMstoreBasic(t *testing.T) {
	// Test: MSTORE(0, 0x1234) - store value at memory offset 0
	code := []byte{
		vm.PUSH2, 0x12, 0x34, // PUSH2 0x1234 (value)
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.MSTORE, // MSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check that memory was written correctly
	// MSTORE stores 32 bytes, with the value right-aligned
	data := d.Memory().Read(0, 32)

	// Expected: 32 bytes with 0x1234 at the end (big-endian)
	expected := make([]byte, 32)
	expected[30] = 0x12
	expected[31] = 0x34

	for i := 0; i < 32; i++ {
		if data[i] != expected[i] {
			t.Fatalf("memory byte %d: expected 0x%02x, got 0x%02x", i, expected[i], data[i])
		}
	}
}

func TestMstoreOffset(t *testing.T) {
	// Test: MSTORE(32, 0xABCD) - store at offset 32
	code := []byte{
		vm.PUSH2, 0xAB, 0xCD, // PUSH2 0xABCD (value)
		vm.PUSH1, 32, // PUSH1 32 (offset)
		vm.MSTORE, // MSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check memory at offset 32
	data := d.Memory().Read(32, 32)

	// Expected: 32 bytes with 0xABCD at the end
	expected := make([]byte, 32)
	expected[30] = 0xAB
	expected[31] = 0xCD

	for i := 0; i < 32; i++ {
		if data[i] != expected[i] {
			t.Fatalf("memory byte %d: expected 0x%02x, got 0x%02x", i, expected[i], data[i])
		}
	}
}

func TestMstoreZero(t *testing.T) {
	// Test: MSTORE(0, 0) - store zero
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (value)
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.MSTORE, // MSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check that 32 bytes of zeros were written
	data := d.Memory().Read(0, 32)
	for i := 0; i < 32; i++ {
		if data[i] != 0 {
			t.Fatalf("memory byte %d: expected 0, got 0x%02x", i, data[i])
		}
	}
}

func TestMstoreLargeValue(t *testing.T) {
	// Test: Store a full 32-byte value
	largeValue := new(uint256.Int).SetBytes([]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeValue)...)
	code = append(code, vm.PUSH1, 0x00) // PUSH1 0 (offset)
	code = append(code, vm.MSTORE)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check that the full value was stored correctly
	data := d.Memory().Read(0, 32)
	expected := largeValue.Bytes32()

	for i := 0; i < 32; i++ {
		if data[i] != expected[i] {
			t.Fatalf("memory byte %d: expected 0x%02x, got 0x%02x", i, expected[i], data[i])
		}
	}
}

func TestMstoreOverlapping(t *testing.T) {
	// Test: Store overlapping values
	code := []byte{
		vm.PUSH2, 0x11, 0x11, // PUSH2 0x1111 (value 1)
		vm.PUSH1, 0x00, // PUSH1 0 (offset 0)
		vm.MSTORE, // MSTORE at offset 0

		vm.PUSH2, 0x22, 0x22, // PUSH2 0x2222 (value 2)
		vm.PUSH1, 0x10, // PUSH1 16 (offset 16 - overlaps)
		vm.MSTORE, // MSTORE at offset 16
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check non-overlapping part of first storage (bytes 0-15)
	data1 := d.Memory().Read(0, 16)
	// Should be all zeros (high-order bytes of 0x1111)
	for i, b := range data1 {
		if b != 0 {
			t.Fatalf("unexpected non-zero byte at position %d: got 0x%02x", i, b)
		}
	}

	// Check second storage (bytes 16-47)
	data2 := d.Memory().Read(16, 32)
	// Should have 0x2222 at the end (bytes 30-31 of this 32-byte chunk, which are bytes 46-47 globally)
	if data2[30] != 0x22 || data2[31] != 0x22 {
		t.Fatalf("second value not stored correctly: got 0x%02x%02x", data2[30], data2[31])
	}

	// Verify overlapping behavior: bytes 16-31 should now be zeros (overwritten by second MSTORE)
	overlapData := d.Memory().Read(16, 16)
	for i, b := range overlapData {
		if b != 0 {
			t.Fatalf("overlapping area not correctly overwritten at byte %d: got 0x%02x", 16+i, b)
		}
	}
}

func TestMstoreMemoryExpansion(t *testing.T) {
	// Test: Store at a high offset to trigger memory expansion
	code := []byte{
		vm.PUSH1, 0xFF, // PUSH1 0xFF (value)
		vm.PUSH2, 0x01, 0x00, // PUSH2 256 (offset)
		vm.MSTORE, // MSTORE
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Memory should have expanded to accommodate offset 256 + 32 bytes
	// Check that the value was stored at the correct location
	data := d.Memory().Read(256, 32)

	// Expected: 32 bytes with 0xFF at the end
	if data[31] != 0xFF {
		t.Fatalf("expected 0xFF at end of stored data, got 0x%02x", data[31])
	}

	// Check that earlier bytes are zero (memory expansion fills with zeros)
	for i := 0; i < 31; i++ {
		if data[i] != 0 {
			t.Fatalf("expected 0 at byte %d, got 0x%02x", i, data[i])
		}
	}
}
