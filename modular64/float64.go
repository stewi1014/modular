package modular64

import (
	"math"
	"math/bits"
)

const (
	fExponentBits = 11
	fFractionBits = 52
	fTotalBits    = 64

	fMaxExp       = (1 << fExponentBits) - 1
	fMinExp       = 0
	fBias         = (1 << (fExponentBits - 1)) - 1
	fSignMask     = 1 << (fTotalBits - 1)
	fExponentMask = (1<<fExponentBits - 1) << fFractionBits
	fFractionMask = (1 << fFractionBits) - 1
)

func frexp(f float64) (uint64, uint) {
	fbits := math.Float64bits(f)
	exp := uint((fbits & fExponentMask) >> 52)
	if exp == 0 {
		return fbits & fFractionMask, 0
	}
	return (fbits & fFractionMask) | (1 << 52), exp
}

func ldexp(fr uint64, exp uint) float64 {
	if exp == 0 || fr == 0 {
		return math.Float64frombits(fr & fFractionMask)
	}
	shift := uint(bits.LeadingZeros64(fr) - 11)
	if shift >= exp {
		shift = exp - 1
		exp = 0
	} else {
		exp -= uint(shift)
	}
	fr = fr << shift
	return math.Float64frombits((uint64(exp) << 52) | (fr & fFractionMask))
}

// fexp returns the float's exponent
func fexp(f float64) uint {
	return uint((math.Float64bits(f) & fExponentMask) >> 52)
}

// n = 2**exp
// modulus cannot be 0
func pow2mod(exp uint, modulus uint64) uint64 {
	if exp == 0 {
		return 1
	}
	if modulus&(1<<63) > 0 {
		// Here for debugging, this function is never called with such a large modulus
		panic("modulus too large; calculations overflow uint64")
	}
	r := uint64(2)
	for i := uint(1); i < exp; i++ {
		r *= 2
		r = r % modulus
	}
	return r
}
