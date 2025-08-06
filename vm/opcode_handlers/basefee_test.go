package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestBaseFeeOpCode_Execute(t *testing.T) {
	code := []byte{
		0x48, // BASEFEE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with block context
	expectedBaseFee := uint256.NewInt(1000000000) // 1 Gwei
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BaseFee: expectedBaseFee,
		},
	}

	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during execution: %v", err)
	}

	// Check that the base fee was pushed onto the stack
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	baseFee, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if baseFee.Cmp(expectedBaseFee) != 0 {
		t.Errorf("Expected base fee %s, got %s", expectedBaseFee.String(), baseFee.String())
	}
}

func TestBaseFeeOpCode_NoContext(t *testing.T) {
	code := []byte{
		0x48, // BASEFEE
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

func TestBaseFeeOpCode_NoBlockContext(t *testing.T) {
	code := []byte{
		0x48, // BASEFEE
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
