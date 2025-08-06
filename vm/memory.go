package vm

import (
	"github.com/holiman/uint256"
)

type Memory struct {
	data            []byte
	highestAccessed int // Tracks highest memory offset accessed (for MSIZE)
}

func NewMemory() *Memory {
	return &Memory{data: make([]byte, 0), highestAccessed: -1}
}

func (m *Memory) Data() []byte {
	return m.data
}

func (m *Memory) Write(offset int, b []byte) {
	if len(b) == 0 {
		return
	}

	end := offset + len(b)
	m.expandTo(end)

	// Update highest accessed offset
	if end-1 > m.highestAccessed {
		m.highestAccessed = end - 1
	}

	copy(m.data[offset:], b)
}

func (m *Memory) Read(offset int, size int) []byte {
	if size == 0 {
		return []byte{}
	}

	end := offset + size

	// Update highest accessed offset
	if end-1 > m.highestAccessed {
		m.highestAccessed = end - 1
	}

	// Expand memory if needed to ensure we can read the full range
	m.expandTo(end)

	// Now we can safely read the full range (memory is zero-initialized)
	return m.data[offset:end]
}

func (m *Memory) ReadWord(offset uint64) *uint256.Int {
	b := m.Read(int(offset), 32) // always 32 bytes, zeros if not written
	return new(uint256.Int).SetBytes(b)
}

func (m *Memory) WriteWord(offset uint64, value *uint256.Int) {
	b := value.Bytes()
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b) // right-align
	m.Write(int(offset), padded)
}

// wordAlignedSize returns the size rounded up to the next 32-byte boundary
func wordAlignedSize(size int) int {
	return ((size + 31) / 32) * 32
}

// expandTo expands memory to the given size using 32-byte word boundaries
func (m *Memory) expandTo(size int) {
	if size <= len(m.data) {
		return
	}

	newSize := wordAlignedSize(size)
	newMem := make([]byte, newSize)
	copy(newMem, m.data)
	m.data = newMem
}

// Size returns the EVM-compliant memory size (highest accessed offset + 1, word-aligned)
func (m *Memory) Size() int {
	if m.highestAccessed < 0 {
		return 0
	}
	// Round up to next 32-byte boundary from highest accessed + 1
	return wordAlignedSize(m.highestAccessed + 1)
}
