package filter

type DirectFormI struct {
	B []float64
	A []float64
	x []float64
	y []float64
}

func New(b, a []float64) *DirectFormI {
	f := &DirectFormI{B: append([]float64(nil), b...), A: append([]float64(nil), a...)}
	f.x = make([]float64, len(b))
	f.y = make([]float64, len(a))
	return f
}

func (f *DirectFormI) Reset() {
	for i := range f.x {
		f.x[i] = 0
	}
	for i := range f.y {
		f.y[i] = 0
	}
}

func (f *DirectFormI) Step(input float64) float64 {
	shiftRight(f.x)
	f.x[0] = input
	out := 0.0
	for i := range f.B {
		out += f.B[i] * f.x[i]
	}
	for i := 1; i < len(f.A); i++ {
		out -= f.A[i] * f.y[i-1]
	}
	shiftRight(f.y)
	f.y[0] = out
	return out
}

func (f *DirectFormI) Process(inputs []float64) []float64 {
	out := make([]float64, len(inputs))
	for i, x := range inputs {
		out[i] = f.Step(x)
	}
	return out
}

func shiftRight(s []float64) {
	for i := len(s) - 1; i > 0; i-- {
		s[i] = s[i-1]
	}
}
