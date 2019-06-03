package modular32_test

import (
	"fmt"
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
			modulus: math32.Float32frombits(4144),
			arg:     math32.Float32frombits(123445),
			want:    math32.Float32frombits(3269),
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
			name:    "Infinte number and modulo test",
			modulus: math32.Inf(1),
			arg:     math32.Inf(1),
			want:    math32.Inf(1),
		},
		{
			name:    "Denormalised edge case",
			modulus: math32.Ldexp(1, -126),
			arg:     math32.Ldexp(1.003, -126),
			want:    math32.Mod(math32.Ldexp(1.003, -126), math32.Ldexp(1, -126)),
		},
		{
			name:    "Denormalised edge case2",
			modulus: math32.Ldexp(1, -127),
			arg:     math32.Ldexp(1.003, -126),
			want:    math32.Mod(math32.Ldexp(1.003, -126), math32.Ldexp(1, -127)),
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

func TestModulus_Dist(t *testing.T) {
	type args struct {
		n1 float32
		n2 float32
	}
	tests := []struct {
		name    string
		modulus float32
		args    args
		want    float32
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
				n1: math32.NaN(),
				n2: 30,
			},
			want: math32.NaN(),
		},
		{
			name:    "NaN args",
			modulus: 100,
			args: args{
				n1: 20,
				n2: math32.NaN(),
			},
			want: math32.NaN(),
		},
		{
			name:    "NaN modulus",
			modulus: math32.NaN(),
			args: args{
				n1: 20,
				n2: 30,
			},
			want: math32.NaN(),
		},
		{
			name:    "Inf modulus",
			modulus: math32.Inf(1),
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
				n1: math32.Inf(1),
				n2: 30,
			},
			want: math32.NaN(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modular32.NewModulus(tt.modulus)
			got := m.Dist(tt.args.n1, tt.args.n2)
			if got != tt.want && !(math32.IsNaN(got) && math32.IsNaN(tt.want)) {
				t.Errorf("Modulus.Dist(%v, %v) = %v, want %v (mod %v)", tt.args.n1, tt.args.n2, got, tt.want, tt.modulus)
			}
		})
	}
}

func TestModulus_GetCongruent(t *testing.T) {
	type args struct {
		n1 float32
		n2 float32
	}
	tests := []struct {
		name    string
		modulus float32
		args    args
		want    float32
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
			m := modular32.NewModulus(tt.modulus)
			if got := m.GetCongruent(tt.args.n1, tt.args.n2); got != tt.want {
				t.Errorf("Modulus.GetCongruent(%v, %v) = %v, want %v (mod %v)", tt.args.n1, tt.args.n2, got, tt.want, tt.modulus)
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
