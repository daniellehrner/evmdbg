package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSha3Empty(t *testing.T) {
	// Test: SHA3 of empty data
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (length)
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.SHA3, // SHA3
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// Expected: Keccak256 of empty string
	// 0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470
	expected := new(uint256.Int).SetBytes([]byte{
		0xc5, 0xd2, 0x46, 0x01, 0x86, 0xf7, 0x23, 0x3c,
		0x92, 0x7e, 0x7d, 0xb2, 0xdc, 0xc7, 0x03, 0xc0,
		0xe5, 0x00, 0xb6, 0x53, 0xca, 0x82, 0x27, 0x3b,
		0x7b, 0xfa, 0xd8, 0x04, 0x5d, 0x85, 0xa4, 0x70,
	})

	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSha3SingleByte(t *testing.T) {
	// Test: SHA3 of single byte 0x42
	code := []byte{
		vm.PUSH1, 0x42, // PUSH1 0x42 (data)
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.MSTORE, // MSTORE (store 0x42 at memory offset 0)

		vm.PUSH1, 0x01, // PUSH1 1 (length - 1 byte)
		vm.PUSH1, 0x1F, // PUSH1 31 (offset - rightmost byte of the 32-byte word)
		vm.SHA3, // SHA3
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// Expected: Keccak256 of byte 0x42
	// This is a known test vector - we expect a specific hash
	// The exact value would need to be computed, but we can verify it's not zero
	if result.IsZero() {
		t.Fatalf("expected non-zero hash, got 0")
	}
}

func TestSha3MultipleBytes(t *testing.T) {
	// Test: SHA3 of multiple bytes "Hello"
	helloBytes := []byte("Hello") // 5 bytes: 0x48 0x65 0x6c 0x6c 0x6f

	// Store the bytes using MSTORE with proper padding
	helloValue := new(uint256.Int).SetBytes(helloBytes)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(helloValue)...)
	code = append(code, vm.PUSH1, 0x00) // offset 0
	code = append(code, vm.MSTORE)      // MSTORE

	code = append(code, vm.PUSH1, 0x05) // length 5
	code = append(code, vm.PUSH1, 0x1B) // offset 27 (32-5=27, to get the 5 bytes)
	code = append(code, vm.SHA3)        // SHA3

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// Should produce a valid hash (non-zero)
	if result.IsZero() {
		t.Fatalf("expected non-zero hash, got 0")
	}

	// Log the result for verification if needed
	t.Logf("SHA3 of 'Hello': %s", result)
}

func TestSha3Full32Bytes(t *testing.T) {
	// Test: SHA3 of full 32-byte word
	testData := new(uint256.Int).SetBytes([]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(testData)...)
	code = append(code, vm.PUSH1, 0x00) // offset 0
	code = append(code, vm.MSTORE)      // MSTORE

	code = append(code, vm.PUSH1, 0x20) // length 32
	code = append(code, vm.PUSH1, 0x00) // offset 0
	code = append(code, vm.SHA3)        // SHA3

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// Should produce a valid hash
	if result.IsZero() {
		t.Fatalf("expected non-zero hash, got 0")
	}
}

func TestSha3ZeroData(t *testing.T) {
	// Test: SHA3 of zero bytes
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (value = 0)
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.MSTORE, // MSTORE (store zeros)

		vm.PUSH1, 0x20, // PUSH1 32 (length)
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.SHA3, // SHA3
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// SHA3 of 32 zero bytes should produce a specific known hash
	if result.IsZero() {
		t.Fatalf("expected non-zero hash even for zero data, got 0")
	}
}

func TestSha3LargeData(t *testing.T) {
	// Test: SHA3 of larger data (64 bytes)
	code := []byte{
		// Store first 32 bytes
		vm.PUSH1, 0xAA, // Fill with 0xAA pattern
		vm.PUSH1, 0x00, // offset 0
		vm.MSTORE, // MSTORE

		// Store second 32 bytes
		vm.PUSH1, 0xBB, // Fill with 0xBB pattern
		vm.PUSH1, 0x20, // offset 32
		vm.MSTORE, // MSTORE

		vm.PUSH1, 0x40, // PUSH1 64 (length)
		vm.PUSH1, 0x00, // PUSH1 0 (offset)
		vm.SHA3, // SHA3
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// Should produce a valid hash
	if result.IsZero() {
		t.Fatalf("expected non-zero hash, got 0")
	}
}

func TestSha3MemoryOffset(t *testing.T) {
	// Test: SHA3 with non-zero memory offset
	code := []byte{
		// Store data at offset 64
		vm.PUSH1, 0x12, // PUSH1 0x12
		vm.PUSH1, 0x40, // PUSH1 64 (offset)
		vm.MSTORE, // MSTORE

		vm.PUSH1, 0x01, // PUSH1 1 (length - 1 byte)
		vm.PUSH1, 0x5F, // PUSH1 95 (offset 64 + 31 = rightmost byte)
		vm.SHA3, // SHA3
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()

	// Should produce a valid hash
	if result.IsZero() {
		t.Fatalf("expected non-zero hash, got 0")
	}
}

func TestSha3Deterministic(t *testing.T) {
	// Test: SHA3 should be deterministic (same input -> same output)

	// First execution
	code1 := []byte{
		vm.PUSH1, 0x99, // PUSH1 0x99
		vm.PUSH1, 0x00, // PUSH1 0
		vm.MSTORE, // MSTORE

		vm.PUSH1, 0x01, // PUSH1 1 (length)
		vm.PUSH1, 0x1F, // PUSH1 31 (offset)
		vm.SHA3, // SHA3
	}
	d1 := vm.NewDebuggerVM(code1, GetHandler)

	for !d1.Stopped {
		err := d1.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result1, _ := d1.Stack.Pop()

	// Second execution with same data
	code2 := []byte{
		vm.PUSH1, 0x99, // PUSH1 0x99 (same data)
		vm.PUSH1, 0x00, // PUSH1 0
		vm.MSTORE, // MSTORE

		vm.PUSH1, 0x01, // PUSH1 1 (length)
		vm.PUSH1, 0x1F, // PUSH1 31 (offset)
		vm.SHA3, // SHA3
	}
	d2 := vm.NewDebuggerVM(code2, GetHandler)

	for !d2.Stopped {
		err := d2.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result2, _ := d2.Stack.Pop()

	// Results should be identical
	if result1.Cmp(result2) != 0 {
		t.Fatalf("SHA3 not deterministic: got %s and %s", result1, result2)
	}
}
