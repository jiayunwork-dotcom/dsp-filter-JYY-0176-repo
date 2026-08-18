package filter

func Energy(h []float64) float64 {
	sum := 0.0
	for _, v := range h {
		sum += v * v
	}
	return sum
}

func TailEnergyRatio(h []float64, cutoff int) float64 {
	if len(h) == 0 {
		return 0
	}
	total := Energy(h)
	if total == 0 {
		return 0
	}
	tail := 0.0
	for i := cutoff; i < len(h); i++ {
		tail += h[i] * h[i]
	}
	return tail / total
}

func MaxAbs(h []float64) float64 {
	max := 0.0
	for _, v := range h {
		if a := abs(v); a > max {
			max = a
		}
	}
	return max
}
