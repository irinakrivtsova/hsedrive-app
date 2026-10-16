package print

import (
	"fmt"
	"io"
)

type Printable interface {
	ForEach(fn func(id int, name string) error) error
}

type Printer struct {
	source Printable
}

func NewPrinter(source Printable) *Printer {
	return &Printer{source: source}
}

func (p *Printer) Print(w io.Writer) error {
	return p.source.ForEach(func(id int, name string) error {
		_, err := fmt.Fprintf(w, "%s в наличии, инвентарный номер: %d\n", name, id)
		return err
	})
}
