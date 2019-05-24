package modular32

import (
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/chewxy/math32"
)

const randomTestNum = 230000

var (
	varNumber float32 = 234
	varMod    float32 = 16
	varSink   float32

	varUintNumber uint = 13267489
	varUintMod    uint = 293
	varUintSink   uint
)

func makeDenormFloat(fr uint32) float32 {
	return math.Float32frombits(fr)
}

func TestModulus_Congruent(t *testing.T) {
	tests := []struct {
		name    string
		modulus float32
		arg     float32
		want    float32
	}{
		{
			name:    "Basic test",
			modulus: 13,
			arg:     58,
			want:    6,
		},
		{
			name:    "Small test",
			modulus: 0.1,
			arg:     0.17,
			want:    0.07,
		},
		{
			name:    "Very small test",
			modulus: makeDenormFloat(4144),
			arg:     makeDenormFloat(123445),
			want:    makeDenormFloat(3269),
		},
		{
			name:    "very big test with small modulus",
			modulus: 10,
			arg:     456897613245865,
			want:    math32.Mod(456897613245865, 10),
		},
		{
			name:    "Negative number",
			modulus: 5,
			arg:     -34,
			want:    1,
		},
		{
			name:    "Negative modulo and number",
			modulus: -5,
			arg:     -3,
			want:    2,
		},
		{
			name:    "NaN number",
			modulus: 92786534,
			arg:     math32.NaN(),
			want:    math32.NaN(),
		},
		{
			name:    "Infinte modulo test",
			modulus: math32.Inf(-1),
			arg:     0.01,
			want:    0.01,
		},
		{
			name:    "Infinte number test",
			modulus: 2,
			arg:     math32.Inf(1),
			want:    math32.NaN(),
		},
		{
			name:    "NaN modulo",
			modulus: math32.NaN(),
			arg:     0.01,
			want:    math32.NaN(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModulus(tt.modulus)
			got := m.Congruent(tt.arg)
			if got != tt.want && !(math32.IsNaN(got) && math32.IsNaN(tt.want)) {
				t.Errorf("Modulus{%v}.Congruent(%v) = %v, want %v", tt.modulus, tt.arg, got, tt.want)
			}
		})
	}
}

func TestModulus_misc(t *testing.T) {
	t.Run("Mod() test", func(t *testing.T) {
		m := NewModulus(varMod)
		got := m.Mod()
		if got != varMod {
			t.Errorf("Modulus.Mod() = %v, want %v", got, varMod)
		}
	})
}

func randomFloat() float32 {
	b := rand.Uint32()
	f := ldexp(b&fFractionMask, uint(b&fExponentMask)>>52)
	if b&fSignMask > 0 {
		f = -f
	}
	if math32.IsNaN(f) || math32.IsInf(f, 0) {
		return randomFloat()
	}
	return f
}

func TestModulus_Congruent_random(t *testing.T) {
	for i := 0; i < randomTestNum; i++ {
		modulus := randomFloat()
		arg := randomFloat()
		want := math32.Mod(arg, modulus)
		if want < 0 {
			want = want + math32.Abs(modulus)
		}
		t.Run("Random Test", func(t *testing.T) {
			m := NewModulus(modulus)
			got := m.Congruent(arg)
			if got != want && !(math32.IsNaN(got) && math32.IsNaN(want)) {
				t.Errorf("Modulus{%v}.Congruent(%v) = %v, want %v", modulus, arg, got, want)
			}
		})
	}
}

func TestModulus_MarshalBinary(t *testing.T) {
	tests := []struct {
		name    string
		modulus Modulus
	}{
		{
			name:    "Basic test",
			modulus: NewModulus(1),
		},
		{
			name:    "Basic test",
			modulus: NewModulus(13425),
		},
		{
			name:    "Basic test",
			modulus: NewModulus(1221313),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, _ := tt.modulus.MarshalBinary() // Never returns an error
			newmod := Modulus{}
			newmod.UnmarshalBinary(data)
			if !reflect.DeepEqual(tt.modulus, newmod) {
				t.Errorf("Modulus.Unmarshal() = %v, want %v", newmod, tt.modulus)
			}
		})
	}
}

func Benchmark_math_Mod(b *testing.B) {
	for i := 0; i < b.N; i++ {
		math32.Mod(varNumber, varMod)
	}
}

func Benchmark_Modulus_Congruent(b *testing.B) {
	mod := NewModulus(varMod)
	for i := 0; i < b.N; i++ {
		mod.Congruent(varNumber)
	}
}
