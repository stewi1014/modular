package modular32_test

import (
	"fmt"
	"testing"

	math "github.com/chewxy/math32"
	"github.com/stewi1014/modular/modular32"
)

func ExampleIndexer() {
	shifts := []string{
		"morning",
		"day",
		"evening",
	}

	// Errors can be ignored so long as we don't feed bad numbers
	indexer, _ := modular32.NewIndexer(24, len(shifts))

	for i := float32(0); i < 100; i += 13 {
		shift := shifts[indexer.Index(i)]
		fmt.Printf("It will be the %v shift in %v hours\n", shift, i)
	}

	// Output:
	// It will be the morning shift in 0 hours
	// It will be the day shift in 13 hours
	// It will be the morning shift in 26 hours
	// It will be the day shift in 39 hours
	// It will be the morning shift in 52 hours
	// It will be the evening shift in 65 hours
	// It will be the morning shift in 78 hours
	// It will be the evening shift in 91 hours
}

var (
	intSink int
)

func TestIndexer_Index(t *testing.T) {
	type args struct {
		modulus float32
		index   int
		n       float32
	}
	type want struct {
		n           int
		creationErr error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Basic test",
			args: args{
				modulus: 15,
				index:   15,
				n:       1,
			},
			want: want{
				n:           1,
				creationErr: nil,
			},
		},
		{
			name: "Different Index and Modulo",
			args: args{
				modulus: 15,
				index:   10,
				n:       1.5,
			},
			want: want{
				n:           1,
				creationErr: nil,
			},
		},
		{
			name: "Number and Modulo same exponent",
			args: args{
				modulus: 120,
				index:   100,
				n:       115,
			},
			want: want{
				n:           95,
				creationErr: nil,
			},
		},
		{
			name: "Negative Number",
			args: args{
				modulus: 200,
				index:   100,
				n:       -2,
			},
			want: want{
				n:           99,
				creationErr: nil,
			},
		},
		{
			name: "Negative Number larger than modulus",
			args: args{
				modulus: 200,
				index:   100,
				n:       -202,
			},
			want: want{
				n:           99,
				creationErr: nil,
			},
		},
		{
			name: "Large number",
			args: args{
				modulus: 10,
				index:   20,
				n:       98723456,
			},
			want: want{
				n:           12,
				creationErr: nil,
			},
		},
		{
			name: "Edge case with number=modulo",
			args: args{
				modulus: 1e+20,
				index:   43848,
				n:       -1e-30,
			},
			want: want{
				n:           43847,
				creationErr: nil,
			},
		},
		{
			name: "Infinite Modulus",
			args: args{
				modulus: math.Inf(1),
				index:   100,
				n:       -2,
			},
			want: want{
				n:           0,
				creationErr: modular32.ErrBadModulo,
			},
		},
		{
			name: "NaN Modulus",
			args: args{
				modulus: math.NaN(),
				index:   10054,
				n:       -2,
			},
			want: want{
				n:           0,
				creationErr: modular32.ErrBadModulo,
			},
		},
		{
			name: "NaN Number",
			args: args{
				modulus: 23,
				index:   10054,
				n:       math.NaN(),
			},
			want: want{
				n:           10054,
				creationErr: nil,
			},
		},
		{
			name: "Index too big",
			args: args{
				modulus: 1.4510462197599293,
				index:   438404225733485,
				n:       1.4510462197599290,
			},
			want: want{
				n:           0,
				creationErr: modular32.ErrBadIndex,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := modular32.NewIndexer(tt.args.modulus, tt.args.index)
			if got := i.Index(tt.args.n); got != tt.want.n || err != tt.want.creationErr {
				t.Errorf("Indexer.Index(%v) = %v, want %v\nNewIndex error: \"%v\", want \"%v\"; Modulus: %v; Index: %v", tt.args.n, got, tt.want.n, err, tt.want.creationErr, tt.args.modulus, tt.args.index)
			}
		})
	}
}

func BenchmarkIndexer(b *testing.B) {
	for _, n := range benchmarks {
		b.Run(fmt.Sprintf("Indexer.Index(%v)", n), func(b *testing.B) {
			ind, _ := modular32.NewIndexer(benchmarkModulo, 100)
			for i := 0; i < b.N; i++ {
				intSink = ind.Index(n)
			}
		})
	}
}
