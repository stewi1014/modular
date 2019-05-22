package modular64

import "testing"

func Test_pow2mod(t *testing.T) {
	type args struct {
		exp     uint
		modulus uint64
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "Basic test",
			args: args{
				exp:     5,
				modulus: 12,
			},
			want: 8,
		},
		{
			name: "2**64 test",
			args: args{
				exp:     64,
				modulus: 12,
			},
			want: 4,
		},
		{
			name: "Large test",
			args: args{
				exp:     128,
				modulus: 36754,
			},
			want: 18344,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pow2mod(tt.args.exp, tt.args.modulus); got != tt.want {
				t.Errorf("pow2mod(%v, %v) = %v, want %v", tt.args.exp, tt.args.modulus, got, tt.want)
			}
		})
	}
}
