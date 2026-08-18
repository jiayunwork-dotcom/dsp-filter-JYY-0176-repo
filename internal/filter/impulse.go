package filter

func ImpulseResponse(f *DirectFormI, n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = f.Step(impulseAt(i))
	}
	return out
}

func impulseAt(i int) float64 {
	if i == 0 {
		return 1
	}
	return 0
}

func StepResponse(f *DirectFormI, n int) []float64 {
	out := make([]float64, n)
	for i := range out {
		out[i] = f.Step(1)
	}
	return out
}

func SteadyGain(b, a []float64) float64 {
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

func SettledAt(values []float64, tol float64) int {
	if len(values) == 0 {
		return -1
	}
	last := values[len(values)-1]
	for i := len(values) - 1; i >= 0; i-- {
		if abs(values[i]-last) > tol {
			return i + 1
		}
	}
	return 0
}

func abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
