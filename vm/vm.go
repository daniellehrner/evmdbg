package vm

import (
	"fmt"
	"math/big"
)

type Handler interface {
	Execute(vm *DebuggerVM) error
}
type HandlerGetter func(b byte) Handler

type DebuggerVM struct {
	Code    []byte
	PC      uint64
	Stack   *Stack
	Memory  *Memory
	Storage map[string]*big.Int
	Stopped bool

	ReturnValue []byte
	Reverted    bool
	Logs        []LogEntry

	Context       *ExecutionContext
	CodeMetadata  *CodeMetadata
	HandlerGetter HandlerGetter
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
	Value   *big.Int

	CallData []byte
	GasPrice *big.Int
	Gas      uint64
	Balance  *big.Int

	Block *BlockContext
}

type BlockContext struct {
	Coinbase   [20]byte
	Timestamp  uint64
	Number     uint64
	Difficulty *big.Int
	GasLimit   uint64
	ChainID    *big.Int
}

type CodeMetadata struct {
	ValidPC   map[uint64]struct{}
	JumpDests map[uint64]struct{}
}

func NewDebuggerVM(code []byte, hg HandlerGetter) *DebuggerVM {
	return &DebuggerVM{
		Code:          code,
		Stack:         NewStack(),
		Memory:        NewMemory(),
		HandlerGetter: hg,
		CodeMetadata:  scanCodeMetadata(code),
	}
}

func (vm *DebuggerVM) Step() error {
	if vm.Stopped || int(vm.PC) >= len(vm.Code) {
		vm.Stopped = true
		return nil
	}

	op := vm.Code[vm.PC]
	vm.PC++

	handler := vm.HandlerGetter(op)
	if handler == nil {
		return fmt.Errorf("unsupported opcode: 0x%x", op)
	}

	return handler.Execute(vm)
}

func (vm *DebuggerVM) RunUntil(breakpoints map[uint64]struct{}) error {
	for {
		if vm.Stopped || int(vm.PC) >= len(vm.Code) {
			vm.Stopped = true
			return nil
		}

		if _, ok := breakpoints[vm.PC]; ok {
			return nil // reached a breakpoint
		}

		// Only execute at valid PC (not in PUSH immediate)
		if _, ok := vm.CodeMetadata.ValidPC[vm.PC]; !ok {
			return fmt.Errorf("invalid PC: 0x%x (likely inside PUSH immediate)", vm.PC)
		}

		err := vm.Step()
		if err != nil {
			return err
		}
	}
}

func (vm *DebuggerVM) ReadCodeByte(offset uint64) (byte, error) {
	pos := vm.PC + offset
	if int(pos) >= len(vm.Code) {
		return 0, fmt.Errorf("code out of bounds at PC + %d", offset)
	}
	return vm.Code[pos], nil
}

func (vm *DebuggerVM) AdvancePC(n uint64) {
	vm.PC += n
}

func (vm *DebuggerVM) ReadCodeSlice(n uint64) ([]byte, error) {
	if vm.PC+n > uint64(len(vm.Code)) {
		return nil, fmt.Errorf("code out of bounds: PC=%d, len=%d, need=%d", vm.PC, len(vm.Code), n)
	}
	return vm.Code[vm.PC : vm.PC+n], nil
}

func (vm *DebuggerVM) RequireStack(n int) error {
	if vm.Stack.Len() < n {
		return fmt.Errorf("stack underflow: need %d, have %d", n, vm.Stack.Len())
	}
	return nil
}

func (vm *DebuggerVM) Pop2() (*big.Int, *big.Int, error) {
	a, err := vm.Stack.Pop()
	if err != nil {
		return nil, nil, err
	}
	b, err := vm.Stack.Pop()
	if err != nil {
		return nil, nil, err
	}
	return a, b, nil
}

func (vm *DebuggerVM) Pop3() (*big.Int, *big.Int, *big.Int, error) {
	a, err := vm.Stack.Pop()
	if err != nil {
		return nil, nil, nil, err
	}
	b, err := vm.Stack.Pop()
	if err != nil {
		return nil, nil, nil, err
	}
	c, err := vm.Stack.Pop()
	if err != nil {
		return nil, nil, nil, err
	}
	return a, b, c, nil
}

func (vm *DebuggerVM) PushUint64(u uint64) error {
	return vm.Stack.Push(new(big.Int).SetUint64(u))
}

func (vm *DebuggerVM) Push(x *big.Int) error {
	err := vm.Stack.Push(x)
	if err != nil {
		return err
	}

	return nil
}

func (vm *DebuggerVM) ReadStorage(slot *big.Int) *big.Int {
	key := fmt.Sprintf("%064x", slot) // 32-byte hex string
	val := vm.Storage[key]
	if val == nil {
		return new(big.Int) // default zero
	}
	return new(big.Int).Set(val)
}

func (vm *DebuggerVM) WriteStorage(slot *big.Int, value *big.Int) {
	key := fmt.Sprintf("%064x", slot)
	vm.Storage[key] = new(big.Int).Set(value)
}

func (vm *DebuggerVM) PushBytes(data []byte) error {
	bi := new(big.Int).SetBytes(data)
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
	_, ok := vm.CodeMetadata.ValidPC[pc]
	return ok
}

func (vm *DebuggerVM) IsJumpDest(pc uint64) bool {
	_, ok := vm.CodeMetadata.JumpDests[pc]
	return ok
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
