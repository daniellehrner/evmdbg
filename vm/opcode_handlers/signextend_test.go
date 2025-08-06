package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSignExtendByte0Negative(t *testing.T) {
	// Test: SIGNEXTEND(0, 0xFF) should extend sign from byte 0
	// 0xFF in byte 0 is negative, so result should be all 1s
	code := []byte{
		vm.PUSH1, 0xFF, // PUSH1 0xFF (value with negative sign bit in byte 0)
		vm.PUSH1, 0x00, // PUSH1 0 (extend from byte 0)
		vm.SIGNEXTEND, // SIGNEXTEND
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := new(uint256.Int).SetAllOne() // All bits set (-1)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSignExtendByte0Positive(t *testing.T) {
	// Test: SIGNEXTEND(0, 0x7F) should not extend (positive)
	// 0x7F in byte 0 is positive, so result should remain 0x7F
	code := []byte{
		vm.PUSH1, 0x7F, // PUSH1 0x7F (value with positive sign bit in byte 0)
		vm.PUSH1, 0x00, // PUSH1 0 (extend from byte 0)
		vm.SIGNEXTEND, // SIGNEXTEND
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(0x7F)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSignExtendByte1Negative(t *testing.T) {
	// Test: SIGNEXTEND(1, 0x80FF) should extend sign from byte 1
	// Byte 1 = 0x80 (negative), so should extend to fill upper bytes
	code := []byte{
		vm.PUSH2, 0x80, 0xFF, // PUSH2 0x80FF (negative in byte 1)
		vm.PUSH1, 0x01, // PUSH1 1 (extend from byte 1)
		vm.SIGNEXTEND, // SIGNEXTEND
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()

	// Expected: 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF80FF
	// Sign bit from byte 1 (0x80) extends to fill all upper bytes
	expected := new(uint256.Int).SetBytes([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0xFF,
	})

	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSignExtendByte1Positive(t *testing.T) {
	// Test: SIGNEXTEND(1, 0x7FFF) should not extend (positive in byte 1)
	code := []byte{
		vm.PUSH2, 0x7F, 0xFF, // PUSH2 0x7FFF (positive in byte 1)
		vm.PUSH1, 0x01, // PUSH1 1 (extend from byte 1)
		vm.SIGNEXTEND, // SIGNEXTEND
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	expected := uint256.NewInt(0x7FFF)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSignExtendByte31NoChange(t *testing.T) {
	// Test: SIGNEXTEND(31, value) should return value unchanged
	// Byte 31 is already the most significant byte
	testValue := new(uint256.Int).SetBytes([]byte{
		0x80, 0x12, 0x34, 0x56, 0x78, 0x9A, 0xBC, 0xDE,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11,
		0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(testValue)...)
	code = append(code, vm.PUSH1, 31) // PUSH1 31 (byte 31)
	code = append(code, vm.SIGNEXTEND)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if result.Cmp(testValue) != 0 {
		t.Fatalf("expected %s, got %s", testValue, result)
	}
}

func TestSignExtendLargeByteNum(t *testing.T) {
	// Test: SIGNEXTEND(32, value) should return value unchanged (k > 30)
	testValue := uint256.NewInt(0x123456)

	code := []byte{
		vm.PUSH4, 0x00, 0x12, 0x34, 0x56, // PUSH4 0x123456
		vm.PUSH1, 32, // PUSH1 32 (k > 30)
		vm.SIGNEXTEND, // SIGNEXTEND
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if result.Cmp(testValue) != 0 {
		t.Fatalf("expected %s, got %s", testValue, result)
	}
}

func TestSignExtendZeroValue(t *testing.T) {
	// Test: SIGNEXTEND(0, 0) should return 0
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (value)
		vm.PUSH1, 0x00, // PUSH1 0 (byte position)
		vm.SIGNEXTEND, // SIGNEXTEND
	}
	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0, got %s", result)
	}
}

func TestSignExtendByte15(t *testing.T) {
	// Test: SIGNEXTEND(15, value) with negative sign in byte 15
	// Create a value with 0x80 in byte 15 (16th byte from right)
	testBytes := make([]byte, 32)
	testBytes[31-15] = 0x80 // Byte 15 from the right (big-endian)
	testBytes[31] = 0xFF    // Some data in lower bytes

	testValue := new(uint256.Int).SetBytes(testBytes)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(testValue)...)
	code = append(code, vm.PUSH1, 15) // PUSH1 15
	code = append(code, vm.SIGNEXTEND)

	d := vm.NewDebuggerVM(code, GetHandler)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack().Pop()

	// Check that bytes 16-31 are all 0xFF (sign extended)
	// and bytes 0-15 contain original data
	resultBytes := result.Bytes32()

	// Bytes 0-15 should be 0xFF (sign extension)
	for i := 0; i < 16; i++ {
		if resultBytes[i] != 0xFF {
			t.Fatalf("byte %d should be 0xFF, got 0x%02X", i, resultBytes[i])
		}
	}

	// Byte 16 should be 0x80 (original sign byte)
	if resultBytes[16] != 0x80 {
		t.Fatalf("byte 16 should be 0x80, got 0x%02X", resultBytes[16])
	}

	// Byte 31 should be 0xFF (original data)
	if resultBytes[31] != 0xFF {
		t.Fatalf("byte 31 should be 0xFF, got 0x%02X", resultBytes[31])
	}
}

func TestSignExtendBoundaryValues(t *testing.T) {
	// Test boundary between positive and negative (0x7F vs 0x80)

	// Test 0x80 (negative boundary)
	code1 := []byte{
		vm.PUSH1, 0x80, // PUSH1 0x80 (negative boundary)
		vm.PUSH1, 0x00, // PUSH1 0 (extend from byte 0)
		vm.SIGNEXTEND, // SIGNEXTEND
	}
	d1 := vm.NewDebuggerVM(code1, GetHandler)

	for !d1.Stopped {
		err := d1.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result1, _ := d1.Stack().Pop()
	expected1 := new(uint256.Int).SetAllOne()                                                                      // Should be -1 (all 1s)
	expected1 = expected1.And(expected1, new(uint256.Int).Sub(new(uint256.Int).SetAllOne(), uint256.NewInt(0x7F))) // Clear lower 7 bits, keep 0x80
	expected1 = expected1.Or(expected1, uint256.NewInt(0x80))

	// Actually, 0x80 sign extended from byte 0 should be 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF80
	expected1 = new(uint256.Int).SetBytes([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80,
	})

	if result1.Cmp(expected1) != 0 {
		t.Fatalf("0x80 sign extend: expected %s, got %s", expected1, result1)
	}

	// Test 0x7F (positive boundary)
	code2 := []byte{
		vm.PUSH1, 0x7F, // PUSH1 0x7F (positive boundary)
		vm.PUSH1, 0x00, // PUSH1 0 (extend from byte 0)
		vm.SIGNEXTEND, // SIGNEXTEND
	}
	d2 := vm.NewDebuggerVM(code2, GetHandler)

	for !d2.Stopped {
		err := d2.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result2, _ := d2.Stack().Pop()
	expected2 := uint256.NewInt(0x7F) // Should remain 0x7F

	if result2.Cmp(expected2) != 0 {
		t.Fatalf("0x7F sign extend: expected %s, got %s", expected2, result2)
	}
}
