package zplane

import (
	"math"
	"testing"
)

func TestByRootsRoundTrip(t *testing.T) {
	roots := []complex128{
		complex(0.5, 0.4), complex(0.5, -0.4),
		complex(-0.7, 0),
	}
	coeff := ByRoots(roots)
	got := Roots(coeff)
	if len(got) != 3 {
		t.Fatalf("round trip root count %d, want 3", len(got))
	}
	for _, g := range got {
		match := false
		for _, r := range roots {
			if math.Abs(real(g)-real(r)) < 1e-6 && math.Abs(imag(g)-imag(r)) < 1e-6 {
				match = true
			}
		}
		if !match {
			t.Fatalf("recovered root %v not in expected set", g)
		}
	}
}

func TestByRootsRealCoefficients(t *testing.T) {
	roots := []complex128{complex(0.5, 0.4), complex(0.5, -0.4), complex(0.2, 0)}
	coeff := ByRoots(roots)
	for _, c := range coeff {
		if math.Abs(c-math.Round(c*1e9)/1e9) > 1e-12 {
		}
	}
	if len(coeff) != 4 {
		t.Fatalf("quadratic*linear must give degree 3, got length %d", len(coeff))
	}
}

func TestUnitCirclePoints(t *testing.T) {
	pts := UnitCirclePoints(8)
	if len(pts) != 8 {
		t.Fatalf("want 8 points, got %d", len(pts))
	}
	for _, p := range pts {
		r := math.Hypot(p.Re, p.Im)
		if math.Abs(r-1) > 1e-12 {
			t.Fatalf("unit circle point radius %v", r)
		}
	}
}

func TestToPolar(t *testing.T) {
	p := ToPolar(Complex{Re: 1, Im: 1})
	if math.Abs(p.Radius-math.Sqrt2) > 1e-12 {
		t.Fatalf("radius = %v", p.Radius)
	}
	if math.Abs(p.Angle-math.Pi/4) > 1e-12 {
		t.Fatalf("angle = %v", p.Angle)
	}
}
