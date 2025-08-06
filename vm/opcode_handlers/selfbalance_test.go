package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSelfBalanceOpCode_Execute(t *testing.T) {
	code := []byte{
		0x47, // SELFBALANCE
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with a balance
	expectedBalance := uint256.NewInt(123456789)
	v.Context = &vm.ExecutionContext{
		Balance: expectedBalance,
	}

	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during execution: %v", err)
	}

	// Check that the balance was pushed onto the stack
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	balance, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if balance.Cmp(expectedBalance) != 0 {
		t.Errorf("Expected balance %s, got %s", expectedBalance.String(), balance.String())
	}
}

func TestSelfBalanceOpCode_NoContext(t *testing.T) {
	code := []byte{
		0x47, // SELFBALANCE
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
