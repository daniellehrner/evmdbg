package opcode_handlers

import (
	"fmt"
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type PushNOpCode struct {
	N int
}

func (p *PushNOpCode) Execute(v *vm.DebuggerVM) error {
	if p.N < 1 || p.N > 32 {
		return fmt.Errorf("invalid PUSHN size: %d", p.N)
	}
	data := make([]byte, p.N)
	for i := 0; i < p.N; i++ {
		b, err := v.ReadCodeByte(uint64(i))
		if err != nil {
			return fmt.Errorf("PUSH%d: %w", p.N, err)
		}
		data[i] = b
	}
	v.AdvancePC(uint64(p.N))
	val := new(big.Int).SetBytes(data)
	return v.Stack.Push(val)
}
