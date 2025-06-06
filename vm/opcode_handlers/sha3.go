package opcode_handlers

import (
	"github.com/daniellehrner/evmdbg/vm"
	"golang.org/x/crypto/sha3"
	"hash"
	"math/big"
)

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func newKeccakState() KeccakState {
	return sha3.NewLegacyKeccak256().(KeccakState)
}

type Sha3OpCode struct{}

func (*Sha3OpCode) Execute(v *vm.DebuggerVM) error {
	// The SHA3 opcode requires two values on the stack: offset and size.
	if err := v.RequireStack(2); err != nil {
		return err
	}

	// Pop the top two items from the stack.
	offset, size, err := v.Pop2()
	if err != nil {
		return err
	}

	// Read the specified memory range and compute the SHA3 hash.
	data := v.Memory.Read(int(offset.Uint64()), int(size.Uint64()))

	// calculate the SHA3 hash of the data
	h := keccak256(data)

	return v.Push(new(big.Int).SetBytes(h))
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
