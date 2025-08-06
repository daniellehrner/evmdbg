package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestBlobBaseFeeOpCode_Execute(t *testing.T) {
	code := []byte{
		0x4a, // BLOBBASEFEE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with block context
	expectedBlobBaseFee := uint256.NewInt(2000000000) // 2 Gwei
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobBaseFee: expectedBlobBaseFee,
		},
	}

	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during execution: %v", err)
	}

	// Check that the blob base fee was pushed onto the stack
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	blobBaseFee, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if blobBaseFee.Cmp(expectedBlobBaseFee) != 0 {
		t.Errorf("Expected blob base fee %s, got %s", expectedBlobBaseFee.String(), blobBaseFee.String())
	}
}

func TestBlobBaseFeeOpCode_ZeroFee(t *testing.T) {
	code := []byte{
		0x4a, // BLOBBASEFEE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with zero blob base fee
	expectedBlobBaseFee := uint256.NewInt(0)
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobBaseFee: expectedBlobBaseFee,
		},
	}

	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during execution: %v", err)
	}

	// Check result
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	blobBaseFee, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !blobBaseFee.IsZero() {
		t.Errorf("Expected zero blob base fee, got %s", blobBaseFee.String())
	}
}

func TestBlobBaseFeeOpCode_NoContext(t *testing.T) {
	code := []byte{
		0x4a, // BLOBBASEFEE
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set context

	err := v.Step()
	if err == nil {
		t.Fatal("Expected error when context is not set, got nil")
	}

	if !strings.Contains(err.Error(), "execution context") {
		t.Errorf("Expected error to mention 'execution context', got: %v", err)
	}
}

func TestBlobBaseFeeOpCode_NoBlockContext(t *testing.T) {
	code := []byte{
		0x4a, // BLOBBASEFEE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set context but no block context
	v.Context = &vm.ExecutionContext{}

	err := v.Step()
	if err == nil {
		t.Fatal("Expected error when block context is not set, got nil")
	}

	if !strings.Contains(err.Error(), "block context") {
		t.Errorf("Expected error to mention 'block context', got: %v", err)
	}
}

func TestBlobBaseFeeOpCode_LargeFee(t *testing.T) {
	code := []byte{
		0x4a, // BLOBBASEFEE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with large blob base fee
	expectedBlobBaseFee, _ := uint256.FromDecimal("123456789012345678901234567890") // Large number
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobBaseFee: expectedBlobBaseFee,
		},
	}

	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during execution: %v", err)
	}

	// Check result
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	blobBaseFee, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if blobBaseFee.Cmp(expectedBlobBaseFee) != 0 {
		t.Errorf("Expected large blob base fee %s, got %s", expectedBlobBaseFee.String(), blobBaseFee.String())
	}
}
