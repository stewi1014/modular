package modular64_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stewi1014/modular/modular64"
)

const randomTestNum = 20000

var (
	float64Sink float64
)

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
			modulus: math.Float64frombits(4144),
			arg:     math.Float64frombits(123445),
			want:    math.Float64frombits(3269),
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
			name:    "Infinte number and modulo test",
			modulus: math.Inf(1),
			arg:     math.Inf(1),
			want:    math.Inf(1),
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
			m := modular64.NewModulus(tt.modulus)
			got := m.Congruent(tt.arg)
			if got != tt.want && !(math.IsNaN(got) && math.IsNaN(tt.want)) {
				t.Errorf("Modulus{%v}.Congruent(%v) = %v, want %v", tt.modulus, tt.arg, got, tt.want)
			}
		})
	}
}

func TestModulus_Misc(t *testing.T) {
	t.Run("Mod() test", func(t *testing.T) {
		m := modular64.NewModulus(15)
		got := m.Mod()
		if got != 15 {
			t.Errorf("Modulus.Mod() = %v, want %v", got, 15)
		}
	})
}

var benchmarks = []float64{
	0,
	1,
	20,
	1e20,
	1e150,
	1e300,
}

func BenchmarkMath_Mod(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Math.Mod(%v)", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				float64Sink = math.Mod(n, 1)
			}
		})
	}
}

func BenchmarkModulus(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Congruent(%v)", n), func(b *testing.B) {
			m := modular64.NewModulus(1)
			for i := 0; i < b.N; i++ {
				float64Sink = m.Congruent(n)
			}
		})
	}
}
