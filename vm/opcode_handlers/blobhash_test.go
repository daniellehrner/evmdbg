package opcode_handlers

import (
	"strings"
	"testing"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/holiman/uint256"
)

func TestBlobHashOpCode_Execute(t *testing.T) {
	code := []byte{
		0x60, 0x01, // PUSH1 0x01 (index 1)
		0x49, // BLOBHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with blob hashes
	blobHash0 := [32]byte{0xaa, 0xbb, 0xcc, 0xdd}
	blobHash1 := [32]byte{0x11, 0x22, 0x33, 0x44}
	blobHash2 := [32]byte{0xff, 0xee, 0xdd, 0xcc}

	for i := 4; i < 32; i++ {
		blobHash0[i] = byte(i)
		blobHash1[i] = byte(i + 10)
		blobHash2[i] = byte(i + 20)
	}

	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobHashes: [][32]byte{blobHash0, blobHash1, blobHash2},
		},
	}

	// Execute PUSH1 and BLOBHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should return blobHash1 (index 1)
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	expectedHash := new(uint256.Int).SetBytes(blobHash1[:])
	if hash.Cmp(expectedHash) != 0 {
		t.Errorf("Expected blob hash %s, got %s", expectedHash.Hex(), hash.Hex())
	}
}

func TestBlobHashOpCode_IndexZero(t *testing.T) {
	code := []byte{
		0x60, 0x00, // PUSH1 0x00 (index 0)
		0x49, // BLOBHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with blob hashes
	blobHash0 := [32]byte{0xde, 0xad, 0xbe, 0xef}
	for i := 4; i < 32; i++ {
		blobHash0[i] = byte(0x42)
	}

	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobHashes: [][32]byte{blobHash0},
		},
	}

	// Execute PUSH1 and BLOBHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	expectedHash := new(uint256.Int).SetBytes(blobHash0[:])
	if hash.Cmp(expectedHash) != 0 {
		t.Errorf("Expected blob hash %s, got %s", expectedHash.Hex(), hash.Hex())
	}
}

func TestBlobHashOpCode_OutOfBounds(t *testing.T) {
	code := []byte{
		0x60, 0x05, // PUSH1 0x05 (index 5, out of bounds)
		0x49, // BLOBHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with only 2 blob hashes
	blobHash0 := [32]byte{0xaa, 0xbb}
	blobHash1 := [32]byte{0xcc, 0xdd}

	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobHashes: [][32]byte{blobHash0, blobHash1},
		},
	}

	// Execute PUSH1 and BLOBHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should be 0 for out of bounds index
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !hash.IsZero() {
		t.Errorf("Expected zero for out of bounds index, got %s", hash.Hex())
	}
}

func TestBlobHashOpCode_EmptyBlobHashes(t *testing.T) {
	code := []byte{
		0x60, 0x00, // PUSH1 0x00 (index 0)
		0x49, // BLOBHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with no blob hashes
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobHashes: [][32]byte{}, // Empty slice
		},
	}

	// Execute PUSH1 and BLOBHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should be 0 since no blob hashes exist
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !hash.IsZero() {
		t.Errorf("Expected zero for empty blob hashes, got %s", hash.Hex())
	}
}

func TestBlobHashOpCode_LargeIndex(t *testing.T) {
	code := []byte{
		0x61, 0xff, 0xff, // PUSH2 0xffff (very large index)
		0x49, // BLOBHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up execution context with one blob hash
	blobHash := [32]byte{0x12, 0x34, 0x56, 0x78}
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobHashes: [][32]byte{blobHash},
		},
	}

	// Execute PUSH2 and BLOBHASH
	for i := 0; i < 2; i++ {
		err := v.Step()
		if err != nil {
			t.Fatalf("Unexpected error during step %d: %v", i, err)
		}
	}

	// Check result - should be 0 for large out of bounds index
	if v.Stack().Len() != 1 {
		t.Fatalf("Expected 1 item on stack, got %d", v.Stack().Len())
	}

	hash, err := v.Stack().Peek(0)
	if err != nil {
		t.Fatalf("Error peeking at stack: %v", err)
	}

	if !hash.IsZero() {
		t.Errorf("Expected zero for large out of bounds index, got %s", hash.Hex())
	}
}

func TestBlobHashOpCode_NoContext(t *testing.T) {
	code := []byte{
		0x60, 0x00, // PUSH1 0x00
		0x49, // BLOBHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)
	// Don't set context

	// Execute PUSH1
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH1: %v", err)
	}

	// BLOBHASH should fail without context
	err = v.Step()
	if err == nil {
		t.Fatal("Expected error when context is not set, got nil")
	}

	if !strings.Contains(err.Error(), "execution context") {
		t.Errorf("Expected error to mention 'execution context', got: %v", err)
	}
}

func TestBlobHashOpCode_NoBlockContext(t *testing.T) {
	code := []byte{
		0x60, 0x00, // PUSH1 0x00
		0x49, // BLOBHASH
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set context but no block context
	v.Context = &vm.ExecutionContext{}

	// Execute PUSH1
	err := v.Step()
	if err != nil {
		t.Fatalf("Unexpected error during PUSH1: %v", err)
	}

	// BLOBHASH should fail without block context
	err = v.Step()
	if err == nil {
		t.Fatal("Expected error when block context is not set, got nil")
	}

	if !strings.Contains(err.Error(), "block context") {
		t.Errorf("Expected error to mention 'block context', got: %v", err)
	}
}

func TestBlobHashOpCode_StackUnderflow(t *testing.T) {
	code := []byte{
		0x49, // BLOBHASH (no index on stack)
	}

	v := vm.NewDebuggerVM(code, GetHandler)

	// Set up context
	v.Context = &vm.ExecutionContext{
		Block: &vm.BlockContext{
			BlobHashes: [][32]byte{},
		},
	}

	// BLOBHASH should fail with stack underflow
	err := v.Step()
	if err == nil {
		t.Fatal("Expected stack underflow error, got nil")
	}
}
