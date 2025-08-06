package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
)

func TestMStore8OpCode_Execute(t *testing.T) {
	tests := []struct {
		name     string
		code     []byte
		expected []byte
	}{
		{
			name: "Store byte 0x42 at address 0x00",
			code: []byte{
				0x60, 0x42, // PUSH1 0x42
				0x60, 0x00, // PUSH1 0x00 (address)
				0x53, // MSTORE8
			},
			expected: []byte{0x42},
		},
		{
			name: "Store byte 0xff at address 0x10",
			code: []byte{
				0x60, 0xff, // PUSH1 0xff
				0x60, 0x10, // PUSH1 0x10 (address)
				0x53, // MSTORE8
			},
			expected: []byte{0xff},
		},
		{
			name: "Store only LSB of larger value",
			code: []byte{
				0x61, 0x01, 0x23, // PUSH2 0x0123 (only 0x23 should be stored)
				0x60, 0x00, // PUSH1 0x00 (address)
				0x53, // MSTORE8
			},
			expected: []byte{0x23},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vm.NewDebuggerVM(tt.code, GetHandler)

			for !v.Stopped {
				err := v.Step()
				if err != nil {
					t.Fatalf("Unexpected error during execution: %v", err)
				}
			}

			// Check memory content
			var addr int
			if tt.name == "Store byte 0xff at address 0x10" {
				addr = 0x10
			} else {
				addr = 0x00
			}

			result := v.Memory().Read(addr, 1)
			if len(result) != 1 || result[0] != tt.expected[0] {
				t.Errorf("Expected memory[%d] = 0x%02x, got 0x%02x", addr, tt.expected[0], result[0])
			}

			// Verify stack is empty
			if v.Stack().Len() != 0 {
				t.Errorf("Expected empty stack, got %d items", v.Stack().Len())
			}
		})
	}
}

func TestMStore8OpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0x60, 0x42, // PUSH1 0x42 (only one item on stack)
		0x53, // MSTORE8 (needs two items)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Step through PUSH1
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH1: %v", err)
	}

	// MSTORE8 should fail with stack underflow
	err = v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}
