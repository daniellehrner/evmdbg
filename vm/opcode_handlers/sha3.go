package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"golang.org/x/crypto/sha3"
	"hash"
	"math/big"
)

type Sha3OpCode struct{}

func (*Sha3OpCode) Execute(v *vm.DebuggerVM) error {
	if err := v.RequireStack(2); err != nil {
		return err
	}
	offset, size, err := v.Pop2()
	if err != nil {
		return err
	}
	data := v.Memory.Read(int(offset.Uint64()), int(size.Uint64()))
	h := keccak256(data)
	return v.Push(new(big.Int).SetBytes(h))
}

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := newKeccakState()
	for _, part := range data {
		d.Write(part)
	}
	d.Read(b)
	return b
}

func newKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}
