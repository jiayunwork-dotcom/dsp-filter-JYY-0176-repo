package zplane

import (
	"math"
	"testing"

	"dsp-filter/internal/design"
)

func TestFIRZerosCount(t *testing.T) {
	b, err := design.DesignFIR(&design.DesignSpec{Kind: design.KindFIR, Order: 30, Cutoff: 0.2, Window: "hamming"})
	if err != nil {
		t.Fatal(err)
	}
	zp := ZeroPoles(b, []float64{1})
	if len(zp.Zeros) != 30 {
		t.Fatalf("order 30 fir must have 30 zeros, got %d", len(zp.Zeros))
	}
	if !zp.Stable {
		t.Fatal("fir must be stable")
	}
}

func TestButterworthStable(t *testing.T) {
	_, a, err := design.DesignIIR(&design.DesignSpec{Kind: design.KindIIR, Order: 4, Cutoff: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	zp := ZeroPoles([]float64{1}, a)
	if len(zp.Poles) != 4 {
		t.Fatalf("order 4 iir must have 4 poles, got %d", len(zp.Poles))
	}
	if !zp.Stable {
		t.Fatal("butterworth poles must be inside the unit circle")
	}
}

func TestUnstableDetection(t *testing.T) {
	// y[n] = 2.1*y[n-1] + x[n]  => pole at 2.1
	zp := ZeroPoles([]float64{1}, []float64{1, -2.1})
	if zp.Stable {
		t.Fatal("pole at 2.1 must be flagged unstable")
	}
}

func TestRootsKnownQuadratic(t *testing.T) {
	// x^2 - 3x + 2 = (x-1)(x-2)
	roots := Roots([]float64{2, -3, 1})
	if len(roots) != 2 {
		t.Fatalf("expected 2 roots, got %d", len(roots))
	}
	seen := map[float64]bool{}
	for _, r := range roots {
		if math.Abs(imag(r)) > 1e-6 {
			t.Fatalf("roots must be real: %v", r)
		}
		seen[math.Round(real(r))] = true
	}
	if !seen[1] || !seen[2] {
		t.Fatalf("roots must be 1 and 2, got %v", roots)
	}
}

func TestPolyDegTrim(t *testing.T) {
	if Deg([]float64{1, 2, 3, 0}) != 2 {
		t.Fatalf("deg of [1 2 3 0] must be 2")
	}
	if Deg([]float64{0}) != 0 {
		t.Fatalf("deg of [0] must be 0")
	}
	trimmed := TrimZeros([]float64{1, 2, 3, 0, 0})
	if len(trimmed) != 3 {
		t.Fatalf("trimmed length = %d, want 3", len(trimmed))
	}
}

func TestPoleClassification(t *testing.T) {
	poles := []complex128{
		complex(0.5, 0), complex(1.5, 0), complex(1, 0),
		complex(0.3, 0.4), complex(0.3, -0.4),
	}
	c := Classify(poles, 1e-6)
	if c.Inside != 3 || c.Outside != 1 || c.OnUnit != 1 {
		t.Fatalf("classification wrong: %+v", c)
	}
	if c.ConjugatePairs != 1 {
		t.Fatalf("expected 1 conjugate pair, got %d", c.ConjugatePairs)
	}
	if c.Real != 3 {
		t.Fatalf("expected 3 real poles, got %d", c.Real)
	}
}

func TestMaxPoleRadius(t *testing.T) {
	poles := []Complex{{Re: 0, Im: 0}, {Re: 3, Im: 4}}
	if r := MaxPoleRadius(poles); math.Abs(r-5) > 1e-9 {
		t.Fatalf("max pole radius = %v, want 5", r)
	}
}
