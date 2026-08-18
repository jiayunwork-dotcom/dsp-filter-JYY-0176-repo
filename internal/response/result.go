package response

type Result struct {
	B           []float64 `json:"b"`
	A           []float64 `json:"a"`
	Freq        []float64 `json:"freq"`
	MagnitudeDB []float64 `json:"mag_db"`
	Phase       []float64 `json:"phase"`
	GroupDelay  []float64 `json:"group_delay"`
}

func Compute(b, a, freq []float64) (*Result, error) {
	if err := ValidateFreq(freq); err != nil {
		return nil, err
	}
	if len(b) == 0 || len(a) == 0 {
		return nil, &Error{Code: ErrBadFrequency, Message: "coefficients must be non-empty"}
	}
	res := &Result{B: b, A: a, Freq: freq}
	res.MagnitudeDB = make([]float64, len(freq))
	res.Phase = make([]float64, len(freq))
	raw := make([]float64, len(freq))
	for i, f := range freq {
		h := Evaluate(b, a, f)
		res.MagnitudeDB[i] = MagnitudeDB(h)
		raw[i] = Phase(h)
	}
	res.Phase = Unwrap(raw)
	res.GroupDelay = GroupDelay(res.Phase, freq)
	return res, nil
}
