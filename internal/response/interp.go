package response

func InterpDB(freq, magDB []float64, f float64) (float64, bool) {
	if len(freq) < 2 {
		return 0, false
	}
	if f <= freq[0] {
		return magDB[0], true
	}
	if f >= freq[len(freq)-1] {
		return magDB[len(magDB)-1], true
	}
	for i := 1; i < len(freq); i++ {
		if freq[i] >= f {
			f0, f1 := freq[i-1], freq[i]
			m0, m1 := magDB[i-1], magDB[i]
			frac := (f - f0) / (f1 - f0)
			return m0 + frac*(m1-m0), true
		}
	}
	return 0, false
}

func FrequencyOfLevel(freq, magDB []float64, levelDB float64, from int) (float64, bool) {
	if from < 0 {
		from = 0
	}
	for i := from + 1; i < len(freq); i++ {
		if magDB[i-1] >= levelDB && magDB[i] <= levelDB {
			f0, f1 := freq[i-1], freq[i]
			m0, m1 := magDB[i-1], magDB[i]
			if m0 == m1 {
				return f1, true
			}
			frac := (m0 - levelDB) / (m0 - m1)
			return f0 + frac*(f1-f0), true
		}
	}
	return 0, false
}
