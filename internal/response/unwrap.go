package response

import "math"

func Unwrap(phase []float64) []float64 {
	out := make([]float64, len(phase))
	if len(phase) == 0 {
		return out
	}
	out[0] = phase[0]
	for i := 1; i < len(phase); i++ {
		delta := phase[i] - phase[i-1]
		out[i] = out[i-1] + wrapDelta(delta)
	}
	return out
}

func wrapDelta(delta float64) float64 {
	for delta > math.Pi {
		delta -= 2 * math.Pi
	}
	for delta < -math.Pi {
		delta += 2 * math.Pi
	}
	return delta
}

func PhaseSlope(phase []float64, freq []float64) float64 {
	if len(phase) < 2 {
		return 0
	}
	last := len(phase) - 1
	dPhase := phase[last] - phase[0]
	dOmega := 2 * math.Pi * (freq[last] - freq[0])
	if dOmega == 0 {
		return 0
	}
	return dPhase / dOmega
}
