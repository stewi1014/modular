package modular32

import (
	"math/bits"
	"unsafe"
)

const (
	nan    = 0x7FE00000
	posinf = 0x7F800000
	neginf = 0xFF800000

	MaxFloat = 0x1p127 * (1 + (1 - 0x1p-23))

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

// Frexp splits a float into it's exponent and fraction component. Sign bit is discarded.
// The 24th implied bit is placed in the fraction if appropriate
func Frexp(f float32) (uint32, uint) {
	fbits := ToBits(f)
	exp := uint((fbits & fExponentMask) >> fFractionBits)
	if exp == 0 {
		return fbits & fFractionMask, 0
	}
	return (fbits & fFractionMask) | (1 << fFractionBits), exp
}

// Ldexp assembles a float from an exponent and fraction component. Sign is ignored.
// Expects the 24th implied bit to be set if appropriate.
func Ldexp(fr uint32, exp uint) float32 {
	if exp == 0 || fr == 0 {
		return FromBits(fr & fFractionMask)
	}
	shift := uint(bits.LeadingZeros32(fr) - fExponentBits)
	if shift >= exp {
		shift = exp - 1
		exp = 0
	} else {
		exp -= uint(shift)
	}
	fr = fr << shift
	return FromBits((uint32(exp) << fFractionBits) | (fr & fFractionMask))
}

func ToBits(n float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&n))
}

func FromBits(bits uint32) float32 {
	return *(*float32)(unsafe.Pointer(&bits))
}

func Abs(n float32) float32 {
	if n <= 0 {
		return -n
	}
	return n
}

func NaN() float32 {
	return FromBits(nan)
}

func IsInf(n float32) bool {
	return n > MaxFloat || n < -MaxFloat
}

func Inf(sign int) float32 {
	if sign < 0 {
		return FromBits(neginf)
	}
	return FromBits(posinf)
}
