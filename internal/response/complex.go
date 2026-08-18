package response

type ComplexResult struct {
	B     []float64     `json:"b"`
	A     []float64     `json:"a"`
	Freq  []float64     `json:"freq"`
	Real  []float64     `json:"real"`
	Imag  []float64     `json:"imag"`
}

func ComputeComplex(b, a, freq []float64) (*ComplexResult, error) {
	if err := ValidateFreq(freq); err != nil {
		return nil, err
	}
	res := &ComplexResult{B: b, A: a, Freq: freq}
	res.Real = make([]float64, len(freq))
	res.Imag = make([]float64, len(freq))
	for i, f := range freq {
		h := Evaluate(b, a, f)
		res.Real[i] = real(h)
		res.Imag[i] = imag(h)
	}
	return res, nil
}

func FromComplex(h []complex128) (re, im []float64) {
	re = make([]float64, len(h))
	im = make([]float64, len(h))
	for i, v := range h {
		re[i] = real(v)
		im[i] = imag(v)
	}
	return
}
