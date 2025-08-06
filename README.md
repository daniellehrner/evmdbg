# evmdbg

A minimal Ethereum Virtual Machine (EVM) execution engine and debugger written in Go. Designed for clarity, simplicity, and fast iteration.

> ⚠️ **Experimental Project**
>
> This project is still under active development and is considered experimental.
> It is **not guaranteed to be correct**, complete, or secure.  
> Use at your own risk, especially for production or security-critical purposes.

## Features

 - Core EVM opcode support (arithmetic, memory, storage, control flow, calls)
 - Call opcode support (`CALL`, `CALLCODE`, `DELEGATECALL`, `STATICCALL`)
 - Frame-based execution model with proper call context switching
 - Static call protection preventing state changes in static contexts
 - In-memory stack, memory, and storage model
 - Simulates EVM bytecode with step-by-step introspection
 - Clean, extensible opcode handler architecture
 - Usable as a CLI tool or embeddable in debugging UIs

## Status

Implemented opcodes: ~90-95% of commonly used opcodes.

**Recently implemented**: `CALL`, `CALLCODE`, `DELEGATECALL`, `STATICCALL`

**Missing**: `CREATE`, `CREATE2`, `EXTCODE*`, `BLOCKHASH`, some precompiled contracts

The full set of implemented opcodes can be found in the [vm/opcodes.go](vm/opcodes.go).

## Using as a Library

Here are some examples of how to use the `evmdbg` package in your Go applications.

### Step-by-Step Execution

You can step through the execution one instruction at a time:
```go
package main

import (
    "fmt"
    "github.com/daniellehrner/evmdbg/evmdbg"
)

func main() {
    code := []byte{
        0x60, 0x02, // PUSH1 0x02
        0x60, 0x03, // PUSH1 0x03
        0x01, // ADD
    }

    v := evmdbg.CreateDebuggerVM(code)

    for !v.Stopped {
        fmt.Printf("PC: %d, Stack: %s\n", v.PC(), v.Stack().String())
        err := v.Step()
        if err != nil {
            fmt.Printf("Execution error: %v\n", err)
            break
        }
    }

    fmt.Printf("Final Stack: %s\n", v.Stack().String())
}

```

### Run Until Breakpoint

You can execute the VM until it reaches a specific program counter (PC) using RunUntil. This avoids stepping into the middle of immediates (e.g., PUSH32 payloads):

```go
package main

import (
    "fmt"
    "github.com/daniellehrner/evmdbg/evmdbg"
)

func main() {
    code := []byte{
        0x60, 0x01, // PUSH1 0x01
        0x60, 0x02, // PUSH1 0x02
        0x01,       // ADD
        0x60, 0x00, // PUSH1 0x00
        0x52,       // MSTORE
        0x60, 0x20, // PUSH1 0x20
        0x60, 0x00, // PUSH1 0x00
        0xf3, // RETURN
    }

    v := evmdbg.CreateDebuggerVM(code)

    breakpoints := map[uint64]struct{}{
        5: {}, // break before MSTORE (PC at 0x05)
    }

    err := v.RunUntil(breakpoints)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Stopped at PC: %d\n", v.PC())
    fmt.Printf("Stack: %s\n", v.Stack().String())
}

```

### Inspecting State During and After Execution

You can interact with the virtual machine’s components directly:

```go
package main

import (
    "fmt"
    "github.com/daniellehrner/evmdbg/evmdbg"
)

func main() {
    // Example program: Store 0xdeadbeef at memory offset 0x00, return it
    code := []byte{
        0x7f, 0xde, 0xad, 0xbe, 0xef, 0x00, 0x00, 0x00, 0x00, // PUSH32 0xdeadbeef...
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
        0x60, 0x00, // PUSH1 0x00
        0x52,       // MSTORE
        0x60, 0x20, // PUSH1 0x20
        0x60, 0x00, // PUSH1 0x00
        0xf3, // RETURN
    }

    vm := evmdbg.CreateDebuggerVM(code)

    // Run the VM
    for !vm.Stopped {
        if err := vm.Step(); err != nil {
            fmt.Printf("Error: %v\n", err)
            break
        }

        // Dump stack after each step
        fmt.Print("Stack:")
        if vm.Stack().Len() == 0 {
            fmt.Println(" <empty>")
        } else {
            for i := 0; i < vm.Stack().Len(); i++ {
                val, _ := vm.Stack().Peek(i)
                fmt.Printf(" [%d]: 0x%x", i, val)
            }
            fmt.Println()
        }
    }

    // Inspect stack
    fmt.Println("Final Stack:")
    for i := 0; i < vm.Stack().Len(); i++ {
        val, _ := vm.Stack().Peek(i)
        fmt.Printf("  [%d]: 0x%x\n", i, val)
    }

    // Inspect memory (first 32 bytes)
    fmt.Println("Memory [0x00..0x20]:", vm.Memory().Read(0, 32))

    // Inspect return value
    fmt.Println("Return value:", vm.ReturnValue)

    // Inspect logs (if any LOG opcodes were used)
    if len(vm.Logs) > 0 {
        fmt.Println("Logs:")
        for _, log := range vm.Logs {
            fmt.Printf("  Address: %x\n", log.Address)
            for i, topic := range log.Topics {
                fmt.Printf("  Topic %d: %x\n", i, topic)
            }
            fmt.Printf("  Data: %x\n", log.Data)
        }
    }
}
```

## Goals
- Tooling: a base for debuggers, linters, language servers, or smart contract playgrounds
- Educational: a readable, hackable EVM core for learning and experimentation

## CLI Examples

The project includes several example programs demonstrating different use cases:

### Basic Example
```bash
go run cmd/examples/basic/main.go
```
Simple step-by-step execution of basic arithmetic operations.

### Inspection Example  
```bash
go run cmd/examples/inspect/main.go
```
Demonstrates full state introspection during execution (stack, memory, storage).

### Run Until Example
```bash
go run cmd/examples/run_until/main.go
```
Shows breakpoint-style execution using `RunUntil()` method.

## Development

### Building and Testing
```bash
# Build the project
go build ./...

# Run all tests
go test ./...

# Run tests for specific package
go test ./vm/opcode_handlers

# Run tests with verbose output
go test -v ./vm/opcode_handlers
```

### Architecture Overview

The project is structured around a clean handler pattern:

- **`vm/`**: Core VM implementation with stack, memory, and execution logic
- **`vm/opcode_handlers/`**: Individual opcode implementations following the `Handler` interface
- **`evmdbg/`**: Public API wrapper for easy library usage
- **`cmd/examples/`**: Example programs demonstrating various use cases

Each opcode is implemented as a separate handler struct, making the codebase modular and easy to extend.

## Contributing

When adding new opcodes:

1. Add the opcode constant to `vm/opcodes.go`
2. Create a handler file in `vm/opcode_handlers/{opcode}.go` implementing the `Handler` interface
3. Register the handler in `vm/opcode_handlers/handlers.go`
4. Add comprehensive tests in `vm/opcode_handlers/{opcode}_test.go`

Tests should follow the pattern:
1. Create bytecode with the opcode and operands
2. Initialize VM and step through execution
3. Verify final state (stack, memory, storage, logs)

## Future Work
 - Support for `CREATE`, `CREATE2`, and `EXTCODE*` opcodes
 - Source mapping and symbolic variable tracking
 - Basic gas accounting
 - Precompiled contract support

## License

This project is experimental and provided as-is for educational and development purposes.