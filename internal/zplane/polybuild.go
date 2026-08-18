package zplane

import "math"

func ByRoots(roots []complex128) []float64 {
	// 从复共轭成对的根重建实系数多项式（首一，按 x 升幂）
	coeff := []float64{1}
	i := 0
	for i < len(roots) {
		r := roots[i]
		if math.Abs(imag(r)) < 1e-9 {
			coeff = mulReal(coeff, []float64{-real(r), 1})
			i++
			continue
		}
		if i+1 < len(roots) {
			conj := roots[i+1]
			if math.Abs(real(r)-real(conj)) < 1e-9 && math.Abs(imag(r)+imag(conj)) < 1e-9 {
				// (x - r)(x - conj(r)) = x^2 - 2 Re(r) x + |r|^2
				quad := []float64{real(r)*real(r) + imag(r)*imag(r), -2 * real(r), 1}
				coeff = mulReal(coeff, quad)
				i += 2
				continue
			}
		}
		coeff = mulReal(coeff, []float64{-real(r), 1})
		i++
	}
	return coeff
}

func mulReal(p, q []float64) []float64 {
	out := make([]float64, len(p)+len(q)-1)
	for i := range p {
		for j := range q {
			out[i+j] += p[i] * q[j]
		}
	}
	return out
}
