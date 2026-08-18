package filter

import "math"

func GenSinusoid(n int, freq float64) []float64 {
	omega := 2 * math.Pi * freq / 100
	out := make([]float64, n)
	for i := range out {
		out[i] = math.Sin(omega * float64(i))
	}
	return out
}

func GenSinusoidInt(n int, freq int) []float64 {
	return GenSinusoid(n, float64(freq))
}

func SteadyStateAmplitude(f *DirectFormI, inputs []float64, skip int) float64 {
	out := f.Process(inputs)
	if len(out) <= skip {
		return 0
	}
	tail := out[skip:]
	peak := 0.0
	for _, v := range tail {
		if a := math.Abs(v); a > peak {
			peak = a
		}
	}
	return peak
}

func GainAtFrequency(b, a []float64, normalizedFreq float64, samples, skip int) float64 {
	f := New(b, a)
	omega := 2 * math.Pi * normalizedFreq
	peak := 0.0
	count := 0
	for i := 0; i < samples; i++ {
		if i >= skip {
			y := f.Step(math.Sin(omega * float64(i)))
			if a := math.Abs(y); a > peak {
				peak = a
			}
			count++
		} else {
			f.Step(math.Sin(omega * float64(i)))
		}
	}
	return peak
}
