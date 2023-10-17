package modular64

import (
	"math/bits"

	"github.com/bmkessler/fastdiv"
)

// NewModulus creates a new Modulus.
//
// Special cases:
//
//	NewModulus(0) = panic(integer divide by zero)
func NewModulus(modulus float64) Modulus {
	modfr, modexp := frexp(modulus)
	fd := fastdiv.NewUint64(modfr)

	const minPowerLen = 65
	powerlen := fMaxExp - modexp
	if powerlen < minPowerLen {
		powerlen = minPowerLen
	}

	powers := make([]uint64, powerlen)
	r := uint64(1)
	if len(powers) > 0 {
		powers[0] = 1
	}
	for i := 1; i < len(powers); i++ {
		r = r << 1
		r = fd.Mod(r)
		powers[i] = uint64(r)
	}

	mod := Modulus{
		fd:     fastdiv.NewUint64(modfr),
		powers: powers,
		mod:    Abs(modulus),
		fr:     modfr,
		exp:    modexp,
	}

	return mod
}

// Modulus defines a modulus.
// It offers greater performance than traditional floating point modulo calculations by pre-computing the inverse of the modulus's fractional component.
// This obviously adds overhead to the creation of a new Modulus, but quickly breaks even after a few calls to Congruent.
type Modulus struct {
	fd     fastdiv.Uint64
	powers []uint64
	mod    float64
	fr     uint64
	exp    uint
}

// Modulus returns the modulus.
func (m Modulus) Modulus() float64 {
	return m.mod
}

// Dist returns the distance and direction of n1 to n2.
func (m Modulus) Dist(n1, n2 float64) float64 {
	d := m.Mod(n2 - n1)
	if d > m.mod/2 {
		return d - m.mod
	}
	return d
}

// GetCongruent returns the closest number to n1 that is congruent to n2.
func (m Modulus) GetCongruent(n1, n2 float64) float64 {
	return n1 - m.Dist(n2, n1)
}

// Mod returns n mod m.
//
// Special cases:
//
//	Modulus{NaN}.Mod(n) = NaN
//	Modulus{±Inf}.Mod(n>=0) = n
//	Modulus{±Inf}.Mod(n<0) = +Inf
//	Modulus{m}.Mod(±Inf) = NaN
//	Modulus{m}.Mod(NaN) = NaN
func (m Modulus) Mod(n float64) float64 {
	if m.mod == 0 || m.mod != m.mod { // 0 or NaN modulus
		return NaN()
	}

	nfr, nexp := frexp(n)

	if n < m.mod && n > -m.mod {
		if n < 0 {
			return n + m.mod
		}
		return n
	}

	if nexp == fMaxExp {
		return NaN()
	}

	expdiff := nexp - m.exp
	if m.exp == 0 && nexp != 0 {
		expdiff-- //We're in denormalised land, skip an exponent.
	}

	rfr := m.modExp(nfr, expdiff)

	r := ldexp(rfr, m.exp)

	if n < 0 && r != 0 {
		r = m.mod - r
	}

	return r
}

// modExp returns n * 2**exp (mod m)
func (m Modulus) modExp(n uint64, exp uint) uint64 {
	switch { // Switch fastest computation method
	case exp <= uint(bits.LeadingZeros64(n)):
		return m.fd.Mod(n << exp)

	default:
		hi, lo := bits.Mul64(n, m.powers[exp])
		_, q := bits.Div64(hi, lo, m.fr)
		return q
	}
}
