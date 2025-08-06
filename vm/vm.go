package vm

import (
	"errors"
	"fmt"

	"github.com/holiman/uint256"
)

// Errors
var (
	ErrInvalidSHA3           = errors.New("invalid SHA3 hash calculation")
	ErrStackUnderflow        = errors.New("stack underflow")
	ErrStackOverflow         = errors.New("stack overflow")
	ErrOutOfGas              = errors.New("out of gas")
	ErrInvalidJump           = errors.New("invalid jump destination")
	ErrCallDepthLimit        = errors.New("call depth limit exceeded")
	ErrStaticCallStateChange = errors.New("state change operation in static call context")
)

type Handler interface {
	Execute(vm *DebuggerVM) error
}
type HandlerGetter func(b byte) Handler

type DebuggerVM struct {
	// Frame stack for call support
	frames []MessageFrame

	// VM state
	Storage          map[string]*uint256.Int
	TransientStorage map[string]*uint256.Int // EIP-1153: Transient storage
	Stopped          bool

	ReturnValue []byte
	Reverted    bool
	Logs        []LogEntry

	Context       *ExecutionContext
	HandlerGetter HandlerGetter
	StateProvider StateProvider

	// Return data from last call
	lastReturnData []byte
}

type LogEntry struct {
	Address [20]byte
	Topics  [][]byte
	Data    []byte
}

type ExecutionContext struct {
	Caller  [20]byte
	Address [20]byte
	Origin  [20]byte
	Value   *uint256.Int

	CallData []byte
	GasPrice *uint256.Int
	Gas      uint64
	Balance  *uint256.Int

	Block *BlockContext
}

type BlockContext struct {
	Coinbase    [20]byte
	Timestamp   uint64
	Number      uint64
	Difficulty  *uint256.Int
	GasLimit    uint64
	ChainID     *uint256.Int
	BaseFee     *uint256.Int
	BlobBaseFee *uint256.Int // EIP-4844: Blob base fee
	BlobHashes  [][32]byte   // EIP-4844: Versioned hashes of the blobs
}

type CodeMetadata struct {
	ValidPC   map[uint64]struct{}
	JumpDests map[uint64]struct{}
}

// CallType represents the type of call being made
type CallType int

const (
	CallTypeCall CallType = iota
	CallTypeCallCode
	CallTypeDelegateCall
	CallTypeStaticCall
)

// MessageFrame represents a single execution frame
type MessageFrame struct {
	Code         []byte
	PC           uint64
	Stack        *Stack
	Memory       *Memory
	ReturnData   []byte
	Gas          uint64
	CallType     CallType
	IsStatic     bool
	CodeMetadata *CodeMetadata
}

// CallContext contains information about a call
type CallContext struct {
	Caller   [20]byte
	Address  [20]byte
	Origin   [20]byte
	Value    *uint256.Int
	CallData []byte
	Gas      uint64
}

// StateProvider interface for accessing blockchain state
type StateProvider interface {
	GetBalance(addr [20]byte) *uint256.Int
	GetCode(addr [20]byte) []byte
	GetStorage(addr [20]byte, key *uint256.Int) *uint256.Int
	SetStorage(addr [20]byte, key *uint256.Int, value *uint256.Int)
	AccountExists(addr [20]byte) bool
	GetBlockHash(blockNumber uint64) [32]byte
	CreateAccount(addr [20]byte, code []byte, balance *uint256.Int) error
	GetNonce(addr [20]byte) uint64
	SetNonce(addr [20]byte, nonce uint64)
}

func NewDebuggerVM(code []byte, hg HandlerGetter) *DebuggerVM {
	stack := NewStack()
	memory := NewMemory()
	codeMetadata := scanCodeMetadata(code)

	// Create initial frame
	initialFrame := MessageFrame{
		Code:         code,
		PC:           0,
		Stack:        stack,
		Memory:       memory,
		ReturnData:   nil,
		Gas:          0, // Will be set via context
		CallType:     CallTypeCall,
		IsStatic:     false,
		CodeMetadata: codeMetadata,
	}

	vm := &DebuggerVM{
		frames:           []MessageFrame{initialFrame},
		Storage:          make(map[string]*uint256.Int),
		TransientStorage: make(map[string]*uint256.Int),
		HandlerGetter:    hg,
	}

	return vm
}

