package modular64

import (
	"bytes"
	"encoding/binary"
	"math/bits"

	"github.com/bmkessler/fastdiv"
	"github.com/chewxy/math32"
)

// NewModulus creates a new Modulus
//
// Special cases:
// NewModulus(0) = panic(integer divide by zero)
// NewModulus(±Inf) = Modulus{+Inf}
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

// Congruent returns n mod m
//
// Special cases:
// Modulus{+Inf}.Congruent(±n) = +n
func (m Modulus) Congruent(n float32) float32 {
	if math32.IsNaN(n) || math32.IsInf(n, 0) || math32.IsNaN(m.mod) {
		return math32.NaN()
	}
	if math32.IsInf(m.mod, 0) {
		return math32.Abs(n)
	}

	if n < m.mod && n > -m.mod {
		if n < 0 {
			return n + m.mod
		}
		return n
	}

	nfr, nexp := frexp(n)
	expdiff := nexp - m.exp
	if m.exp == 0 && nexp > 0 {
		expdiff-- //If we're going to shift into denormalised-land, we need to note that denormals have
	}

	//Iterativly apply exponent, trying to take the lagest possible chunk every iteration
	for {
		shift := uint(bits.LeadingZeros32(nfr)) // Find the maximum amount we can shift
		if shift > expdiff {                    // Don't want to shift too far
			shift = expdiff
		}
		nfr = nfr << shift  // iteratively apply exponent
		nfr = m.fd.Mod(nfr) // apply mod
		expdiff -= shift
		if expdiff == 0 { // we've moved as much as we need to
			break
		}
	}

	r := ldexp(nfr, m.exp)

	if n < 0 && r != 0 {
		r = -r + m.mod // correctly handle negatives
	}

	return r
}

// Mod returns the modulus
func (m Modulus) Mod() float32 {
	return m.mod
}

// MarshalBinary implements binary.BinaryMarshaler
func (m Modulus) MarshalBinary() ([]byte, error) {
	buff := new(bytes.Buffer)
	binary.Write(buff, binary.LittleEndian, m.mod)
	return buff.Bytes(), nil
}

// UnmarshalBinary implements binary.BinaryUnmarshaler
func (m *Modulus) UnmarshalBinary(data []byte) error {
	buff := bytes.NewReader(data)
	var f float32
	binary.Read(buff, binary.LittleEndian, &f)
	*m = NewModulus(f)
	return nil
}
