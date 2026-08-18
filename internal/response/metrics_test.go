package response

import (
	"math"
	"testing"

	"dsp-filter/internal/design"
)

func TestBandEdge3dB(t *testing.T) {
	b, a, err := design.DesignIIR(&design.DesignSpec{Kind: design.KindIIR, Order: 4, Cutoff: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compute(b, a, Grid(256))
	if err != nil {
		t.Fatal(err)
	}
	edge, ok := BandEdge(res.Freq, res.MagnitudeDB, -3.01)
	if !ok {
		t.Fatal("no band edge found")
	}
	if math.Abs(edge-0.2) > 0.02 {
		t.Fatalf("band edge %v must sit near cutoff 0.2", edge)
	}
}

func TestInterpDB(t *testing.T) {
	freq := []float64{0, 0.25, 0.5}
	mag := []float64{0, -6, -12}
	v, ok := InterpDB(freq, mag, 0.125)
	if !ok || math.Abs(v+3) > 1e-9 {
		t.Fatalf("interp at 0.125 = %v, want -3", v)
	}
}

func TestFrequencyOfLevel(t *testing.T) {
	freq := []float64{0, 0.1, 0.2, 0.3}
	mag := []float64{0, -3, -6, -9}
	f, ok := FrequencyOfLevel(freq, mag, -6, 0)
	if !ok || math.Abs(f-0.2) > 1e-9 {
		t.Fatalf("frequency of -6 dB = %v, want 0.2", f)
	}
}

func TestEvaluatePoleBlowsUp(t *testing.T) {
	_, a, err := design.DesignIIR(&design.DesignSpec{Kind: design.KindIIR, Order: 2, Cutoff: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	b := []float64{1}
	bad := EvaluatePole(b, a, complex(0.9, 0))
	if bad < 0.1 {
		t.Fatalf("H near a pole should be large, got %v", bad)
	}
}
