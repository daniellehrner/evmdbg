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
