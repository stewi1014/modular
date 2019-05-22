package modular64

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"

	"github.com/stewi1014/things"
)

const randomTestNum = 2000

var (
	varNumber float64 = 13
	varMod    float64 = 16
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
			name:    "Failed case from random (fixed)",
			modulus: -1.16358406231669e-309,
			arg:     -2.3400420480579135e+145,
			want:    math.Mod(-2.3400420480579135e+145, 1.16358406231669e-309) + 1.16358406231669e-309,
		},
		{
			name:    "Failed case from random2",
			modulus: 1.0863201545657832e-307,
			arg:     1.3303463489150531e+13,
			want:    math.Mod(1.3303463489150531e+13, 1.0863201545657832e-307),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModulus(tt.modulus)
			got := m.Congruent(tt.arg)
			if got != tt.want && !(math.IsNaN(got) && math.IsNaN(tt.want)) {
				fmt.Println("num", things.FormatFloat64(tt.arg))
				fmt.Println("mod", things.FormatFloat64(tt.modulus))
				fmt.Println("got", things.FormatFloat64(got))
				fmt.Println("wan", things.FormatFloat64(tt.want))
				fmt.Println("")
				t.Errorf("Modulus.Congruent(%v) = %v, want %v", tt.arg, got, tt.want)
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

func TestModulus_CongruentRandom(t *testing.T) {
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
				fmt.Println("num", arg)
				fmt.Println("mod", modulus)
				fmt.Println("got", got)
				fmt.Println("want", want)
				fmt.Println("")
				t.Errorf("Modulus.Congruent(%v) = %v, want %v", arg, got, want)
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
