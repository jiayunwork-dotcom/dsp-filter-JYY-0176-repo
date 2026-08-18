package zplane

type Complex struct {
	Re float64 `json:"re"`
	Im float64 `json:"im"`
}

type ZPResult struct {
	Type   string    `json:"type"`
	Zeros  []Complex `json:"zeros"`
	Poles  []Complex `json:"poles"`
	Stable bool      `json:"stable"`
}

const StabilityEps = 1e-6

func ZeroPoles(b, a []float64) *ZPResult {
	zeros := Roots(Reversed(TrimZeros(b)))
	poles := Roots(Reversed(TrimZeros(a)))
	res := &ZPResult{
		Type:   classifyType(a),
		Zeros:  toComplex(zeros),
		Poles:  toComplex(poles),
		Stable: AllInside(poles, StabilityEps),
	}
	return res
}

func classifyType(a []float64) string {
	if len(TrimZeros(a)) == 1 {
		return "fir"
	}
	return "iir"
}

func toComplex(roots []complex128) []Complex {
	out := make([]Complex, len(roots))
	for i, r := range roots {
		out[i] = Complex{Re: real(r), Im: imag(r)}
	}
	return out
}

func AllInside(poles []complex128, eps float64) bool {
	for _, p := range poles {
		if absC(p) > 1+eps {
			return false
		}
	}
	return true
}
