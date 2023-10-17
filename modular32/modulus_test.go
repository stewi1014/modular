package modular32_test

import (
	"fmt"
	"testing"

	"github.com/stewi1014/modular/modular32"
)

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
			modulus: modular32.FromBits(4144),
			arg:     modular32.FromBits(123445),
			want:    modular32.FromBits(3269),
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
		//		Modulus{NaN}.Mod(n) = NaN
		{
			name:    "NaN modulus",
			modulus: modular32.NaN(),
			arg:     0,
			want:    modular32.NaN(),
		},
		// 		Modulus{±Inf}.Mod(n>=0) = n
		{
			name:    "Inf modulus, positive number",
			modulus: modular32.Inf(1),
			arg:     0,
			want:    0,
		},
		//		Modulus{±Inf}.Mod(n<0) = +Inf
		{
			name:    "Inf modulus, negative number",
			modulus: modular32.Inf(1),
			arg:     -1,
			want:    modular32.Inf(1),
		},
		//		Modulus{m}.Mod(±Inf) = NaN
		{
			name:    "Inf number",
			modulus: 1,
			arg:     modular32.Inf(1),
			want:    modular32.NaN(),
		},
		//		Modulus{m}.Mod(NaN) = NaN
		{
			name:    "NaN number",
			modulus: 1,
			arg:     modular32.NaN(),
			want:    modular32.NaN(),
		},
		{
			name:    "Denormalised edge case",
			modulus: 1.1754944e-38,
			arg:     1.1790209e-38,
			want:    3.5265e-41,
		},
		{
			name:    "Denormalised edge case2",
			modulus: 5.877472e-39,
			arg:     5.877472e-39,
			want:    3.5265e-41,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modular32.NewModulus(tt.modulus)
			got := m.Mod(tt.arg)
			if got != tt.want && !(got != got && tt.want != tt.want) {
				t.Errorf("Modulus{%v}.Mod(%v) = %v, want %v", tt.modulus, tt.arg, got, tt.want)
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
				n1: modular32.NaN(),
				n2: 30,
			},
			want: modular32.NaN(),
		},
		{
			name:    "NaN args",
			modulus: 100,
			args: args{
				n1: 20,
				n2: modular32.NaN(),
			},
			want: modular32.NaN(),
		},
		{
			name:    "NaN modulus",
			modulus: modular32.NaN(),
			args: args{
				n1: 20,
				n2: 30,
			},
			want: modular32.NaN(),
		},
		{
			name:    "Inf modulus",
			modulus: modular32.Inf(1),
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
				n1: modular32.Inf(1),
				n2: 30,
			},
			want: modular32.NaN(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := modular32.NewModulus(tt.modulus)
			got := m.Dist(tt.args.n1, tt.args.n2)
			if got != tt.want && !(got != got && tt.want != tt.want) {
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
		got := m.Modulus()
		if got != 15 {
			t.Errorf("Modulus.Mod() = %v, want %v", got, 15)
		}
	})
}

var benchmarkModulo = float32(1e-25)
var benchmarks = []float32{
	0,
	2.5e-25,
	1,
	1e25,
}

func BenchmarkMath_Mod(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Math.Mod(%v)", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				float32Sink = math32mod(n, benchmarkModulo)
			}
		})
	}
}

func math32mod(x, y float32) float32 {
	if y == 0 || modular32.IsInf(x) || x != x || y != y {
		return modular32.NaN()
	}
	y = modular32.Abs(y)

	yfr, yexp := math32frexp(y)
	r := x
	if x < 0 {
		r = -x
	}

	for r >= y {
		rfr, rexp := math32frexp(r)
		if rfr < yfr {
			rexp = rexp - 1
		}
		r = r - math32ldexp(y, rexp-yexp)
	}
	if x < 0 {
		r = -r
	}
	return r
}

func math32ldexp(frac float32, exp int) float32 {
	// special cases
	switch {
	case frac == 0:
		return frac // correctly return -0
	case modular32.IsInf(frac) || frac != frac:
		return frac
	}
	frac, e := math32normalize(frac)
	exp += e
	x := modular32.ToBits(frac)
	exp += int(x>>23)&0xFF - 127
	if exp < -149 {
		return math32copysign(0, frac) // underflow
	}
	if exp > 127 { // overflow
		if frac < 0 {
			return modular32.Inf(-1)
		}
		return modular32.Inf(1)
	}
	var m float32 = 1
	if exp < -(127 - 1) { // denormal
		exp += 23
		m = 1.0 / (1 << 23) // 1/(2**-23)
	}
	x &^= 0xFF << 23
	x |= uint32(exp+127) << 23
	return m * modular32.FromBits(x)
}

func math32copysign(x, y float32) float32 {
	const sign = 1 << 31
	return modular32.FromBits(modular32.ToBits(x)&^sign | modular32.ToBits(y)&sign)
}

func math32frexp(f float32) (frac float32, exp int) {
	// special cases
	switch {
	case f == 0:
		return f, 0 // correctly return -0
	case modular32.IsInf(f) || f != f:
		return f, 0
	}
	f, exp = math32normalize(f)
	x := modular32.ToBits(f)
	exp += int((x>>23)&0xFF) - 127 + 1
	x &^= 0xFF << 23
	x |= (-1 + 127) << 23
	frac = modular32.FromBits(x)
	return
}

func math32normalize(x float32) (y float32, exp int) {
	const SmallestNormal = 1.1754943508222875079687365e-38 // 2**-(127 - 1)
	if modular32.Abs(x) < SmallestNormal {
		return x * (1 << 23), -23
	}
	return x, 0
}

func BenchmarkModulus(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Mod(%v)", n), func(b *testing.B) {
			m := modular32.NewModulus(benchmarkModulo)
			for i := 0; i < b.N; i++ {
				float32Sink = m.Mod(n)
			}
		})
	}
}
