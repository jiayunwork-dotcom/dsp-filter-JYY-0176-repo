package response

func BandEdge(freq, magDB []float64, targetDB float64) (float64, bool) {
	if len(freq) < 2 {
		return 0, false
	}
	if magDB[0] <= targetDB {
		return freq[0], true
	}
	for i := 1; i < len(freq); i++ {
		if magDB[i] <= targetDB {
			f0, f1 := freq[i-1], freq[i]
			m0, m1 := magDB[i-1], magDB[i]
			if m0 == m1 {
				return f1, true
			}
			frac := (m0 - targetDB) / (m0 - m1)
			return f0 + frac*(f1-f0), true
		}
	}
	return freq[len(freq)-1], false
}

func PassbandRipple(magDB []float64, upTo int) float64 {
	if upTo > len(magDB) {
		upTo = len(magDB)
	}
	if upTo <= 1 {
		return 0
	}
	maxV := magDB[0]
	minV := magDB[0]
	for i := 1; i < upTo; i++ {
		if magDB[i] > maxV {
			maxV = magDB[i]
		}
		if magDB[i] < minV {
			minV = magDB[i]
		}
	}
	return maxV - minV
}

func StopbandAttenuation(magDB []float64, from int) float64 {
	if from >= len(magDB) {
		return 0
	}
	maxV := magDB[from]
	for i := from; i < len(magDB); i++ {
		if magDB[i] > maxV {
			maxV = magDB[i]
		}
	}
	return maxV
}
