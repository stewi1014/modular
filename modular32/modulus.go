package modular32

import (
	"github.com/bmkessler/fastdiv"
	math "github.com/chewxy/math32"
)

// NewModulus creates a new Modulus.
// An Infinite modulus has no effect other than to waste CPU time
//
// Special cases:
// 		NewModulus(0) = panic(integer divide by zero)
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
		mod:    math.Abs(modulus),
		exp:    modexp,
	}

	return mod
}

// Modulus defines a modulus.
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

// Dist returns the distance and direction of n1 to n2.
func (m Modulus) Dist(n1, n2 float32) float32 {
	d := m.Congruent(n2 - n1)
	if d > m.mod/2 {
		return d - m.mod
	}
	return d
}

// GetCongruent returns the closest number to n1 that is congruent to n2.
func (m Modulus) GetCongruent(n1, n2 float32) float32 {
	return n1 - m.Dist(n2, n1)
}

// Congruent returns n mod m.
//
// Special cases:
//		Modulus{NaN}.Congruent(n) = NaN
// 		Modulus{±Inf}.Congruent(n>=0) = n
//		Modulus{±Inf}.Congruent(n<0) = +Inf
//		Modulus{m}.Congruent(±Inf) = NaN
//		Modulus{m}.Congruent(NaN) = NaN
func (m Modulus) Congruent(n float32) float32 {
	if m.mod == 0 || m.mod != m.mod { // 0 or NaN modulus
		return math.NaN()
	}

	nfr, nexp := frexp(n)

	if n < m.mod && n > -m.mod {
		if n < 0 {
			return n + m.mod
		}
		return n
	}

	if nexp == fMaxExp {
		return math.NaN()
	}

	expdiff := nexp - m.exp
	if m.exp == 0 && expdiff > 0 {
		expdiff-- //We're in denormalised land, skip an exponent.
	}

	nfr = m.modExp(nfr, expdiff)

	r := ldexp(nfr, m.exp)

	if n < 0 && r != 0 {
		r = m.mod - r // correctly handle negatives
	}

	return r
}

// modExp returns a * 2**exp (mod m)
func (m Modulus) modExp(a uint32, exp uint) uint32 {
	return uint32(m.fd.Mod(uint64(a) * m.powers[exp]))
}
