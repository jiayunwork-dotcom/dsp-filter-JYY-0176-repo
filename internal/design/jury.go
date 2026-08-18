package design

import "math"

func JuryStable(a []float64) bool {
	n := len(a) - 1
	if n == 0 {
		return true
	}
	if evalZPoly(a, 1) <= 0 {
		return false
	}
	sign := 1.0
	if n%2 == 1 {
		sign = -1
	}
	if sign*evalZPoly(a, -1) <= 0 {
		return false
	}
	if math.Abs(a[0]) <= math.Abs(a[n]) {
		return false
	}
	row := make([]float64, n+1)
	copy(row, a)
	for m := n; m > 1; m-- {
		alpha := row[m] / row[0]
		next := make([]float64, m)
		for i := 0; i < m; i++ {
			next[i] = row[i] - alpha*row[m-i]
		}
		if next[0] <= 0 || math.Abs(next[0]) <= math.Abs(next[m-1]) {
			return false
		}
		row = next
	}
	return true
}

func evalZPoly(a []float64, z float64) float64 {
	acc := a[0]
	for i := 1; i < len(a); i++ {
		acc = acc*z + a[i]
	}
	return acc
}