func (vm *DebuggerVM) Step() error {
	frame := vm.currentFrame()
	if frame == nil {
		return fmt.Errorf("no execution frame")
	}

	if vm.Stopped || int(frame.PC) >= len(frame.Code) {
		vm.Stopped = true
		return nil
	}

	op := frame.Code[frame.PC]
	frame.PC++

	handler := vm.HandlerGetter(op)
	if handler == nil {
		return fmt.Errorf("unsupported opcode: 0x%x", op)
	}

	return handler.Execute(vm)
}

func (vm *DebuggerVM) RunUntil(breakpoints map[uint64]struct{}) error {
	for {
		frame := vm.currentFrame()
		if frame == nil {
			return fmt.Errorf("no execution frame")
		}

		if vm.Stopped || int(frame.PC) >= len(frame.Code) {
			vm.Stopped = true
			return nil
		}

		if _, ok := breakpoints[frame.PC]; ok {
			return nil // reached a breakpoint
		}

		// Only execute at valid PC (not in PUSH immediate)
		if _, ok := frame.CodeMetadata.ValidPC[frame.PC]; !ok {
			return fmt.Errorf("invalid PC: 0x%x (likely inside PUSH immediate)", frame.PC)
		}

		err := vm.Step()
		if err != nil {
			return err
		}
	}
}

func (vm *DebuggerVM) ReadCodeByte(offset uint64) (byte, error) {
	frame := vm.currentFrame()
	if frame == nil {
		return 0, fmt.Errorf("no execution frame")
	}

	pos := frame.PC + offset
	if int(pos) >= len(frame.Code) {
		return 0, fmt.Errorf("code out of bounds at PC + %d", offset)
	}
	return frame.Code[pos], nil
}

func (vm *DebuggerVM) AdvancePC(n uint64) {
	frame := vm.currentFrame()
	if frame != nil {
		frame.PC += n
	}
}

func (vm *DebuggerVM) ReadCodeSlice(n uint64) ([]byte, error) {
	frame := vm.currentFrame()
	if frame == nil {
		return nil, fmt.Errorf("no execution frame")
	}

	if frame.PC+n > uint64(len(frame.Code)) {
		return nil, fmt.Errorf("code out of bounds: PC=%d, len=%d, need=%d", frame.PC, len(frame.Code), n)
	}
	return frame.Code[frame.PC : frame.PC+n], nil
}

func (vm *DebuggerVM) RequireStack(n int) error {
	frame := vm.currentFrame()
	if frame == nil {
		return fmt.Errorf("no execution frame")
	}

	if frame.Stack.Len() < n {
		return fmt.Errorf("stack underflow: need %d, have %d", n, frame.Stack.Len())
	}
	return nil
}

func (vm *DebuggerVM) Pop2() (*uint256.Int, *uint256.Int, error) {
	frame := vm.currentFrame()
	if frame == nil {
		return nil, nil, fmt.Errorf("no execution frame")
	}

	a, err := frame.Stack.Pop()
	if err != nil {
		return nil, nil, err
	}
	b, err := frame.Stack.Pop()
	if err != nil {
		return nil, nil, err
	}
	return a, b, nil
}

func (vm *DebuggerVM) Pop3() (*uint256.Int, *uint256.Int, *uint256.Int, error) {
	frame := vm.currentFrame()
	if frame == nil {
		return nil, nil, nil, fmt.Errorf("no execution frame")
	}

	a, err := frame.Stack.Pop()
	if err != nil {
		return nil, nil, nil, err
	}
	b, err := frame.Stack.Pop()
	if err != nil {
		return nil, nil, nil, err
	}
	c, err := frame.Stack.Pop()
	if err != nil {
		return nil, nil, nil, err
	}
	return a, b, c, nil
}

