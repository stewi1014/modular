package modular64_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stewi1014/modular/modular64"
)

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
			modulus: modular64.FromBits(4144),
			arg:     modular64.FromBits(123445),
			want:    modular64.FromBits(3269),
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
		//		Modulus{NaN}.Congruent(n) = NaN
		{
			name:    "NaN modulus",
			modulus: modular64.NaN(),
			arg:     0,
			want:    modular64.NaN(),
		},
		// 		Modulus{±Inf}.Congruent(n>=0) = n
		{
			name:    "Inf modulus, positive number",
			modulus: modular64.Inf(1),
			arg:     0,
			want:    0,
		},
		//		Modulus{±Inf}.Congruent(n<0) = +Inf
		{
			name:    "Inf modulus, negative number",
			modulus: modular64.Inf(1),
			arg:     -1,
			want:    modular64.Inf(1),
		},
		//		Modulus{m}.Congruent(±Inf) = NaN
		{
			name:    "Inf number",
			modulus: 1,
			arg:     modular64.Inf(1),
			want:    modular64.NaN(),
		},
		//		Modulus{m}.Congruent(NaN) = NaN
		{
			name:    "NaN number",
			modulus: 1,
			arg:     modular64.NaN(),
			want:    modular64.NaN(),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modular64.NewModulus(tt.modulus)
			got := m.Mod(tt.arg)
			if got != tt.want && !(got != got && tt.want != tt.want) {
				t.Errorf("Modulus{%v}.Mod(%v) = %v, want %v", tt.modulus, tt.arg, got, tt.want)
			}
		})
	}
}

func TestModulus_Dist(t *testing.T) {
	type args struct {
		n1 float64
		n2 float64
	}
	tests := []struct {
		name    string
		modulus float64
		args    args
		want    float64
	}{
		{
			name:    "Basic test",
			modulus: 100,
			args: args{
				n1: 10,
				n2: 20,
			},
			want: 10,
		},
		{
			name:    "Forwards over 0",
			modulus: 100,
			args: args{
				n1: 90,
				n2: 20,
			},
			want: 30,
		},
		{
			name:    "Backwards over 0",
			modulus: 100,
			args: args{
				n1: 10,
				n2: 90,
			},
			want: -20,
		},
		{
			name:    "Backwards",
			modulus: 100,
			args: args{
				n1: 40,
				n2: 30,
			},
			want: -10,
		},
		{
			name:    "NaN args",
			modulus: 100,
			args: args{
				n1: modular64.NaN(),
				n2: 30,
			},
			want: modular64.NaN(),
		},
		{
			name:    "NaN args",
			modulus: 100,
			args: args{
				n1: 20,
				n2: modular64.NaN(),
			},
			want: modular64.NaN(),
		},
		{
			name:    "NaN modulus",
			modulus: modular64.NaN(),
			args: args{
				n1: 20,
				n2: 30,
			},
			want: modular64.NaN(),
		},
		{
			name:    "Inf modulus",
			modulus: modular64.Inf(1),
			args: args{
				n1: 20,
				n2: 30,
			},
			want: 10,
		},
		{
			name:    "Inf arg",
			modulus: 100,
			args: args{
				n1: modular64.Inf(1),
				n2: 30,
			},
			want: modular64.NaN(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modular64.NewModulus(tt.modulus)
			got := m.Dist(tt.args.n1, tt.args.n2)
			if got != tt.want && !(got != got && tt.want != tt.want) {
				t.Errorf("Modulus.Dist(%v, %v) = %v, want %v (mod %v)", tt.args.n1, tt.args.n2, got, tt.want, tt.modulus)
			}
		})
	}
}

func TestModulus_GetCongruent(t *testing.T) {
	type args struct {
		n1 float64
		n2 float64
	}
	tests := []struct {
		name    string
		modulus float64
		args    args
		want    float64
	}{
		{
			name:    "Backwards",
			modulus: 100,
			args: args{
				n1: 230,
				n2: 20,
			},
			want: 220,
		},
		{
			name:    "Forward",
			modulus: 100,
			args: args{
				n1: 210,
				n2: 20,
			},
			want: 220,
		},
		{
			name:    "Negative",
			modulus: 100,
			args: args{
				n1: -350,
				n2: 20,
			},
			want: -380,
		},
		{
			name:    "Over 0",
			modulus: 100,
			args: args{
				n1: -310,
				n2: 20,
			},
			want: -280,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modular64.NewModulus(tt.modulus)
			if got := m.GetCongruent(tt.args.n1, tt.args.n2); got != tt.want {
				t.Errorf("Modulus.GetCongruent(%v, %v) = %v, want %v (mod %v)", tt.args.n1, tt.args.n2, got, tt.want, tt.modulus)
			}
		})
	}
}

func TestModulus_Misc(t *testing.T) {
	t.Run("Mod() test", func(t *testing.T) {
		m := modular64.NewModulus(15)
		got := m.Modulus()
		if got != 15 {
			t.Errorf("Modulus.Modulus() = %v, want %v", got, 15)
		}
	})
}

var benchmarkModulo = float64(1e-25)
var benchmarks = []float64{
	0,
	2.5e-25,
	1,
	1e300,
}

func BenchmarkMath_Mod(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Math.Mod(%v)", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				float64Sink = math.Mod(n, benchmarkModulo)
			}
		})
	}
}

func BenchmarkModulus(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Mod(%v)", n), func(b *testing.B) {
			m := modular64.NewModulus(benchmarkModulo)
			for i := 0; i < b.N; i++ {
				float64Sink = m.Mod(n)
			}
		})
	}
}
