package matrix

import "github.com/stewi1014/modular"

func Identity[T modular.Number]() Matrix[T] {
	return ident[T]{}
}

type ident[T modular.Number] struct{}

func (ident[T]) At(row, col int) T {
	if row == col {
		return 1
	}
	return 0
}

func (ident[T]) NumRows() int {
	return -1
}

func (ident[T]) NumCols() int {
	return -1
}
