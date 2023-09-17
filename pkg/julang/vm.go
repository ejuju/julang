package julang

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type VM struct {
	s *Stack
	d *Dictionary

	// IO
	stdout, stderr io.Writer
	stdin          io.Reader
}

func NewVM(s *Stack, d *Dictionary) *VM {
	if s == nil {
		s = NewStack(0)
	}
	if d == nil {
		d = NewDictionary(Builtins...)
	}
	return &VM{
		s:      s,
		d:      d,
		stdout: os.Stdout,
		stderr: os.Stderr,
		stdin:  os.Stdin,
	}
}

func (vm *VM) Exec(src io.Reader) error {
	ts := NewTokenStream(src)
	for {
		tok, err := ts.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		switch tok.Type {
		default:
			panic(fmt.Errorf("unhandled token type %q", tok.Type))
		case TokenTypeQuotation:
			err = vm.s.Push(NewCellQuotation(tok.Value))
			if err != nil {
				return err
			}
		case TokenTypeWord:
			w, ok := vm.d.Get(tok.Value)
			if !ok {
				num, err := strconv.ParseUint(tok.Value, 10, 8)
				if err != nil {
					return fmt.Errorf("unknown word %q: %w", tok.Value, err)
				}
				err = vm.s.Push(NewCellUint8(uint8(num)))
				if err != nil {
					return err
				}
				continue
			}
			err = w.Do(vm)
			if err != nil {
				return err
			}
		case TokenTypeLiteralText:
			err = vm.s.Push(NewCellBytes([]byte(tok.Value)))
			if err != nil {
				return err
			}
		}
	}
	return nil
}
