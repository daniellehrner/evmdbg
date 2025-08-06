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
