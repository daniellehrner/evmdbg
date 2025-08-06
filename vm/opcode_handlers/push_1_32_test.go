package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestPush1(t *testing.T) {
	// Test: PUSH1 with single byte
	code := []byte{
		vm.PUSH1, 0x42, // PUSH1 0x42
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x42)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestPush2(t *testing.T) {
	// Test: PUSH2 with two bytes
	code := []byte{
		vm.PUSH2, 0x12, 0x34, // PUSH2 0x1234
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x1234)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestPush4(t *testing.T) {
	// Test: PUSH4 with four bytes
	code := []byte{
		vm.PUSH4, 0x12, 0x34, 0x56, 0x78, // PUSH4 0x12345678
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x12345678)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestPush8(t *testing.T) {
	// Test: PUSH8 with eight bytes
	code := []byte{
		vm.PUSH8, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, // PUSH8
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0).SetBytes([]byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88})
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestPush16(t *testing.T) {
	// Test: PUSH16 with sixteen bytes
	data := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
	}

	code := []byte{vm.PUSH16}
	code = append(code, data...)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0).SetBytes(data)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestPush32(t *testing.T) {
	// Test: PUSH32 with full 32 bytes
	data := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
		0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
	}

	code := []byte{vm.PUSH32}
	code = append(code, data...)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0).SetBytes(data)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestPushZeros(t *testing.T) {
	// Test: PUSH operations with zero bytes
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0
		vm.PUSH2, 0x00, 0x00, // PUSH2 0
		vm.PUSH4, 0x00, 0x00, 0x00, 0x00, // PUSH4 0
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// All three values should be zero
	for i := 0; i < 3; i++ {
		result, _ := d.Stack.Pop()
		if !result.IsZero() {
			t.Fatalf("expected 0 at position %d, got %s", i, result)
		}
	}
}

func TestPushMaxValues(t *testing.T) {
	// Test: PUSH operations with maximum values
	code := []byte{
		vm.PUSH1, 0xFF, // PUSH1 0xFF
		vm.PUSH2, 0xFF, 0xFF, // PUSH2 0xFFFF
		vm.PUSH4, 0xFF, 0xFF, 0xFF, 0xFF, // PUSH4 0xFFFFFFFF
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Check PUSH4 result (0xFFFFFFFF)
	result4, _ := d.Stack.Pop()
	expected4 := uint256.NewInt(0xFFFFFFFF)
	if result4.Cmp(expected4) != 0 {
		t.Fatalf("PUSH4: expected %s, got %s", expected4, result4)
	}

	// Check PUSH2 result (0xFFFF)
	result2, _ := d.Stack.Pop()
	expected2 := uint256.NewInt(0xFFFF)
	if result2.Cmp(expected2) != 0 {
		t.Fatalf("PUSH2: expected %s, got %s", expected2, result2)
	}

	// Check PUSH1 result (0xFF)
	result1, _ := d.Stack.Pop()
	expected1 := uint256.NewInt(0xFF)
	if result1.Cmp(expected1) != 0 {
		t.Fatalf("PUSH1: expected %s, got %s", expected1, result1)
	}
}

func TestPushSequence(t *testing.T) {
	// Test: Sequence of PUSH operations of different sizes
	code := []byte{
		vm.PUSH1, 0x01, // PUSH1 1
		vm.PUSH2, 0x02, 0x02, // PUSH2 0x0202
		vm.PUSH3, 0x03, 0x03, 0x03, // PUSH3 0x030303
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Stack should have [1, 0x0202, 0x030303] (most recent on top)
	if d.Stack.Len() != 3 {
		t.Fatalf("expected 3 items on stack, got %d", d.Stack.Len())
	}

	// Check PUSH3 result (0x030303)
	result3, _ := d.Stack.Pop()
	expected3 := uint256.NewInt(0x030303)
	if result3.Cmp(expected3) != 0 {
		t.Fatalf("PUSH3: expected %s, got %s", expected3, result3)
	}

	// Check PUSH2 result (0x0202)
	result2, _ := d.Stack.Pop()
	expected2 := uint256.NewInt(0x0202)
	if result2.Cmp(expected2) != 0 {
		t.Fatalf("PUSH2: expected %s, got %s", expected2, result2)
	}

	// Check PUSH1 result (1)
	result1, _ := d.Stack.Pop()
	expected1 := uint256.NewInt(1)
	if result1.Cmp(expected1) != 0 {
		t.Fatalf("PUSH1: expected %s, got %s", expected1, result1)
	}
}

func TestPush32MaxValue(t *testing.T) {
	// Test: PUSH32 with maximum 256-bit value (all 0xFF)
	data := make([]byte, 32)
	for i := range data {
		data[i] = 0xFF
	}

	code := []byte{vm.PUSH32}
	code = append(code, data...)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := new(uint256.Int).SetAllOne()
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}
