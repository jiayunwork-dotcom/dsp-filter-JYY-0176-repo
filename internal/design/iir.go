package design

import "math"

func DesignIIR(spec *DesignSpec) (b, a []float64, err error) {
	if err := Validate(spec); err != nil {
		return nil, nil, err
	}
	omega := Prewarp(spec.Cutoff)
	den := analogButterworth(spec.Order, omega)
	bn, an := bilinearDenominator(den, spec.Order)
	scale := math.Pow(omega, float64(spec.Order))
	for i := range bn {
		bn[i] *= scale
	}
	a0 := an[0]
	for i := range an {
		an[i] /= a0
	}
	for i := range bn {
		bn[i] /= a0
	}
	return bn, an, nil
}

func Prewarp(fc float64) float64 {
	return 2 * math.Tan(math.Pi*fc)
}

func analogButterworth(n int, omega float64) []complex128 {
	// A(s) = prod(s - pk), pk = omega * exp(j*pi*(2k+n-1)/(2n))
	coeff := []complex128{1}
	for k := 1; k <= n; k++ {
		angle := math.Pi * float64(2*k+n-1) / float64(2*n)
		pk := complex(omega*math.Cos(angle), omega*math.Sin(angle))
		next := make([]complex128, len(coeff)+1)
		for i, c := range coeff {
			next[i+1] += c
			next[i] -= c * pk
		}
		coeff = next
	}
	return coeff
}

func bilinearDenominator(analog []complex128, n int) ([]float64, []float64) {
	// H(z) = (1+z^-1)^n / sum_k c_k * 2^k * (1-z^-1)^k * (1+z^-1)^(n-k)
	num := powPoly([]float64{1, 1}, n)
	den := make([]float64, n+1)
	for k := 0; k <= n; k++ {
		ck := real(analog[k])
		if ck == 0 {
			continue
		}
		term := mulPoly(powPoly([]float64{1, -1}, k), powPoly([]float64{1, 1}, n-k))
		scale := ck * math.Pow(2, float64(k))
		for i := range term {
			den[i] += scale * term[i]
		}
	}
	trim := len(den)
	for trim > 1 && math.Abs(den[trim-1]) < 1e-15 {
		trim--
	}
	return num, den[:trim]
}

func NormalizeA(a []float64) []float64 {
	if len(a) == 0 || a[0] == 0 {
		return a
	}
	out := make([]float64, len(a))
	for i, v := range a {
		out[i] = v / a[0]
	}
	return out
}

func mulPoly(p, q []float64) []float64 {
	out := make([]float64, len(p)+len(q)-1)
	for i := range p {
		for j := range q {
			out[i+j] += p[i] * q[j]
		}
	}
	return out
}

func powPoly(base []float64, n int) []float64 {
	out := []float64{1}
	for i := 0; i < n; i++ {
		out = mulPoly(out, base)
	}
	return out
}
