package response

import (
	"math"
	"testing"
)

func TestComputeComplex(t *testing.T) {
	res, err := ComputeComplex([]float64{1}, []float64{1}, []float64{0, 0.25, 0.5})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Real) != 3 || len(res.Imag) != 3 {
		t.Fatalf("complex result length mismatch")
	}
	if math.Abs(res.Real[0]-1) > 1e-12 || math.Abs(res.Imag[0]) > 1e-12 {
		t.Fatalf("H(0) must be 1, got %v + %v j", res.Real[0], res.Imag[0])
	}
}

func TestFromComplex(t *testing.T) {
	h := []complex128{complex(1, 2), complex(3, -4)}
	re, im := FromComplex(h)
	if re[0] != 1 || im[0] != 2 || re[1] != 3 || im[1] != -4 {
		t.Fatalf("from complex mismatch: %v %v", re, im)
	}
}
