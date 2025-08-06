package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
)

func TestExtCodeCopyOpCode_Execute(t *testing.T) {
	// Test address: 0x1234567890123456789012345678901234567890
	testAddr := [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
	testCode := []byte{0x60, 0x01, 0x60, 0x02, 0x01, 0x60, 0x03, 0x60, 0x04} // 8 bytes of bytecode

	code := []byte{
		// Push arguments for EXTCODECOPY: address, destOffset, offset, size
		0x60, 0x04, // PUSH1 0x04 (size - copy 4 bytes)
		0x60, 0x02, // PUSH1 0x02 (offset - start from byte 2)
		0x60, 0x00, // PUSH1 0x00 (destOffset - write to memory offset 0)
		0x73,                                                       // PUSH20 (address)
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, // address bytes
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90,
		0x3c, // EXTCODECOPY
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider
	mock := &mockStateProvider{
		codeMap: map[[20]byte][]byte{
			testAddr: testCode,
		},
	}
	v.StateProvider = mock

	// Execute all instructions
	for !v.Stopped {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during execution: %v", err)
		}
	}

	// Check memory - should contain bytes 2-5 from testCode: [0x60, 0x02, 0x01, 0x60]
	result := v.Memory().Read(0, 4)
	expected := []byte{0x60, 0x02, 0x01, 0x60}

	if len(result) != 4 {
		t.Fatalf("Expected 4 bytes in memory, got %d", len(result))
	}

	for i := 0; i < 4; i++ {
		if result[i] != expected[i] {
			t.Errorf("Byte %d: expected 0x%02x, got 0x%02x", i, expected[i], result[i])
		}
	}

	// Stack should be empty
	if v.Stack().Len() != 0 {
		t.Errorf("Expected empty stack, got %d items", v.Stack().Len())
	}
}

func TestExtCodeCopyOpCode_OffsetBeyondCode(t *testing.T) {
	// Test copying beyond the code size (should fill with zeros)
	testAddr := [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
	testCode := []byte{0x60, 0x01} // Only 2 bytes

	code := []byte{
		// Copy 4 bytes starting from offset 10 (beyond code size)
		0x60, 0x04, // PUSH1 0x04 (size)
		0x60, 0x0a, // PUSH1 0x0a (offset = 10)
		0x60, 0x00, // PUSH1 0x00 (destOffset)
		0x73,                                                       // PUSH20 (address)
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, // address bytes
		0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90,
		0x3c, // EXTCODECOPY
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up mock state provider
	mock := &mockStateProvider{
		codeMap: map[[20]byte][]byte{
			testAddr: testCode,
		},
	}
	v.StateProvider = mock

	// Execute all instructions
	for !v.Stopped {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during execution: %v", err)
		}
	}

	// Check memory - should be all zeros since offset is beyond code
	result := v.Memory().Read(0, 4)
	expected := []byte{0x00, 0x00, 0x00, 0x00}

	for i := 0; i < 4; i++ {
		if result[i] != expected[i] {
			t.Errorf("Byte %d: expected 0x%02x, got 0x%02x", i, expected[i], result[i])
		}
	}
}

func TestExtCodeCopyOpCode_ZeroSize(t *testing.T) {
	code := []byte{
		0x60, 0x00, // PUSH1 0x00 (size = 0)
		0x60, 0x00, // PUSH1 0x00 (offset)
		0x60, 0x00, // PUSH1 0x00 (destOffset)
		0x60, 0x01, // PUSH1 0x01 (dummy address)
		0x3c, // EXTCODECOPY
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Execute all instructions
	for !v.Stopped {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during execution: %v", err)
		}
	}

	// Memory should remain empty since size is 0
	if v.Memory().Size() != 0 {
		t.Errorf("Expected memory size 0, got %d", v.Memory().Size())
	}

	// Stack should be empty
	if v.Stack().Len() != 0 {
		t.Errorf("Expected empty stack, got %d items", v.Stack().Len())
	}
}
