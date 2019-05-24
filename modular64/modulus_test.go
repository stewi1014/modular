package modular64

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
)

const randomTestNum = 20000

var (
	varNumber float64 = 234
	varMod    float64 = 16
	varSink   float64

	varUintNumber uint = 13267489
	varUintMod    uint = 293
	varUintSink   uint
)

func makeDenormFloat(fr uint64) float64 {
	return math.Float64frombits(fr)
}

func TestModulus_Congruent(t *testing.T) {
	tests := []struct {
		name    string
		modulus float64
		arg     float64
		want    float64
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
			want:    math.Mod(456897613245865, 10),
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
			arg:     math.NaN(),
			want:    math.NaN(),
		},
		{
			name:    "Infinte modulo test",
			modulus: math.Inf(-1),
			arg:     0.01,
			want:    0.01,
		},
		{
			name:    "Infinte number test",
			modulus: 2,
			arg:     math.Inf(1),
			want:    math.NaN(),
		},
		{
			name:    "Denormalised edge case",
			modulus: math.Ldexp(1, -1022),
			arg:     math.Ldexp(1.003, -1022),
			want:    math.Mod(math.Ldexp(1.003, -1022), math.Ldexp(1, -1022)),
		},
		{
			name:    "Generated case 1",
			modulus: -1.16358406231669e-309,
			arg:     -2.3400420480579135e+145,
			want:    math.Mod(-2.3400420480579135e+145, 1.16358406231669e-309) + 1.16358406231669e-309,
		},
		{
			name:    "Generated case 2",
			modulus: 1.0863201545657832e-307,
			arg:     1.3303463489150531e+13,
			want:    math.Mod(1.3303463489150531e+13, 1.0863201545657832e-307),
		},
		{
			name:    "Generated case 3",
			modulus: 2.039381663448266e-229,
			arg:     -1.370217367318819e-267,
			want:    math.Mod(-1.370217367318819e-267, 2.039381663448266e-229) + 2.039381663448266e-229,
		},
		{
			name:    "NaN modulo",
			modulus: math.NaN(),
			arg:     0.01,
			want:    math.NaN(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModulus(tt.modulus)
			got := m.Congruent(tt.arg)
			if got != tt.want && !(math.IsNaN(got) && math.IsNaN(tt.want)) {
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

func randomFloat() float64 {
	b := rand.Uint64()
	f := ldexp(b&fFractionMask, uint(b&fExponentMask)>>52)
	if b&fSignMask > 0 {
		f = -f
	}
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return randomFloat()
	}
	return f
}

func TestModulus_Congruent_random(t *testing.T) {
	for i := 0; i < randomTestNum; i++ {
		modulus := randomFloat()
		arg := randomFloat()
		want := math.Mod(arg, modulus)
		if want < 0 {
			want = want + math.Abs(modulus)
		}
		t.Run("Random Test", func(t *testing.T) {
			m := NewModulus(modulus)
			got := m.Congruent(arg)
			if got != want && !(math.IsNaN(got) && math.IsNaN(want)) {
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
		math.Mod(varNumber, varMod)
	}
}

func Benchmark_Modulus_Congruent(b *testing.B) {
	mod := NewModulus(varMod)
	for i := 0; i < b.N; i++ {
		mod.Congruent(varNumber)
	}
}
