package design

import "math"

func RoundCoeffs(coeff []float64, digits int) []float64 {
	factor := math.Pow(10, float64(digits))
	out := make([]float64, len(coeff))
	for i, v := range coeff {
		out[i] = math.Round(v*factor) / factor
	}
	return out
}

func MaxAbsDiff(a, b []float64) float64 {
	n := minInt(len(a), len(b))
	max := 0.0
	for i := 0; i < n; i++ {
		if d := math.Abs(a[i] - b[i]); d > max {
			max = d
		}
	}
	return max
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
