package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SGTOpCode struct{}

func (*SGTOpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	a, b, err := v.Pop2()
	if err != nil {
		return err
	}
	if toSigned(a, uint256Half).Cmp(toSigned(b, uint256Half)) > 0 {
		return v.Push(big.NewInt(1))
	}
	return v.Push(big.NewInt(0))
}
