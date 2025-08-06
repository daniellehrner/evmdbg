# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`evmdbg` is a minimal Ethereum Virtual Machine (EVM) execution engine and debugger written in Go. It's designed for tooling development and EVM bytecode analysis. The project provides step-by-step EVM execution with full state introspection capabilities.

## Development Commands

### Building and Running
```bash
# Build the project
go build ./...

# Run tests
go test ./...

# Run tests for a specific package
go test ./vm/opcode_handlers

# Run example programs
go run cmd/examples/basic/main.go
go run cmd/examples/inspect/main.go
go run cmd/examples/run_until/main.go
```

### Testing
- Individual opcode tests are in `vm/opcode_handlers/*_test.go`
- Tests follow the pattern: create VM with bytecode, step through execution, verify final state
- Use `go test -v ./vm/opcode_handlers` for verbose test output

## Architecture

### Core Components

**VM Package (`vm/`)**
- `vm.go`: Main `DebuggerVM` struct with execution logic, step-by-step execution, breakpoint support
- `stack.go`: EVM stack implementation with 256-bit integer operations
- `memory.go`: EVM memory model with dynamic expansion
- `opcodes.go`: OpCode constants matching EVM specification

**Opcode Handlers (`vm/opcode_handlers/`)**
- `handlers.go`: Central opcode-to-handler mapping and registry
- Individual opcode implementations: `add.go`, `mul.go`, `push_1_32.go`, etc.
- Each handler implements the `Handler` interface with `Execute(vm *DebuggerVM) error`
- Test files follow pattern `{opcode}_test.go`

**Public API (`evmdbg/`)**
- `evmdbg.go`: Simple wrapper providing `CreateDebuggerVM(code []byte)` convenience function

### Key Design Patterns

1. **Handler Pattern**: Each opcode is implemented as a separate handler struct implementing the `Handler` interface
2. **Metadata Scanning**: Code is pre-scanned to identify valid PC positions and JUMPDEST locations
3. **State Introspection**: Full VM state (stack, memory, storage, logs) is accessible at any execution point
4. **Breakpoint Support**: `RunUntil()` method allows execution until specific PC addresses

### Execution Context

The VM supports execution context including:
- Caller, origin, and contract addresses
- Call data and call value
- Block context (coinbase, timestamp, number, difficulty, gas limit, chain ID)
- Gas tracking and storage operations

### Memory Model

- **Stack**: 256-bit integers with standard EVM operations (push, pop, peek, dup, swap)
- **Memory**: Byte-addressable with automatic expansion, accessed via MLOAD/MSTORE
- **Storage**: Key-value mapping using 256-bit keys and values, accessed via SLOAD/SSTORE

## Implementation Status

**Implemented**: ~80-90% of commonly used opcodes including:
- Arithmetic: ADD, SUB, MUL, DIV, MOD, EXP, etc.
- Bitwise: AND, OR, XOR, NOT, BYTE, SHL/SHR/SAR
- Comparison: LT, GT, EQ, ISZERO
- Stack: PUSH0-PUSH32, DUP1-DUP16, SWAP1-SWAP16, POP
- Memory: MLOAD, MSTORE, MSIZE
- Storage: SLOAD, SSTORE  
- Control flow: JUMP, JUMPI, JUMPDEST, PC
- Context: ADDRESS, CALLER, CALLVALUE, CALLDATALOAD, etc.
- Block: COINBASE, TIMESTAMP, NUMBER, DIFFICULTY, etc.
- Logging: LOG0-LOG4
- Termination: STOP, RETURN, REVERT

**Missing**: CALL family, CREATE family, EXTCODE* opcodes, BLOCKHASH

## Working with the Codebase

### Adding New Opcodes
1. Add opcode constant to `vm/opcodes.go`
2. Create handler file `vm/opcode_handlers/{opcode}.go` implementing `Handler` interface
3. Register handler in `vm/opcode_handlers/handlers.go`
4. Add comprehensive tests in `vm/opcode_handlers/{opcode}_test.go`

### Testing Pattern
Tests typically:
1. Create bytecode array with the opcode and operands
2. Initialize VM with `vm.NewDebuggerVM(code, GetHandler)`
3. Execute with stepping loop until `vm.Stopped`
4. Verify final state (stack, memory, storage, logs as appropriate)

### Debugging and Analysis
- Use `vm.Stack.String()` for stack visualization
- Access `vm.Memory.Read(offset, length)` for memory inspection
- Check `vm.Storage` map for storage state
- Examine `vm.Logs` for emitted events
- Use `vm.PC` and `vm.CodeMetadata` for execution tracking