package modular64

import (
	"math/bits"
	"unsafe"
)

const (
	nan    = 0x7FF8000000000001
	posinf = 0x7FF0000000000000
	neginf = 0xFFF0000000000000

	MaxFloat = 0x1p1023 * (1 + (1 - 0x1p-52))

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

// shiftSub shifts n up by up-down
func shiftSub(up, down uint, n uint64) uint64 {
	if up > down {
		return n << (up - down)
	}
	return n >> (down - up)
}

// frexp splits a float into it's exponent and fraction component. Sign bit is discarded.
// The 53rd implied bit is placed in the fraction if appropriate
func frexp(f float64) (uint64, uint) {
	fbits := ToBits(f)
	exp := uint((fbits & fExponentMask) >> fFractionBits)
	if exp == 0 {
		return fbits & fFractionMask, 0
	}
	return (fbits & fFractionMask) | (1 << fFractionBits), exp
}

// ldexp assembles a float from an exponent and fraction component. Sign is ignored.
// Expects the 53rd implied bit to be set if appropriate.
func ldexp(fr uint64, exp uint) float64 {
	if exp == 0 || fr == 0 {
		return FromBits(fr & fFractionMask)
	}
	shift := uint(bits.LeadingZeros64(fr) - fExponentBits)
	if shift >= exp {
		shift = exp - 1
		exp = 0
	} else {
		exp -= uint(shift)
	}
	fr = fr << shift
	return FromBits((uint64(exp) << fFractionBits) | (fr & fFractionMask))
}

func ToBits(n float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&n))
}

func FromBits(bits uint64) float64 {
	return *(*float64)(unsafe.Pointer(&bits))
}

func Abs(n float64) float64 {
	if n <= 0 {
		return -n
	}
	return n
}

func NaN() float64 {
	return FromBits(nan)
}

func IsInf(n float64) bool {
	return n > MaxFloat || n < -MaxFloat
}

func Inf(sign int) float64 {
	if sign < 0 {
		return FromBits(neginf)
	}
	return FromBits(posinf)
}
