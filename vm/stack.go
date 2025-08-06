package vm

import (
	"fmt"

	"github.com/holiman/uint256"
)

type Stack struct {
	data []*uint256.Int
}

func NewStack() *Stack {
	return &Stack{data: make([]*uint256.Int, 0, 1024)}
}

func (s *Stack) Push(x *uint256.Int) error {
	if len(s.data) >= 1024 {
		return fmt.Errorf("stack overflow")
	}
	s.data = append(s.data, new(uint256.Int).Set(x))
	return nil
}

func (s *Stack) Pop() (*uint256.Int, error) {
	n := len(s.data)
	if n == 0 {
		return nil, fmt.Errorf("stack underflow")
	}
	x := s.data[n-1]
	s.data = s.data[:n-1]
	return x, nil
}

func (s *Stack) Len() int {
	return len(s.data)
}

func (s *Stack) String() string {
	out := "["
	for i, x := range s.data {
		if i > 0 {
			out += " "
		}
		out += fmt.Sprintf("0x%x", x)
	}
	return out + "]"
}

func (s *Stack) Peek(n int) (*uint256.Int, error) {
	if n < 0 || n >= len(s.data) {
		return nil, fmt.Errorf("stack underflow on peek(%d): size=%d", n, len(s.data))
	}
	// Top of stack is the end of the slice
	index := len(s.data) - 1 - n
	return s.data[index], nil
}

func (s *Stack) Swap(n int) error {
	if n < 1 || n >= len(s.data) {
		return fmt.Errorf("stack underflow on swap(%d): size=%d", n, len(s.data))
	}
	top := len(s.data) - 1
	other := top - n
	s.data[top], s.data[other] = s.data[other], s.data[top]
	return nil
}
