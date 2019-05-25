package modular32

import (
	"math/bits"

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
	mod := Modulus{
		fd:  fastdiv.NewUint32(modfr),
		mod: math32.Abs(modulus),
		exp: modexp,
	}

	return mod
}

// Modulus defines a modulus
// It offers greater performance than traditional floating point modulo calculations by pre-computing the inverse of the modulus's fractional component.
// This obviously adds overhead to the creation of a new Modulus, but quickly breaks even after a few calls to Congruent.
type Modulus struct {
	fd  fastdiv.Uint32
	mod float32
	exp uint
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

	//Iterativly apply exponent to the fraction, trying to take the lagest possible chunk every iteration
	for {
		shift := uint(bits.LeadingZeros32(nfr)) // Find the maximum chunk we can take
		if shift > exp {                        // Don't want to shift too far
			shift = exp
		}
		nfr = nfr << shift  // apply a chunk of exponent
		nfr = m.fd.Mod(nfr) // apply mod
		exp -= shift
		if exp == 0 { // we've applied all the exponent
			break
		}
	}

	return nfr
}
