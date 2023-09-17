package julang

import (
	"errors"
)

type Stack struct{ mem []Cell }

func NewStack(capacity int) *Stack {
	if capacity <= 0 {
		capacity = 4096
	}
	return &Stack{mem: make([]Cell, 0, capacity)}
}

var (
	ErrStackOverflow  = errors.New("stack overflow")
	ErrStackUnderflow = errors.New("stack underflow")
)

func (s *Stack) Length() int { return len(s.mem) }

func (s *Stack) Push(c Cell) error {
	if len(s.mem) == cap(s.mem) {
		return ErrStackOverflow
	}
	s.mem = append(s.mem, c)
	return nil
}

func (s *Stack) Pop() (Cell, error) {
	if len(s.mem) == 0 {
		return Cell{}, ErrStackUnderflow
	}
	b := s.mem[len(s.mem)-1]
	s.mem = s.mem[:len(s.mem)-1]
	return b, nil
}

func (s *Stack) Peek() (Cell, error) {
	if len(s.mem) == 0 {
		return Cell{}, ErrStackUnderflow
	}
	return s.mem[len(s.mem)-1], nil
}
