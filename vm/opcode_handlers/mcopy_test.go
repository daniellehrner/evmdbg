package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
)

func TestMCopyOpCode_Execute(t *testing.T) {
	tests := []struct {
		name string
		code []byte
	}{
		{
			name: "Copy 4 bytes from offset 0 to offset 32",
			code: []byte{
				// First store some data at offset 0
				0x7f, 0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe, // PUSH32 0xdeadbeefcafebabe...
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x60, 0x00, // PUSH1 0x00 (memory offset)
				0x52, // MSTORE

				// Now copy 4 bytes from offset 0 to offset 32
				0x60, 0x04, // PUSH1 0x04 (size)
				0x60, 0x00, // PUSH1 0x00 (source offset)
				0x60, 0x20, // PUSH1 0x20 (dest offset = 32)
				0x5e, // MCOPY
			},
		},
		{
			name: "Copy zero bytes (should be no-op)",
			code: []byte{
				0x60, 0x00, // PUSH1 0x00 (size = 0)
				0x60, 0x10, // PUSH1 0x10 (source offset)
				0x60, 0x20, // PUSH1 0x20 (dest offset)
				0x5e, // MCOPY
			},
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

			if tt.name == "Copy 4 bytes from offset 0 to offset 32" {
				// Check that the data was copied correctly
				source := v.Memory().Read(0, 4)
				dest := v.Memory().Read(32, 4)

				if len(source) != 4 || len(dest) != 4 {
					t.Fatalf("Expected 4 bytes, got source=%d dest=%d", len(source), len(dest))
				}

				for i := 0; i < 4; i++ {
					if source[i] != dest[i] {
						t.Errorf("Byte %d: source=0x%02x, dest=0x%02x", i, source[i], dest[i])
					}
				}

				// Verify expected values
				expected := []byte{0xde, 0xad, 0xbe, 0xef}
				for i := 0; i < 4; i++ {
					if dest[i] != expected[i] {
						t.Errorf("Expected dest[%d]=0x%02x, got 0x%02x", i, expected[i], dest[i])
					}
				}
			}

			// Verify stack is empty
			if v.Stack().Len() != 0 {
				t.Errorf("Expected empty stack, got %d items", v.Stack().Len())
			}
		})
	}
}

func TestMCopyOpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0x60, 0x04, // PUSH1 0x04
		0x60, 0x00, // PUSH1 0x00 (only two items on stack)
		0x5e, // MCOPY (needs three items)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Step through PUSH instructions
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during PUSH: %v", err)
		}
	}

	// MCOPY should fail with stack underflow
	err := v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}
