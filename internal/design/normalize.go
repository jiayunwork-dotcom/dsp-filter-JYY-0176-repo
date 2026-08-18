package design

func NormalizeDC(b, a []float64) []float64 {
	gain := dcGain(b, a)
	if gain == 0 {
		return append([]float64(nil), b...)
	}
	out := make([]float64, len(b))
	for i, v := range b {
		out[i] = v / gain
	}
	return out
}

func dcGain(b, a []float64) float64 {
	bSum := 0.0
	for _, v := range b {
		bSum += v
	}
	aSum := 0.0
	for _, v := range a {
		aSum += v
	}
	if aSum == 0 {
		return 0
	}
	return bSum / aSum
}

func MaxCoeff(coeff []float64) float64 {
	max := 0.0
	for _, v := range coeff {
		if av := absF(v); av > max {
			max = av
		}
	}
	return max
}

func absF(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
