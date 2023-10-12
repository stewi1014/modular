// Package matrix aims to abstract matrix math just enough to allow significant optimisations to be transparantly acheived.
//
// TODO; eleborate: This package uses matricies to represent vectors.
package matrix

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/stewi1014/modular"
)

// assert that mgl32 implements our matrix interface
var _ Matrix[float64] = mgl64.Mat2{}
var _ Matrix[float64] = mgl64.Mat2x3{}
var _ Matrix[float64] = mgl64.Mat2x4{}
var _ Matrix[float64] = mgl64.Mat3x2{}
var _ Matrix[float64] = mgl64.Mat3{}
var _ Matrix[float64] = mgl64.Mat3x4{}
var _ Matrix[float64] = mgl64.Mat4x2{}
var _ Matrix[float64] = mgl64.Mat4x3{}
var _ Matrix[float64] = mgl64.Mat4{}
var _ Matrix[float64] = &mgl64.MatMxN{}

// Matrix is the centrepiece of this package.
//
// It aims to include matricies from packages like go-gl/mgl32 and go-gl/mgl64,
// while keeping the possibility for matricies to remain *undefined in memory*.
// The perfect example being the Identity matrix. The Identity() function does, quite literally, nothing.
// The returned matrix does not exist, and is generated "At" the given index by an single-branch inlineable function.
type Matrix[T modular.Number] interface {
	// At returns the number at the given index
	At(row, col int) T

	// NumRows returns the number of rows in the matrix
	NumRows() int

	// NumCols returns the number of colums in the matrix
	NumCols() int
}

// AbstractMatrix is an interface only implemented by matrices that do not have defined dimensions.
type AbstractMatrix[T modular.Number] interface {
	Evaluate(rows, cols int) Matrix[T]
}
