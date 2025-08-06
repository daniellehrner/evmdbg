package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestAdd(t *testing.T) {
	code := []byte{vm.PUSH1, 0x03, vm.PUSH1, 0x05, vm.ADD}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	if d.Stack.Len() != 1 {
		t.Fatalf("expected 1 item on the stack, got %d", d.Stack.Len())
	}

	top, _ := d.Stack.Pop()
	if top.Cmp(uint256.NewInt(8)) != 0 {
		t.Fatalf("expected 8 on stack, got %s", top)
	}
}

func TestAddOverflow(t *testing.T) {
	m := new(uint256.Int).Sub(new(uint256.Int).Lsh(uint256.NewInt(1), 256), uint256.NewInt(1)) // 2^256 - 1
	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(uint256.NewInt(1))...) // b = 1
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(m)...) // a = 2^256 - 1
	code = append(code, vm.ADD)

	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = map[string]*uint256.Int{}

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

	// EVM ADD is modulo 2^256: (2^256 - 1 + 1) % 2^256 = 0
	if got.Sign() != 0 {
		t.Fatalf("expected 0, got %s", got.String())
	}
}

// helper for PUSH32 values
func bytes32WithValue(val *uint256.Int) []byte {
	b := val.Bytes()
	out := make([]byte, 32)
	copy(out[32-len(b):], b)
	return out
}
