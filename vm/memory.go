package vm

import "math/big"

type Memory struct {
	data []byte
}

func NewMemory() *Memory {
	return &Memory{data: make([]byte, 0)}
}

func (m *Memory) Data() []byte {
	return m.data
}

func (m *Memory) Write(offset int, b []byte) {
	end := offset + len(b)
	if end > len(m.data) {
		newMem := make([]byte, end)
		copy(newMem, m.data)
		m.data = newMem
	}
	copy(m.data[offset:], b)
}

func (m *Memory) Read(offset int, size int) []byte {
	end := offset + size
	if end > len(m.data) {
		return append([]byte{}, m.data[offset:]...)
	}
	return m.data[offset:end]
}

func (m *Memory) ReadWord(offset uint64) *big.Int {
	b := m.Read(int(offset), 32) // always 32 bytes
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b) // right-align if b is shorter
	return new(big.Int).SetBytes(padded)
}

func (m *Memory) WriteWord(offset uint64, value *big.Int) {
	b := value.Bytes()
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b) // right-align
	m.Write(int(offset), padded)
}
