package julang

import (
	"fmt"
	"strings"
)

type Dictionary struct{ words []Word }

func NewDictionary(words ...Word) *Dictionary { return &Dictionary{words: words} }

type Word struct {
	Name string
	Do   func(vm *VM) error
}

func (d *Dictionary) Append(w Word) { d.words = append(d.words, w) }
func (d *Dictionary) Get(name string) (Word, bool) {
	for i := len(d.words) - 1; i >= 0; i-- {
		w := d.words[i]
		if w.Name == name {
			return w, true
		}
	}
	return Word{}, false
}

var Builtins = []Word{
	{
		Name: "define",
		Do: func(vm *VM) error {
			c, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop quotation: %w", err)
			}
			if c.Type != CellTypeQuotation {
				return fmt.Errorf("got type %q instead of quotation", c.Type)
			}
			code := c.AsQuotation()

			c, err = vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop name: %w", err)
			}
			if c.Type != CellTypeBytes {
				return fmt.Errorf("got type %q instead of bytes", c.Type)
			}
			name := string(c.AsBytes())
			vm.d.Append(Word{Name: name, Do: func(vm *VM) error { return vm.Exec(strings.NewReader(code)) }})
			return nil
		},
	},
	{
		Name: "if",
		Do: func(vm *VM) error {
			c, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop false-quotation: %w", err)
			}
			if c.Type != CellTypeQuotation {
				return fmt.Errorf("got type %q instead of quotation", c.Type)
			}
			onFalse := c.AsQuotation()

			c, err = vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop true-quotation: %w", err)
			}
			if c.Type != CellTypeQuotation {
				return fmt.Errorf("got type %q instead of quotation", c.Type)
			}
			onTrue := c.AsQuotation()

			c, err = vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop boolean: %w", err)
			}
			if c.Type != CellTypeUint8 {
				return fmt.Errorf("got type %q instead of uint8", c.Type)
			}
			ok := c.AsUint8()
			if ok == 0 {
				return vm.Exec(strings.NewReader(onTrue))
			}
			return vm.Exec(strings.NewReader(onFalse))
		},
	},
	{
		Name: "loop",
		Do: func(vm *VM) error {
			c, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop callback: %w", err)
			}
			if c.Type != CellTypeQuotation {
				return fmt.Errorf("got type %q instead of quotation", c.Type)
			}
			callback := c.AsQuotation()

			c, err = vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop count: %w", err)
			}
			if c.Type != CellTypeUint8 {
				return fmt.Errorf("got type %q instead of uint8", c.Type)
			}
			count := c.AsUint8()

			for i := 0; i < int(count); i++ {
				err = vm.Exec(strings.NewReader(callback))
				if err != nil {
					return err
				}
			}
			return nil
		},
	},
	{
		Name: "noop",
		Do:   func(vm *VM) error { return nil },
	},
	{
		Name: "print",
		Do: func(vm *VM) error {
			c, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop bytes: %w", err)
			}
			_, err = vm.stdout.Write(append(c.Data, '\n'))
			if err != nil {
				return fmt.Errorf("write bytes: %w", err)
			}
			return nil
		},
	},
	{
		Name: "dup",
		Do: func(vm *VM) error {
			c, err := vm.s.Peek()
			if err != nil {
				return err
			}
			return vm.s.Push(c)
		},
	},
	{
		Name: "drop",
		Do: func(vm *VM) error {
			_, err := vm.s.Pop()
			if err != nil {
				return err
			}
			return nil
		},
	},
	{
		Name: "+",
		Do: func(vm *VM) error {
			b, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop b: %w", err)
			}
			a, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop a: %w", err)
			}
			err = vm.s.Push(NewCellUint8(a.AsUint8() + b.AsUint8()))
			if err != nil {
				return fmt.Errorf("push a+b: %w", err)
			}
			return nil
		},
	},
	{
		Name: "-",
		Do: func(vm *VM) error {
			b, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop b: %w", err)
			}
			a, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop a: %w", err)
			}
			err = vm.s.Push(NewCellUint8(a.AsUint8() - b.AsUint8()))
			if err != nil {
				return fmt.Errorf("push a-b: %w", err)
			}
			return nil
		},
	},
	{
		Name: "*",
		Do: func(vm *VM) error {
			b, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop b: %w", err)
			}
			a, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop a: %w", err)
			}
			err = vm.s.Push(NewCellUint8(a.AsUint8() * b.AsUint8()))
			if err != nil {
				return fmt.Errorf("push a*b: %w", err)
			}
			return nil
		},
	},
	{
		Name: "/",
		Do: func(vm *VM) error {
			b, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop b: %w", err)
			}
			a, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop a: %w", err)
			}
			err = vm.s.Push(NewCellUint8(a.AsUint8() / b.AsUint8()))
			if err != nil {
				return fmt.Errorf("push a/b: %w", err)
			}
			return nil
		},
	},
	{
		Name: "%",
		Do: func(vm *VM) error {
			b, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop b: %w", err)
			}
			a, err := vm.s.Pop()
			if err != nil {
				return fmt.Errorf("pop a: %w", err)
			}
			err = vm.s.Push(NewCellUint8(a.AsUint8() % b.AsUint8()))
			if err != nil {
				return fmt.Errorf("push a mod b: %w", err)
			}
			return nil
		},
	},
}
