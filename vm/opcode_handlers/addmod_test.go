package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestAddMod(t *testing.T) {
	code := []byte{
		vm.PUSH1, 0x05, // modulus
		vm.PUSH1, 0x03, // b
		vm.PUSH1, 0x04, // a
		vm.ADDMOD,
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	expected := uint256.NewInt((4 + 3) % 5)
	actual, _ := d.Stack.Pop()
	if actual.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}
