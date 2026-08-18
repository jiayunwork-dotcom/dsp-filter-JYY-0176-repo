package zplane

func Deg(coeff []float64) int {
	n := len(coeff) - 1
	for n > 0 && coeff[n] == 0 {
		n--
	}
	return n
}

func EvalPoly(coeff []float64, x complex128) complex128 {
	acc := 0 + 0i
	for i := len(coeff) - 1; i >= 0; i-- {
		acc = acc*x + complex(coeff[i], 0)
	}
	return acc
}

func EvalPolyReal(coeff []float64, x float64) float64 {
	acc := 0.0
	for i := len(coeff) - 1; i >= 0; i-- {
		acc = acc*x + coeff[i]
	}
	return acc
}

func TrimZeros(coeff []float64) []float64 {
	n := Deg(coeff)
	if n+1 == len(coeff) {
		return coeff
	}
	out := make([]float64, n+1)
	copy(out, coeff[:n+1])
	return out
}

func Reversed(coeff []float64) []float64 {
	out := make([]float64, len(coeff))
	for i, v := range coeff {
		out[len(coeff)-1-i] = v
	}
	return out
}
