package modular64

import (
	"math"
	"math/bits"

	"github.com/bmkessler/fastdiv"
)

// NewModulus creates a new Modulus.
//
// An Infinite modulus has no effect other than to waste CPU time.
//
// Special cases:
//		NewModulus(0) = panic(integer divide by zero)
func NewModulus(modulus float64) Modulus {
	modfr, modexp := frexp(modulus)
	mod := Modulus{
		fd:  fastdiv.NewUint64(modfr),
		mod: math.Abs(modulus),
		exp: modexp,
	}

	return mod
}

// Modulus defines a modulus.
// It offers greater performance than traditional floating point modulo calculations by pre-computing the inverse of the modulus's fractional component.
// This obviously adds overhead to the creation of a new Modulus, but quickly breaks even after a few calls to Congruent.
type Modulus struct {
	fd  fastdiv.Uint64
	mod float64
	exp uint
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
func (m Modulus) Congruent(n float64) float64 {
	if math.IsNaN(n) || math.IsInf(n, 0) || math.IsNaN(m.mod) {
		return math.NaN()
	}
	if math.IsInf(m.mod, 0) {
		return math.Abs(n)
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

	//Iterativly apply exponent to the fraction, trying to take the lagest possible chunk every iteration
	for {
		shift := uint(bits.LeadingZeros64(nfr)) // Find the maximum chunk we can take
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
