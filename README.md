# evmdbg

A minimal Ethereum Virtual Machine (EVM) execution engine and debugger written in Go. Designed for clarity, simplicity, and fast iteration.

## Features

 - Core EVM opcode support (arithmetic, memory, storage, control flow, etc.)
 - In-memory stack, memory, and storage model
 - Simulates EVM bytecode with step-by-step introspection
 - Clean, extensible opcode handler architecture
 - Usable as a CLI tool or embeddable in debugging UIs

## Status

Implemented opcodes: ~80–90% of commonly used opcodes.

Missing: `CALL`, `CALLCODE`, `DELEGATECALL`, `STATICCALL`, `CREATE`, `CREATE2`, `EXTCODE*`, `BLOCKHASH`, etc.

The full set of implemented opcodes can be found in the [vm/opcodes.go](vm/opcodes.go).

## Using as a Library

```go
package main

import (
	"github.com/daniellehrner/evmdbg/vm"
	"github.com/daniellehrner/evmdbg/vm/opcode_handlers"
)

func NewDebuggerVM(code []byte) *vm.DebuggerVM {
	d := vm.NewDebuggerVM(code, opcode_handlers.GetHandler)
	// Optionally configure execution context:
	d.Context.CallData = []byte{ /* ... */ }
	d.Context.Address = [20]byte{ /* ... */ }
	return d
}
```

Execute bytecode:

```go
for !d.Stopped {
	err := d.Step()
	if err != nil {
		// Handle execution errors
		break
	}
}
```

Inspect state after or during execution:
 - d.Stack for current stack state
 - d.Memory and d.Storage
 - d.ReturnValue after RETURN
 - d.Logs for emitted logs

## Goals
 - Educational: a readable, hackable EVM core for learning and experimentation
 - Tooling: a base for debuggers, linters, language servers, or smart contract playgrounds

## Future Work
 - Support for CALL and other multi-context opcodes
 - WASM-compatible build for browser usage
 - Source mapping and symbolic variable tracking
 - Basic gas accounting