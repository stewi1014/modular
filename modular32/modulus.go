package modular32

import (
	"github.com/bmkessler/fastdiv"
	"github.com/chewxy/math32"
)

// NewModulus creates a new Modulus
//
// Special cases:
// NewModulus(0) = panic(integer divide by zero)
//
// An Infinite modulus has no effect other than to waste CPU time
func NewModulus(modulus float32) Modulus {
	modfr, modexp := frexp(modulus)

	powers := make([]uint64, fMaxExp-modexp)
	r := uint32(1)
	if len(powers) > 0 {
		powers[0] = 1
	}
	for i := 1; i < len(powers); i++ {
		r = r << 1
		r = r % modfr
		powers[i] = uint64(r)
	}

	mod := Modulus{
		fd:     fastdiv.NewUint64(uint64(modfr)),
		powers: powers,
		mod:    math32.Abs(modulus),
		exp:    modexp,
	}

	return mod
}

// Modulus defines a modulus
// It offers greater performance than traditional floating point modulo calculations by pre-computing the inverse of the modulus's fractional component,
// and pre-computing a lookup table for different exponents in the given modulus, allowing direct computation of n mod m - no iteration or recursion is used.
// This obviously adds overhead to the creation of a new Modulus, but quickly breaks even after a few calls to Congruent.
type Modulus struct {
	fd     fastdiv.Uint64
	powers []uint64
	mod    float32
	exp    uint
}

// Mod returns the modulus
func (m Modulus) Mod() float32 {
	return m.mod
}

// Congruent returns n mod m
//
// Special cases:
// Modulus{Inf}.Congruent(±n) = +n
// Modulus{NaN}.Congruent(±n) = NaN
func (m Modulus) Congruent(n float32) float32 {
	if math32.IsNaN(n) || math32.IsInf(n, 0) || math32.IsNaN(m.mod) {
		return math32.NaN()
	}
	if math32.IsInf(m.mod, 0) {
		return math32.Abs(n)
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
func (m Modulus) modFrExp(nfr uint32, exp uint) uint32 {
	if m.exp == 0 && exp > 0 {
		exp-- //We're in denormalised land, skip an exponent.
	}

	return uint32(m.fd.Mod(uint64(nfr) * m.powers[exp])) //Hooray for direct computation
}
