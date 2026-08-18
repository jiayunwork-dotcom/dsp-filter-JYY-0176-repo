package design

import "math"

func DesignFIR(spec *DesignSpec) ([]float64, error) {
	if err := Validate(spec); err != nil {
		return nil, err
	}
	win, _ := ParseWindow(spec.Window)
	n := spec.Order
	taps := n + 1
	mid := float64(n) / 2
	b := make([]float64, taps)
	for i := 0; i < taps; i++ {
		b[i] = idealLowpass(float64(i)-mid, spec.Cutoff) * win.Sample(i, taps)
	}
	return b, nil
}

func idealLowpass(x, fc float64) float64 {
	if x == 0 {
		return 2 * fc
	}
	return math.Sin(2*math.Pi*fc*x) / (math.Pi * x)
}

func SymmetricCheck(b []float64, tol float64) bool {
	n := len(b)
	for i := 0; i < n/2; i++ {
		if math.Abs(b[i]-b[n-1-i]) > tol {
			return false
		}
	}
	return true
}
