package modular64

import (
	"math"
	"math/bits"
)

const (
	fExponentBits = 8
	fFractionBits = 23
	fTotalBits    = 32

	fMaxExp       = (1 << fExponentBits) - 1
	fMinExp       = 0
	fBias         = (1 << (fExponentBits - 1)) - 1
	fSignMask     = 1 << (fTotalBits - 1)
	fExponentMask = (1<<fExponentBits - 1) << fFractionBits
	fFractionMask = (1 << fFractionBits) - 1
)

func frexp(f float32) (uint32, uint) {
	fbits := math.Float32bits(f)
	exp := uint((fbits & fExponentMask) >> fFractionBits)
	if exp == 0 {
		return fbits & fFractionMask, 0
	}
	return (fbits & fFractionMask) | (1 << fFractionBits), exp
}

func ldexp(fr uint32, exp uint) float32 {
	if exp == 0 || fr == 0 {
		return math.Float32frombits(fr & fFractionMask)
	}
	shift := uint(bits.LeadingZeros32(fr) - fExponentBits)
	if shift >= exp {
		shift = exp - 1
		exp = 0
	} else {
		exp -= uint(shift)
	}
	fr = fr << shift
	return math.Float32frombits((uint32(exp) << fFractionBits) | (fr & fFractionMask))
}

// fexp returns the float's exponent
func fexp(f float32) uint {
	return uint((math.Float32bits(f) & fExponentMask) >> fFractionBits)
}

// n = 2**exp
// modulus cannot be 0
// this is slow
func pow2mod(exp uint, modulus uint32) uint32 {
	if exp == 0 {
		return 1
	}
	if modulus&(1<<31) > 0 {
		// Here for debugging, this function is never called with such a large modulus
		panic("modulus too large; calculations overflow uint64")
	}
	r := uint32(2)
	for i := uint(1); i < exp; i++ {
		r *= 2
		r = r % modulus
	}
	return r
}
