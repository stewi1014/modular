package integeru64

import "testing"

// TODO; add more tests
func TestSqrt(t *testing.T) {
	tests := []struct {
		name          string
		arg           uint64
		wantSqrt      uint64
		wantRemainder uint64
	}{
		{
			name:          "sqrt(4) = 2,0",
			arg:           4,
			wantSqrt:      2,
			wantRemainder: 0,
		},
		{
			name:          "sqrt(9) = 3,0",
			arg:           9,
			wantSqrt:      3,
			wantRemainder: 0,
		},
		{
			name:          "sqrt(10) = 3,1",
			arg:           10,
			wantSqrt:      3,
			wantRemainder: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSqrt, gotRemainder := Sqrt(tt.arg)
			if gotSqrt != tt.wantSqrt {
				t.Errorf("Sqrt() gotSqrt = %v, want %v", gotSqrt, tt.wantSqrt)
			}
			if gotRemainder != tt.wantRemainder {
				t.Errorf("Sqrt() gotRemainder = %v, want %v", gotRemainder, tt.wantRemainder)
			}
		})
	}
}
