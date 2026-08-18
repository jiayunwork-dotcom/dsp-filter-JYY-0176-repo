package response

import "math"

func GroupDelay(phase []float64, freq []float64) []float64 {
	n := len(phase)
	out := make([]float64, n)
	if n < 2 {
		return out
	}
	for i := 1; i < n-1; i++ {
		dPhase := phase[i+1] - phase[i-1]
		dOmega := 2 * math.Pi * (freq[i+1] - freq[i-1])
		if dOmega == 0 {
			out[i] = math.NaN()
			continue
		}
		out[i] = -dPhase / dOmega
	}
	if n >= 2 {
		dOmega0 := 2 * math.Pi * (freq[1] - freq[0])
		if dOmega0 != 0 {
			out[0] = -(phase[1] - phase[0]) / dOmega0
		}
		dOmegaN := 2 * math.Pi * (freq[n-1] - freq[n-2])
		if dOmegaN != 0 {
			out[n-1] = -(phase[n-1] - phase[n-2]) / dOmegaN
		}
	}
	return out
}

func AverageGroupDelay(gd []float64) float64 {
	if len(gd) == 0 {
		return 0
	}
	sum := 0.0
	count := 0
	for _, v := range gd {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			continue
		}
		sum += v
		count++
	}
	if count == 0 {
		return math.NaN()
	}
	return sum / float64(count)
}
