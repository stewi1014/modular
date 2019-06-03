package modular64

import (
	"errors"
	"math"

	"github.com/bmkessler/fastdiv"
)

// Error types
var (
	ErrBadModulo = errors.New("bad modulus")
	ErrBadIndex  = errors.New("bad index")
)

// NewIndexer creates a new Indexer.
//
// index must not be larger than 2**32, and modulus must be a normalised float.
//
// Special cases:
//		NewIndexer(m, 0) = panic(integer divide by zero)
//		NewIndexer(m, i > 2**32) = ErrBadIndex
//		NewIndexer(0, i) = ErrBadModulo
//		NewIndexer(±Inf, i) = ErrBadModulo
//		NewIndexer(NaN, i) = ErrBadModulo
//		NewIndexer(m, i) = ErrBadModulo for |m| < 2**-1022
func NewIndexer(modulus float64, index int) (Indexer, error) {
	mod := NewModulus(modulus)
	return mod.NewIndexer(index)
}

// NewIndexer creates a new indexer from the Modulus.
func (m Modulus) NewIndexer(index int) (Indexer, error) {
	if math.IsInf(m.mod, 0) || math.IsNaN(m.mod) || m.exp == 0 {
		return Indexer{}, ErrBadModulo
	}
	if index > (1<<32) || index < 1 {
		return Indexer{}, ErrBadIndex
	}

	modfr, _ := frexp(m.mod)
	r := modfr << fExponentBits //r - range; is shifted fExponentBits to get a little more
	rDivisor := r / uint64(index)
	return Indexer{
		Modulus: m,
		fdr:     fastdiv.NewUint64(rDivisor),
		r:       r,
		i:       index,
	}, nil
}

// Indexer provides a fast method for mapping a floating point modulus to a range of integers.
type Indexer struct {
	Modulus
	fdr fastdiv.Uint64
	r   uint64
	i   int
}

// Index indexes n.
//
// If n is NaN or ±Inf, it returns the index.
// Otherwise, it always satisfies 0 <= num < index
//
// Special cases:
//		Index(NaN) = index
//		Index(±Inf) = index
func (i Indexer) Index(n float64) int {
	if math.IsNaN(n) || math.IsInf(n, 0) || i.i == 0 {
		return i.i
	}

	nfr, nexp := frexp(n)
	var nr uint64
	switch {
	case n > i.mod:
		expdiff := nexp - i.exp
		nr = i.modExp(nfr, expdiff) << fExponentBits
	case n < -i.mod:
		expdiff := nexp - i.exp
		nr = i.modExp(nfr, expdiff) << fExponentBits
		if nr != 0 {
			nr = i.r - nr
		}
	case n < 0:
		nr = shiftSub(fExponentBits, i.exp-nexp, nfr)
		if nr == 0 {
			return i.i - 1
		}
		nr = i.r - nr
	default:
		nr = shiftSub(fExponentBits, i.exp-nexp, nfr)
	}
	return int(i.fdr.Div(nr))
}
