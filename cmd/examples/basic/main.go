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
