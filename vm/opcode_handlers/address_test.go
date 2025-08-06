package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"

	"testing"
)

func TestAddress(t *testing.T) {
	addr := [20]byte{0xde, 0xad, 0xbe, 0xef}
	code := []byte{vm.ADDRESS}

	d := vm.NewDebuggerVM(code, GetHandler)
	d.Context = &vm.ExecutionContext{Address: addr}

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	got, err := d.Stack.Pop()
	if err != nil {
		t.Fatalf("stack error: %v", err)
	}

	expected := new(uint256.Int).SetBytes(addr[:])
	if got.Cmp(expected) != 0 {
		t.Fatalf("expected %x, got %x", expected, got)
	}
}

func TestAddressWithoutContext(t *testing.T) {
	d := vm.NewDebuggerVM([]byte{vm.ADDRESS}, GetHandler)

	err := d.Step()
	if err == nil {
		t.Fatal("expected error when executing ADDRESS without context, got nil")
	}

	expectedErr := "address op code requires the execution context to be set"
	if err.Error() != expectedErr {
		t.Fatalf("expected error '%s', got '%s'", expectedErr, err.Error())
	}
}
