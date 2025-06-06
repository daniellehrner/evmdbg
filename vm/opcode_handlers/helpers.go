package opcode_handlers

import "math/big"

// eqivalent to 2^256 - 1
var uint256Mask = new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

var uint256 = new(big.Int).Lsh(big.NewInt(1), 256)
var uint256Half = new(big.Int).Rsh(uint256, 1)

func toSigned(x *big.Int, half *big.Int) *big.Int {
	if x.Cmp(half) >= 0 {
		return new(big.Int).Sub(x, half.Mul(half, big.NewInt(2)))
	}
	return new(big.Int).Set(x)
}

func padTo256Bytes(b []byte) []byte {
	if len(b) >= 32 {
		return b[len(b)-32:]
	}
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}
