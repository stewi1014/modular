package modular32_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/chewxy/math32"
	"github.com/stewi1014/modular/modular32"
)

const randomTestNum = 20000

var (
	float32Sink float32
)

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
			name:    "No change test",
			modulus: 435,
			arg:     434,
			want:    434,
		},
		{
			name:    "Small test",
			modulus: 0.1,
			arg:     0.17,
			want:    0.07,
		},
		{
			name:    "Very small test",
			modulus: math.Float32frombits(4144),
			arg:     math.Float32frombits(123445),
			want:    math.Float32frombits(3269),
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
			name:    "Denormalised edge case",
			modulus: math32.Ldexp(1, -126),
			arg:     math32.Ldexp(1.003, -126),
			want:    math32.Mod(math32.Ldexp(1.003, -126), math32.Ldexp(1, -126)),
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
			m := modular32.NewModulus(tt.modulus)
			got := m.Congruent(tt.arg)
			if got != tt.want && !(math32.IsNaN(got) && math32.IsNaN(tt.want)) {
				t.Errorf("Modulus{%v}.Congruent(%v) = %v, want %v", tt.modulus, tt.arg, got, tt.want)
			}
		})
	}
}

func TestModulus_Misc(t *testing.T) {
	t.Run("Mod() test", func(t *testing.T) {
		m := modular32.NewModulus(15)
		got := m.Mod()
		if got != 15 {
			t.Errorf("Modulus.Mod() = %v, want %v", got, 15)
		}
	})
}

var benchmarks = []float32{
	0,
	1,
	20,
	1e20,
}

func BenchmarkMath_Mod(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Math.Mod(%v)", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				float32Sink = math32.Mod(n, 1e-20)
			}
		})
	}
}

func BenchmarkModulus(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Congruent(%v)", n), func(b *testing.B) {
			m := modular32.NewModulus(1e-20)
			for i := 0; i < b.N; i++ {
				float32Sink = m.Congruent(n)
			}
		})
	}
}
