package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestSloadBasic(t *testing.T) {
	// Test: SSTORE then SLOAD
	code := []byte{
		vm.PUSH1, 0x42, // PUSH1 0x42 (value)
		vm.PUSH1, 0x01, // PUSH1 0x01 (key)
		vm.SSTORE, // SSTORE

		vm.PUSH1, 0x01, // PUSH1 0x01 (key)
		vm.SLOAD, // SLOAD
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Should have loaded the stored value
	if d.Stack.Len() != 1 {
		t.Fatalf("expected 1 item on stack, got %d", d.Stack.Len())
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x42)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSloadUninitialized(t *testing.T) {
	// Test: SLOAD from uninitialized storage (should return 0)
	code := []byte{
		vm.PUSH1, 0x99, // PUSH1 0x99 (key)
		vm.SLOAD, // SLOAD
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0 from uninitialized storage, got %s", result)
	}
}

func TestSloadZeroValue(t *testing.T) {
	// Test: SSTORE zero then SLOAD
	code := []byte{
		vm.PUSH1, 0x00, // PUSH1 0 (value)
		vm.PUSH1, 0x05, // PUSH1 5 (key)
		vm.SSTORE, // SSTORE

		vm.PUSH1, 0x05, // PUSH1 5 (key)
		vm.SLOAD, // SLOAD
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	if !result.IsZero() {
		t.Fatalf("expected 0, got %s", result)
	}
}

func TestSloadLargeKey(t *testing.T) {
	// Test: SLOAD with large key
	largeKey := new(uint256.Int).SetAllOne()
	largeValue := uint256.NewInt(0x12345)

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeValue)...) // value
	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeKey)...) // key
	code = append(code, vm.SSTORE)                     // SSTORE

	code = append(code, vm.PUSH32)
	code = append(code, bytes32WithValue(largeKey)...) // key
	code = append(code, vm.SLOAD)                      // SLOAD

	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	if result.Cmp(largeValue) != 0 {
		t.Fatalf("expected %s, got %s", largeValue, result)
	}
}

func TestSloadLargeValue(t *testing.T) {
	// Test: SLOAD of large stored value
	largeValue := new(uint256.Int).SetBytes([]byte{
		0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88,
		0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF,
		0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10,
	})

	code := []byte{vm.PUSH32}
	code = append(code, bytes32WithValue(largeValue)...) // value
	code = append(code, vm.PUSH1, 0x20)                  // PUSH1 32 (key)
	code = append(code, vm.SSTORE)                       // SSTORE

	code = append(code, vm.PUSH1, 0x20) // PUSH1 32 (key)
	code = append(code, vm.SLOAD)       // SLOAD

	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	if result.Cmp(largeValue) != 0 {
		t.Fatalf("expected %s, got %s", largeValue, result)
	}
}

func TestSloadMultipleKeys(t *testing.T) {
	// Test: Store multiple values and load them back
	code := []byte{
		// Store at key 10
		vm.PUSH1, 0xAA, // PUSH1 0xAA
		vm.PUSH1, 0x0A, // PUSH1 10
		vm.SSTORE, // SSTORE

		// Store at key 20
		vm.PUSH1, 0xBB, // PUSH1 0xBB
		vm.PUSH1, 0x14, // PUSH1 20
		vm.SSTORE, // SSTORE

		// Store at key 30
		vm.PUSH1, 0xCC, // PUSH1 0xCC
		vm.PUSH1, 0x1E, // PUSH1 30
		vm.SSTORE, // SSTORE

		// Load from key 20 (middle one)
		vm.PUSH1, 0x14, // PUSH1 20
		vm.SLOAD, // SLOAD
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0xBB) // Should load value from key 20
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSloadAfterOverwrite(t *testing.T) {
	// Test: SLOAD after overwriting a value
	code := []byte{
		// Initial store
		vm.PUSH1, 0x11, // PUSH1 0x11
		vm.PUSH1, 0x07, // PUSH1 7 (key)
		vm.SSTORE, // SSTORE

		// Overwrite
		vm.PUSH1, 0x22, // PUSH1 0x22
		vm.PUSH1, 0x07, // PUSH1 7 (same key)
		vm.SSTORE, // SSTORE

		// Load
		vm.PUSH1, 0x07, // PUSH1 7 (key)
		vm.SLOAD, // SLOAD
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x22) // Should be the overwritten value
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSloadZeroKey(t *testing.T) {
	// Test: SLOAD from key 0
	code := []byte{
		vm.PUSH1, 0x88, // PUSH1 0x88 (value)
		vm.PUSH1, 0x00, // PUSH1 0 (key)
		vm.SSTORE, // SSTORE

		vm.PUSH1, 0x00, // PUSH1 0 (key)
		vm.SLOAD, // SLOAD
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	result, _ := d.Stack.Pop()
	expected := uint256.NewInt(0x88)
	if result.Cmp(expected) != 0 {
		t.Fatalf("expected %s, got %s", expected, result)
	}
}

func TestSloadSequential(t *testing.T) {
	// Test: Multiple sequential SLOAD operations
	code := []byte{
		// Store three values
		vm.PUSH1, 0x11, // PUSH1 0x11
		vm.PUSH1, 0x01, // PUSH1 1
		vm.SSTORE, // SSTORE

		vm.PUSH1, 0x22, // PUSH1 0x22
		vm.PUSH1, 0x02, // PUSH1 2
		vm.SSTORE, // SSTORE

		vm.PUSH1, 0x33, // PUSH1 0x33
		vm.PUSH1, 0x03, // PUSH1 3
		vm.SSTORE, // SSTORE

		// Load all three in reverse order
		vm.PUSH1, 0x03, // PUSH1 3
		vm.SLOAD, // SLOAD (loads 0x33)

		vm.PUSH1, 0x02, // PUSH1 2
		vm.SLOAD, // SLOAD (loads 0x22)

		vm.PUSH1, 0x01, // PUSH1 1
		vm.SLOAD, // SLOAD (loads 0x11)
	}
	d := vm.NewDebuggerVM(code, GetHandler)
	d.Storage = make(map[string]*uint256.Int)

	for !d.Stopped {
		err := d.Step()
		if err != nil {
			t.Fatalf("execution error: %v", err)
		}
	}

	// Stack should have [0x33, 0x22, 0x11] (most recent on top)
	if d.Stack.Len() != 3 {
		t.Fatalf("expected 3 items on stack, got %d", d.Stack.Len())
	}

	// Check in LIFO order
	result1, _ := d.Stack.Pop() // 0x11
	expected1 := uint256.NewInt(0x11)
	if result1.Cmp(expected1) != 0 {
		t.Fatalf("expected %s, got %s", expected1, result1)
	}

	result2, _ := d.Stack.Pop() // 0x22
	expected2 := uint256.NewInt(0x22)
	if result2.Cmp(expected2) != 0 {
		t.Fatalf("expected %s, got %s", expected2, result2)
	}

	result3, _ := d.Stack.Pop() // 0x33
	expected3 := uint256.NewInt(0x33)
	if result3.Cmp(expected3) != 0 {
		t.Fatalf("expected %s, got %s", expected3, result3)
	}
}
