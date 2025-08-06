package opcode_handlers

import (
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestTransientStorageClear(t *testing.T) {
	// Test that ClearTransientStorage properly clears all transient storage
	v := vm.NewDebuggerVM([]byte{}, GetHandler)

	// Store some values in transient storage
	slot1 := uint256.NewInt(0x10)
	slot2 := uint256.NewInt(0x20)
	value1 := uint256.NewInt(0xdead)
	value2 := uint256.NewInt(0xbeef)

	v.WriteTransientStorage(slot1, value1)
	v.WriteTransientStorage(slot2, value2)

	// Verify values are stored
	storedValue1 := v.ReadTransientStorage(slot1)
	storedValue2 := v.ReadTransientStorage(slot2)

	if storedValue1.Cmp(value1) != 0 {
		t.Errorf("Expected stored value1 %s, got %s", value1.Hex(), storedValue1.Hex())
	}
	if storedValue2.Cmp(value2) != 0 {
		t.Errorf("Expected stored value2 %s, got %s", value2.Hex(), storedValue2.Hex())
	}

	// Clear transient storage
	v.ClearTransientStorage()

	// Verify values are cleared (should return zero)
	clearedValue1 := v.ReadTransientStorage(slot1)
	clearedValue2 := v.ReadTransientStorage(slot2)

	if !clearedValue1.IsZero() {
		t.Errorf("Expected zero after clear, got %s", clearedValue1.Hex())
	}
	if !clearedValue2.IsZero() {
		t.Errorf("Expected zero after clear, got %s", clearedValue2.Hex())
	}
}

func TestTransientStorageIndependenceFromPersistentStorage(t *testing.T) {
	// Test that transient storage is independent from persistent storage
	v := vm.NewDebuggerVM([]byte{}, GetHandler)

	slot := uint256.NewInt(0x42)
	persistentValue := uint256.NewInt(0x1111)
	transientValue := uint256.NewInt(0x2222)

	// Store in both persistent and transient storage with same slot
	v.WriteStorage(slot, persistentValue)
	v.WriteTransientStorage(slot, transientValue)

	// Verify both storages have different values
	storedPersistent := v.ReadStorage(slot)
	storedTransient := v.ReadTransientStorage(slot)

	if storedPersistent.Cmp(persistentValue) != 0 {
		t.Errorf("Expected persistent value %s, got %s", persistentValue.Hex(), storedPersistent.Hex())
	}
	if storedTransient.Cmp(transientValue) != 0 {
		t.Errorf("Expected transient value %s, got %s", transientValue.Hex(), storedTransient.Hex())
	}

	// Clear transient storage
	v.ClearTransientStorage()

	// Persistent storage should be unaffected
	storedPersistentAfter := v.ReadStorage(slot)
	storedTransientAfter := v.ReadTransientStorage(slot)

	if storedPersistentAfter.Cmp(persistentValue) != 0 {
		t.Errorf("Persistent storage should be unaffected, expected %s, got %s",
			persistentValue.Hex(), storedPersistentAfter.Hex())
	}
	if !storedTransientAfter.IsZero() {
		t.Errorf("Transient storage should be cleared, expected zero, got %s", storedTransientAfter.Hex())
	}
}

func TestMultipleTransactionSimulation(t *testing.T) {
	// Simulate multiple transactions with transient storage clearing between them
	code := []byte{
		0x60, 0x99, // PUSH1 0x99 (value)
		0x60, 0x01, // PUSH1 0x01 (slot)
		0x5d,       // TSTORE
		0x60, 0x01, // PUSH1 0x01 (slot)
		0x5c, // TLOAD
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Execute first "transaction"
	for !v.Stopped {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during first transaction: %v", err)
		}
	}

	// Should have stored and loaded the value
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack after first transaction, got %d", v.Stack().Len())
	}

	value, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	expectedValue := uint256.NewInt(0x99)
	if value.Cmp(expectedValue) != 0 {
		t.Errorf("Expected value %s in first transaction, got %s", expectedValue.Hex(), value.Hex())
	}

	// Clear transient storage to simulate end of transaction
	v.ClearTransientStorage()

	// Reset VM for second transaction (same code)
	v2 := vm.NewDebuggerVM(code, GetHandler)

	// Execute second "transaction" - transient storage should be empty
	for !v2.Stopped {
		err := v2.Step()
		if err != nil {
			t.Fatalf("Unexpected error during second transaction: %v", err)
		}
	}

	// Should have stored 0x99, but when loading, should get 0x99 (fresh transient storage)
	if v2.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack after second transaction, got %d", v2.Stack().Len())
	}

	value2, err := v2.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack in second transaction: %v", err)
	}

	// Second transaction should also get 0x99 since it stores then loads in same transaction
	if value2.Cmp(expectedValue) != 0 {
		t.Errorf("Expected value %s in second transaction, got %s", expectedValue.Hex(), value2.Hex())
	}

	// But the first VM's transient storage should be cleared
	slot := uint256.NewInt(0x01)
	clearedValue := v.ReadTransientStorage(slot)
	if !clearedValue.IsZero() {
		t.Errorf("First VM's transient storage should be cleared, expected zero, got %s", clearedValue.Hex())
	}
}
