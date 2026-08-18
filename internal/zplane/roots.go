package zplane

import "math"

const (
	rootMaxIter = 300
	rootTol     = 1e-13
)

func Roots(coeff []float64) []complex128 {
	n := Deg(coeff)
	if n == 0 {
		return nil
	}
	roots := make([]complex128, n)
	for i := range roots {
		roots[i] = initialGuess(i, n)
	}
	for iter := 0; iter < rootMaxIter; iter++ {
		maxDelta := 0.0
		for i := range roots {
			pv := EvalPoly(coeff, roots[i])
			den := 1 + 0i
			for j := range roots {
				if j == i {
					continue
				}
				den *= roots[i] - roots[j]
			}
			if den == 0 {
				continue
			}
			delta := pv / den
			roots[i] -= delta
			if d := absC(delta); d > maxDelta {
				maxDelta = d
			}
		}
		if maxDelta < rootTol {
			break
		}
	}
	return roots
}

func initialGuess(i, n int) complex128 {
	r := 0.4 + 0.9*float64(i%2)
	theta := 2 * math.Pi * float64(i) / float64(n)
	return complex(r*math.Cos(theta), r*math.Sin(theta))
}

func absC(z complex128) float64 {
	return math.Hypot(real(z), imag(z))
}
