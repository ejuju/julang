package julang

import (
	"errors"
	"testing"
)

func TestStack(t *testing.T) {
	t.Run("can push and pop one", func(t *testing.T) {
		s := NewStack(1)
		err := s.Push(NewCellUint8(0x69))
		if err != nil {
			t.Fatal(err)
		}
		c, err := s.Pop()
		if err != nil {
			t.Fatal(err)
		}
		if c.AsUint8() != 0x69 {
			t.Fatalf("got %v instead of [0x69]", c)
		}
	})

	t.Run("fails on stack overflow and underflow", func(t *testing.T) {
		s := NewStack(1)

		err := s.Push(NewCellUint8(0x69))
		if err != nil {
			panic(err)
		}
		err = s.Push(NewCellUint8(0x69))
		if !errors.Is(err, ErrStackOverflow) {
			t.Fatalf("got %q instead of stack overflow", err)
		}

		_, err = s.Pop()
		if err != nil {
			panic(err)
		}
		_, err = s.Pop()
		if !errors.Is(err, ErrStackUnderflow) {
			t.Fatalf("got %q instead of stack underflow", err)
		}
	})
}
