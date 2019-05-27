package modular64

import (
	"math"
	"math/bits"
)

// NewModulus creates a new Modulus.
//
// An Infinite modulus has no effect other than to waste CPU time.
//
// Special cases:
//		NewModulus(0) = panic(integer divide by zero)
func NewModulus(modulus float64) Modulus {
	modfr, modexp := frexp(modulus)

	powers := make([]uint64, fMaxExp-modexp)
	r := uint64(1)
	if len(powers) > 0 {
		powers[0] = 1
	}
	for i := 1; i < len(powers); i++ {
		r = r << 1
		r = r % modfr
		powers[i] = uint64(r)
	}

	mod := Modulus{
		modfr:  modfr,
		powers: powers,
		mod:    math.Abs(modulus),
		exp:    modexp,
	}

	return mod
}

// Modulus defines a modulus.
// It offers greater performance than traditional floating point modulo calculations by pre-computing the inverse of the modulus's fractional component.
// This obviously adds overhead to the creation of a new Modulus, but quickly breaks even after a few calls to Congruent.
type Modulus struct {
	modfr  uint64
	powers []uint64
	mod    float64
	exp    uint
}

// Mod returns the modulus.
func (m Modulus) Mod() float64 {
	return m.mod
}

// Congruent returns n mod m.
//
// Special cases:
//		Modulus{Inf}.Congruent(±n) = +n
//		Modulus{NaN}.Congruent(±n) = NaN
//		Modulus.Congruent(NaN) = NaN
//		Modulus.Congruent(±Inf) = NaN
func (m Modulus) Congruent(n float64) float64 {
	if math.IsInf(m.mod, 0) {
		return math.Abs(n)
	}
	if math.IsNaN(n) || math.IsInf(n, 0) || math.IsNaN(m.mod) {
		return math.NaN()
	}

	if n < m.mod && n > -m.mod {
		if n < 0 {
			r := n + m.mod
			return r
		}
		return n
	}

	nfr, nexp := frexp(n)
	expdiff := nexp - m.exp
	nfr = m.modFrExp(nfr, expdiff)
	r := ldexp(nfr, m.exp)

	if n < 0 && r != 0 {
		r = m.mod - r // correctly handle negatives
	}

	return r
}

// after doing other checks and optimisations, this is what really does the modulo calulation.
func (m Modulus) modFrExp(nfr uint64, exp uint) uint64 {
	if m.exp == 0 && exp > 0 {
		exp-- //We're in denormalised land, skip an exponent.
	}

	hi, lo := bits.Mul64(nfr, m.powers[exp])

	_, r := bits.Div64(hi, lo, m.modfr)
	return r
}
