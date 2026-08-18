package response

import "math"

func Evaluate(b, a []float64, freq float64) complex128 {
	omega := 2 * math.Pi * freq
	return complexDiv(polyEval(b, omega), polyEval(a, omega))
}

func polyEval(coeff []float64, omega float64) complex128 {
	var acc complex128
	for k, c := range coeff {
		acc += complex(c*math.Cos(-omega*float64(k)), c*math.Sin(-omega*float64(k)))
	}
	return acc
}

func complexDiv(num, den complex128) complex128 {
	if den == 0 {
		return complex(math.NaN(), math.NaN())
	}
	d := real(den)*real(den) + imag(den)*imag(den)
	return complex((real(num)*real(den)+imag(num)*imag(den))/d,
		(imag(num)*real(den)-real(num)*imag(den))/d)
}

func MagnitudeDB(h complex128) float64 {
	m := math.Hypot(real(h), imag(h))
	if m == 0 {
		return -math.Inf(1)
	}
	return 20 * math.Log10(m)
}

func Phase(h complex128) float64 {
	return math.Atan2(imag(h), real(h))
}
