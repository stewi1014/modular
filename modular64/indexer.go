package modular64

import (
	"errors"
	"math"

	"github.com/bmkessler/fastdiv"
)

// Error types
var (
	ErrIndexTooBig = errors.New("index is too big")
	ErrBadModulo   = errors.New("bad modulus")
)

// NewIndexer creates a new Indexer
//
// Large indexs that approach or exceed the resolution of floats require more computation to stay accurate.
// The same is also true for denormalised moduli. I don't see these as common use cases.
// As such, in keeping this library fast, I've chosen 2**32 as an upper limit on the size of the index,
// and normalised floats for moduli; >=2**-1022, not NaN and not ±Inf as limits for moduli.
// This should be far more than needed. The primary use case of this is to index lists.
// If you need to cast floats to an integer range larger than this, you should do something else.
//
// Special cases:
// NewIndexer(m, 0) = panic(integer divide by zero)
// NewIndexer(m, i > 2**32) = ErrIndexTooBig
// NewIndexer(0, i) = ErrBadModulo
// NewIndexer(±Inf, i) = ErrBadModulo
// NewIndexer(NaN, i) = ErrBadModulo
// NewIndexer(m, i) = ErrBadModulo for |m| < 2**-1022
func NewIndexer(modulus float64, index uint) (Indexer, error) {
	mod := NewModulus(modulus)
	return mod.NewIndexer(index)
}

// NewIndexer creates a new indexer from the Modulus
func (m Modulus) NewIndexer(index uint) (Indexer, error) {
	if math.IsInf(m.mod, 0) || math.IsNaN(m.mod) || m.exp == 0 {
		return Indexer{}, ErrBadModulo
	}
	if index > (1 << 32) {
		return Indexer{}, ErrIndexTooBig
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

// Indexer provides a fast method for mapping a floating point modulus to a range of integers
type Indexer struct {
	Modulus
	fdr fastdiv.Uint64
	r   uint64
	i   uint
}

// Index indexes n to the integer range 0 <= num < index
//
// Special cases:
// Index(NaN) = 0
// Index(±Inf) = 0
func (i Indexer) Index(n float64) uint {
	if math.IsNaN(n) || math.IsInf(n, 0) {
		return 0
	}

	nfr, nexp := frexp(n)
	var nr uint64
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
	return uint(i.fdr.Div(nr))
}