func (vm *DebuggerVM) PushUint64(u uint64) error {
	frame := vm.currentFrame()
	if frame == nil {
		return fmt.Errorf("no execution frame")
	}
	return frame.Stack.Push(new(uint256.Int).SetUint64(u))
}

func (vm *DebuggerVM) Push(x *uint256.Int) error {
	frame := vm.currentFrame()
	if frame == nil {
		return fmt.Errorf("no execution frame")
	}

	return frame.Stack.Push(x)
}

func (vm *DebuggerVM) ReadStorage(slot *uint256.Int) *uint256.Int {
	key := fmt.Sprintf("%064x", slot) // 32-byte hex string
	val := vm.Storage[key]
	if val == nil {
		return new(uint256.Int) // default zero
	}
	return new(uint256.Int).Set(val)
}

func (vm *DebuggerVM) WriteStorage(slot *uint256.Int, value *uint256.Int) {
	key := fmt.Sprintf("%064x", slot)
	vm.Storage[key] = new(uint256.Int).Set(value)
}

func (vm *DebuggerVM) ReadTransientStorage(slot *uint256.Int) *uint256.Int {
	key := fmt.Sprintf("%064x", slot) // 32-byte hex string
	val := vm.TransientStorage[key]
	if val == nil {
		return new(uint256.Int) // default zero
	}
	return new(uint256.Int).Set(val)
}

func (vm *DebuggerVM) WriteTransientStorage(slot *uint256.Int, value *uint256.Int) {
	key := fmt.Sprintf("%064x", slot)
	vm.TransientStorage[key] = new(uint256.Int).Set(value)
}

func (vm *DebuggerVM) ClearTransientStorage() {
	vm.TransientStorage = make(map[string]*uint256.Int)
}

func (vm *DebuggerVM) PushBytes(data []byte) error {
	bi := new(uint256.Int).SetBytes(data)
	return vm.Push(bi)
}

func (vm *DebuggerVM) UseGas(amount uint64) error {
	if vm.Context.Gas < amount {
		return fmt.Errorf("out of gas")
	}
	vm.Context.Gas -= amount
	return nil
}

func (vm *DebuggerVM) IsValidPC(pc uint64) bool {
	frame := vm.currentFrame()
	if frame == nil || frame.CodeMetadata == nil {
		return false
	}
	_, ok := frame.CodeMetadata.ValidPC[pc]
	return ok
}

func (vm *DebuggerVM) IsJumpDest(pc uint64) bool {
	frame := vm.currentFrame()
	if frame == nil || frame.CodeMetadata == nil {
		return false
	}
	_, ok := frame.CodeMetadata.JumpDests[pc]
	return ok
}

func (vm *DebuggerVM) RequireContext() error {
	if vm.Context == nil {
		return fmt.Errorf("execution context not set")
	}
	return nil
}

func scanCodeMetadata(code []byte) *CodeMetadata {
	validPC := make(map[uint64]struct{})
	jumpDests := make(map[uint64]struct{})

	for pc := 0; pc < len(code); {
		validPC[uint64(pc)] = struct{}{}

		op := code[pc]
		if op == 0x5b { // JUMPDEST
			jumpDests[uint64(pc)] = struct{}{}
			pc++
		} else if op >= 0x60 && op <= 0x7f { // PUSH1 to PUSH32
			pushLen := int(op - 0x5f)
			pc += 1 + pushLen
		} else {
			pc++
		}
	}

	return &CodeMetadata{
		ValidPC:   validPC,
		JumpDests: jumpDests,
	}
}

// Frame management methods

// currentFrame returns the current execution frame
func (vm *DebuggerVM) currentFrame() *MessageFrame {
	if len(vm.frames) == 0 {
		return nil
	}
	return &vm.frames[len(vm.frames)-1]
}

// pushFrame adds a new execution frame
func (vm *DebuggerVM) pushFrame(frame MessageFrame) error {
	const maxCallDepth = 1024
	if len(vm.frames) >= maxCallDepth {
		return ErrCallDepthLimit
	}

	// Add new frame
	vm.frames = append(vm.frames, frame)
	return nil
}

