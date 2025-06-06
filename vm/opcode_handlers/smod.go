package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"math/big"
)

type SModOpCode struct{}

func (*SModOpCode) Execute(vm *vm.DebuggerVM) error {
	if err := vm.RequireStack(2); err != nil {
		return err
	}
	a, b, err := vm.Pop2()
	if err != nil {
		return err
	}

	if b.Sign() == 0 {
		return vm.Push(big.NewInt(0))
	}

	sa := new(big.Int).Set(a)
	if sa.Cmp(uint256Half) >= 0 {
		sa.Sub(sa, uint256)
	}

	sb := new(big.Int).Set(b)
	if sb.Cmp(uint256Half) >= 0 {
		sb.Sub(sb, uint256)
	}

	mod := new(big.Int).Mod(sa, sb)
	if mod.Sign() < 0 {
		mod.Add(mod, uint256)
	}
	return vm.Push(mod)
}
