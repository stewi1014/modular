package modular64

import (
	"math"
	"testing"
)

var (
	varIndex uint = 56
)

func TestIndexer_Index(t *testing.T) {
	type args struct {
		modulus float64
		index   uint
		n       float64
	}
	type want struct {
		n   uint
		err error
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
				n:   1,
				err: nil,
			},
		},
		{
			name: "Different Index",
			args: args{
				modulus: 15,
				index:   10,
				n:       1.5,
			},
			want: want{
				n:   1,
				err: nil,
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
				n:   99,
				err: nil,
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
				n:   12,
				err: nil,
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
				n:   0,
				err: ErrBadModulo,
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
				n:   0,
				err: ErrBadModulo,
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
				n:   0,
				err: ErrIndexTooBig,
			},
		},
		{
			name: "Edge case with number=modulo",
			args: args{
				modulus: 1.4510462197599293e+120,
				index:   4384042,
				n:       -1.4510462197599290e-140,
			},
			want: want{
				n:   4384041,
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := NewIndexer(tt.args.modulus, tt.args.index)
			if got := i.Index(tt.args.n); got != tt.want.n || err != tt.want.err {
				t.Errorf("Indexer.Index(%v) = %v, want %v\nNewIndex error: \"%v\", want \"%v\"; Modulus: %v; Index: %v", tt.args.n, got, tt.want.n, err, tt.want.err, tt.args.modulus, tt.args.index)
			}
		})
	}
}

func BenchmarkIndexer_Index(b *testing.B) {
	ind, _ := NewIndexer(varMod, varIndex)
	for i := 0; i < b.N; i++ {
		varUintSink = ind.Index(varNumber)
	}
}
