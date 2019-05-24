package modular32

import (
	"math"
	"math/bits"
)

// Becuase this makes porting to float32 so much easier
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

// getFractionAt returns the fraction at the given exponent. Truncates bits too high or too low.
func getFractionAt(f float32, exp uint) uint32 {
	ffr, fexp := frexp(f)
	switch {
	case fexp == exp:
		return ffr // That was easy

	case fexp < exp:
		shift := exp - fexp
		if fexp == 0 {
			shift-- // We're in denormalised land
		}
		return ffr >> shift

	case fexp > exp:
		shift := fexp - exp
		if exp == 0 {
			shift-- // Another denormal
		}
		return ffr << shift

	}
	panic("integers have gone crazy - How can a==b, a<b and a>b all be false for an integer?")
}

// frexp splits a float into it's exponent and fraction component. Sign bit is discarded.
// The 24th implied bit is placed in the fraction if appropriate
func frexp(f float32) (uint32, uint) {
	fbits := math.Float32bits(f)
	exp := uint((fbits & fExponentMask) >> fFractionBits)
	if exp == 0 {
		return fbits & fFractionMask, 0
	}
	return (fbits & fFractionMask) | (1 << fFractionBits), exp
}

// ldexp assembles a float from an exponent and fraction component. Sign is ignored.
// Expects the 24th implied bit to be set if appropriate.
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
