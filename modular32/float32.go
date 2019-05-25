package modular32

import (
	"math"
	"math/bits"

	"github.com/chewxy/math32"
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

// shiftSub shifts n up by up-down
func shiftSub(up, down uint, n uint32) uint32 {
	if up > down {
		return n << (up - down)
	}
	return n >> (down - up)
}

// frexp splits a float into it's exponent and fraction component. Sign bit is discarded.
// The 24th implied bit is placed in the fraction if appropriate
func frexp(f float32) (uint32, uint) {
	fbits := math32.Float32bits(f)
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
		return math32.Float32frombits(fr & fFractionMask)
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
