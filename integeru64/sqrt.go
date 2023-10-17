package integeru64

// Sqrt returns the square root of n, rounding down and returning the remainder.
func Sqrt(n uint64) (sqrt uint64, remainder uint64) {
	// TODO; optimise

	if n < 2 {
		return n, 0
	}

	shift := 2
	for (n >> shift) != 0 {
		shift += 2
	}

	result := uint64(0)
	for shift >= 0 {
		result = result << 1

		large_cand := result + 1

		if large_cand*large_cand <= n>>shift {
			result = large_cand
		}

		shift -= 2
	}

	return result, n - (result * result)
}
