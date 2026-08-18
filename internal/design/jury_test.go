package design

import "testing"

func TestJuryStablePoleInside(t *testing.T) {
	if !JuryStable([]float64{1, -0.5}) {
		t.Fatal("pole at 0.5 must be stable")
	}
	if !JuryStable([]float64{1, 0, 0.25}) {
		t.Fatal("poles at +-0.5i must be stable")
	}
}

func TestJuryUnstable(t *testing.T) {
	if JuryStable([]float64{1, -2.1}) {
		t.Fatal("pole at 2.1 must be unstable")
	}
	if JuryStable([]float64{1, -3, 2}) {
		t.Fatal("poles at 1 and 2 must be unstable")
	}
}

func TestJuryButterworth(t *testing.T) {
	_, a, err := DesignIIR(&DesignSpec{Kind: KindIIR, Order: 4, Cutoff: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	if !JuryStable(a) {
		t.Fatal("butterworth denominator must be stable")
	}
}

func TestNormalizeDC(t *testing.T) {
	_, a, _ := DesignIIR(&DesignSpec{Kind: KindIIR, Order: 2, Cutoff: 0.2})
	b, _ := DesignFIR(&DesignSpec{Kind: KindFIR, Order: 30, Cutoff: 0.2, Window: "hamming"})
	norm := NormalizeDC(b, a)
	g := dcGain(norm, a)
	if g > 1.0001 || g < 0.9999 {
		t.Fatalf("normalized dc gain = %v, want 1", g)
	}
}

func TestCheckCoeffs(t *testing.T) {
	if err := CheckCoeffs([]float64{1, 2}, []float64{1}); err != nil {
		t.Fatal(err)
	}
	if err := CheckCoeffs(nil, []float64{1}); err == nil {
		t.Fatal("empty b must be rejected")
	}
	if err := CheckCoeffs([]float64{1}, []float64{0, 1}); err == nil {
		t.Fatal("a[0]=0 must be rejected")
	}
}

func TestBlackmanBartlettWindows(t *testing.T) {
	for _, name := range []string{"blackman", "bartlett"} {
		b, err := DesignFIR(&DesignSpec{Kind: KindFIR, Order: 30, Cutoff: 0.2, Window: name})
		if err != nil {
			t.Fatalf("window %s: %v", name, err)
		}
		if !SymmetricCheck(b, 1e-9) {
			t.Fatalf("window %s must stay type-I symmetric", name)
		}
	}
}
