package modular32

import (
	"errors"

	"github.com/bmkessler/fastdiv"
	"github.com/chewxy/math32"
)

// Error types
var (
	ErrBadModulo = errors.New("bad modulus")
	ErrBadIndex  = errors.New("bad index")
)

// NewIndexer creates a new Indexer
//
// Index must not be larger than 2**16
// Modulus must be a normalised float
//
// Special cases:
// NewIndexer(m, 0) = panic(integer divide by zero)
// NewIndexer(m, i > 2**16) = ErrIndexTooBig
// NewIndexer(0, i) = ErrBadModulo
// NewIndexer(±Inf, i) = ErrBadModulo
// NewIndexer(NaN, i) = ErrBadModulo
// NewIndexer(m, i) = ErrBadModulo for |m| < 2**-126
func NewIndexer(modulus float32, index int) (Indexer, error) {
	mod := NewModulus(modulus)
	return mod.NewIndexer(index)
}

// NewIndexer creates a new indexer from the Modulus
func (m Modulus) NewIndexer(index int) (Indexer, error) {
	if math32.IsInf(m.mod, 0) || math32.IsNaN(m.mod) || m.exp == 0 {
		return Indexer{}, ErrBadModulo
	}
	if index > (1<<16) || index < 1 {
		return Indexer{}, ErrBadIndex
	}

	modfr, _ := frexp(m.mod)
	r := modfr << fExponentBits //r - range; is shifted fExponentBits to get a little more
	rDivisor := r / uint32(index)
	return Indexer{
		Modulus: m,
		fdr:     fastdiv.NewUint32(rDivisor),
		r:       r,
		i:       index,
	}, nil
}

// Indexer provides a fast method for mapping a floating point modulus to a range of integers
type Indexer struct {
	Modulus
	fdr fastdiv.Uint32
	r   uint32
	i   int
}

// Index indexes n to the integer range 0 <= num < index
// If n is NaN or ±Inf, it returns the index.
//
// Special cases:
// Index(NaN) = index
// Index(±Inf) = index
func (i Indexer) Index(n float32) int {
	if math32.IsNaN(n) || math32.IsInf(n, 0) {
		return i.i
	}

	nfr, nexp := frexp(n)
	var nr uint32
	switch {
	case n > i.mod:
		expdiff := nexp - i.exp
		nr = i.modFrExp(nfr, expdiff) << fExponentBits
	case n < -i.mod:
		expdiff := nexp - i.exp
		nr = i.modFrExp(nfr, expdiff) << fExponentBits
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
