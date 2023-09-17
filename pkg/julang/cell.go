package julang

import (
	"fmt"
)

type CellType uint8

const (
	_ CellType = iota
	CellTypeQuotation
	CellTypeUint8
	CellTypeBytes
)

var CellTypeLabels = [...]string{
	CellTypeQuotation: "quotation",
	CellTypeUint8:     "uint8",
	CellTypeBytes:     "bytes",
}

func (c CellType) String() string { return CellTypeLabels[c] }

type Cell struct {
	Type CellType
	Data []byte
}

func NewCell(typ CellType, data ...byte) Cell { return Cell{Type: typ, Data: data} }

func NewCellUint8(v uint8) Cell { return NewCell(CellTypeUint8, v) }
func (c Cell) AsUint8() uint8   { c.mustHaveType(CellTypeUint8); return c.Data[0] }

func NewCellBytes(v []byte) Cell { return NewCell(CellTypeBytes, v...) }
func (c Cell) AsBytes() []byte   { c.mustHaveType(CellTypeBytes); return c.Data }

func NewCellQuotation(v string) Cell { return NewCell(CellTypeQuotation, []byte(v)...) }
func (c Cell) AsQuotation() string   { c.mustHaveType(CellTypeQuotation); return string(c.Data) }

func (c Cell) mustHaveType(typ CellType) {
	if c.Type != typ {
		panic(fmt.Errorf("cell should be of type %q instead of %q", typ, c.Type))
	}
}
