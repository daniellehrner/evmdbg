package main

import (
	"fmt"
	"math/big"

	"github.com/daniellehrner/evmdbg/vm"
	"github.com/daniellehrner/evmdbg/vm/opcode_handlers"
)

func main() {
	// Program: PUSH1 0x03 PUSH1 0x05 ADD PUSH1 0x00 MSTORE PUSH1 0x20 PUSH1 0x00 RETURN
	code := []byte{
		0x60, 0x03, // PUSH1 3
		0x60, 0x05, // PUSH1 5
		0x01,       // ADD
		0x60, 0x00, // PUSH1 0 (offset)
		0x52,       // MSTORE
		0x60, 0x20, // PUSH1 32 (size)
		0x60, 0x00, // PUSH1 0 (offset)
		0xf3, // RETURN
	}

	evmdbg := vm.NewDebuggerVM(code, opcode_handlers.GetHandler)

	for !evmdbg.Stopped {
		if err := evmdbg.Step(); err != nil {
			fmt.Println("VM Error:", err)
			break
		}
	}

	if evmdbg.Reverted {
		fmt.Println("Execution reverted. Output:")
	} else {
		fmt.Println("Execution successful. Output:")
	}

	fmt.Printf("0x%x\n", evmdbg.ReturnValue)
	val := new(big.Int).SetBytes(evmdbg.ReturnValue)
	fmt.Printf("Value as decimal: %s\n", val.String())
}