// popFrame removes the current execution frame
func (vm *DebuggerVM) popFrame() error {
	if len(vm.frames) <= 1 {
		return fmt.Errorf("cannot pop the root frame")
	}

	// Save return data from current frame
	current := vm.currentFrame()
	if current != nil {
		vm.lastReturnData = current.ReturnData
	}

	// Remove current frame
	vm.frames = vm.frames[:len(vm.frames)-1]
	return nil
}

// CallDepth returns the current call depth
func (vm *DebuggerVM) CallDepth() int {
	return len(vm.frames)
}

// ReturnData returns the return data from the last call
func (vm *DebuggerVM) ReturnData() []byte {
	if len(vm.lastReturnData) == 0 {
		return nil
	}
	return vm.lastReturnData
}

// ReturnDataSize returns the size of return data from the last call
func (vm *DebuggerVM) ReturnDataSize() *uint256.Int {
	return uint256.NewInt(uint64(len(vm.lastReturnData)))
}

// PushFrame adds a new execution frame (public method for opcodes)
func (vm *DebuggerVM) PushFrame(frame MessageFrame) error {
	return vm.pushFrame(frame)
}

// PopFrame removes the current execution frame (public method for opcodes)
func (vm *DebuggerVM) PopFrame() error {
	return vm.popFrame()
}

// ScanCodeMetadata is a public wrapper for scanCodeMetadata
func ScanCodeMetadata(code []byte) *CodeMetadata {
	return scanCodeMetadata(code)
}

// CurrentFrame returns the current execution frame (public method for opcodes)
func (vm *DebuggerVM) CurrentFrame() *MessageFrame {
	return vm.currentFrame()
}

// Properties that delegate to current frame

// Stack returns the current frame's stack
func (vm *DebuggerVM) Stack() *Stack {
	frame := vm.currentFrame()
	if frame == nil {
		return nil
	}
	return frame.Stack
}

// Memory returns the current frame's memory
func (vm *DebuggerVM) Memory() *Memory {
	frame := vm.currentFrame()
	if frame == nil {
		return nil
	}
	return frame.Memory
}

// PC returns the current frame's program counter
func (vm *DebuggerVM) PC() uint64 {
	frame := vm.currentFrame()
	if frame == nil {
		return 0
	}
	return frame.PC
}

// Code returns the current frame's code
func (vm *DebuggerVM) Code() []byte {
	frame := vm.currentFrame()
	if frame == nil {
		return nil
	}
	return frame.Code
}

// SetPC sets the current frame's program counter
func (vm *DebuggerVM) SetPC(pc uint64) {
	frame := vm.currentFrame()
	if frame != nil {
		frame.PC = pc
	}
}

// ExecuteCall executes the current frame until completion or revert
func (vm *DebuggerVM) ExecuteCall() error {
	// Save the stopped state and reset it for the call execution
	originalStopped := vm.Stopped
	vm.Stopped = false

	frame := vm.currentFrame()
	if frame == nil {
		vm.Stopped = originalStopped
		return fmt.Errorf("no execution frame")
	}

	// Execute until we hit RETURN, REVERT, or an error
	for !vm.Stopped && int(frame.PC) < len(frame.Code) {
		// Check if we're at a valid PC
		if _, ok := frame.CodeMetadata.ValidPC[frame.PC]; !ok {
			vm.Stopped = originalStopped
			return fmt.Errorf("invalid PC: 0x%x (likely inside PUSH immediate)", frame.PC)
		}

		err := vm.Step()
		if err != nil {
			// Restore original stopped state
			vm.Stopped = originalStopped
			return err
		}

		// Refresh frame reference after step (in case of frame changes)
		frame = vm.currentFrame()
		if frame == nil {
			vm.Stopped = originalStopped
			return fmt.Errorf("execution frame disappeared")
		}

		// Check if execution completed normally
		if frame.PC >= uint64(len(frame.Code)) {
			vm.Stopped = true
			break
		}
	}

	// Restore original stopped state
	vm.Stopped = originalStopped
	return nil
}
