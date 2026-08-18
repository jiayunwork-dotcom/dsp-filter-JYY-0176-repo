package response

import "math"

func EvaluateZ(b, a []float64, z complex128) complex128 {
	var num, den complex128
	zk := 1 + 0i
	for _, c := range b {
		num += complex(c, 0) * zk
		zk /= z
	}
	zk = 1 + 0i
	for _, c := range a {
		den += complex(c, 0) * zk
		zk /= z
	}
	return complexDiv(num, den)
}

func EvaluatePole(b, a []float64, pole complex128) float64 {
	if pole == 0 {
		return math.NaN()
	}
	h := EvaluateZ(b, a, pole)
	return math.Hypot(real(h), imag(h))
}
