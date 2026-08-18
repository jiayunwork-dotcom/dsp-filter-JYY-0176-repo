package response

import (
	"math"
	"testing"

	"dsp-filter/internal/design"
)

func TestButterworthCutoff3dB(t *testing.T) {
	b, a, err := design.DesignIIR(&design.DesignSpec{Kind: design.KindIIR, Order: 4, Cutoff: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	freq := []float64{0.2}
	res, err := Compute(b, a, freq)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(res.MagnitudeDB[0]+3.01) > 0.5 {
		t.Fatalf("magnitude at prewarped cutoff must be near -3 dB, got %v", res.MagnitudeDB[0])
	}
}

func TestResponseEmptyGrid(t *testing.T) {
	if _, err := Compute([]float64{1}, []float64{1}, nil); !IsError(err, ErrEmptyGrid) {
		t.Fatalf("empty grid must be rejected, got %v", err)
	}
}

func TestResponseNaN(t *testing.T) {
	freq := []float64{0, math.NaN(), 0.3}
	if _, err := Compute([]float64{1}, []float64{1}, freq); !IsError(err, ErrBadFrequency) {
		t.Fatalf("NaN grid must be rejected, got %v", err)
	}
}

func TestResponseOutOfRangeFreq(t *testing.T) {
	freq := []float64{0, 0.6}
	if _, err := Compute([]float64{1}, []float64{1}, freq); !IsError(err, ErrBadFrequency) {
		t.Fatalf("out of range grid must be rejected, got %v", err)
	}
}

func TestFIRGroupDelayMatchesN2(t *testing.T) {
	order := 30
	b, err := design.DesignFIR(&design.DesignSpec{Kind: design.KindFIR, Order: order, Cutoff: 0.2, Window: "hamming"})
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compute(b, []float64{1}, Grid(256))
	if err != nil {
		t.Fatal(err)
	}
	sum := 0.0
	count := 0
	for i, f := range res.Freq {
		if f < 0.19 {
			sum += res.GroupDelay[i]
			count++
		}
	}
	avg := sum / float64(count)
	if math.Abs(avg-float64(order)/2) > 0.05 {
		t.Fatalf("fir passband group delay must equal N/2 = %v, got %v", float64(order)/2, avg)
	}
}

func TestPhaseUnwrapSmooth(t *testing.T) {
	b, a, err := design.DesignIIR(&design.DesignSpec{Kind: design.KindIIR, Order: 4, Cutoff: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compute(b, a, Grid(128))
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i < len(res.Phase); i++ {
		if math.Abs(res.Phase[i]-res.Phase[i-1]) > math.Pi {
			t.Fatalf("unwrapped phase must be continuous at %d: %v -> %v", i, res.Phase[i-1], res.Phase[i])
		}
	}
}

func TestGridUniform(t *testing.T) {
	g := Grid(5)
	if len(g) != 5 {
		t.Fatalf("grid length = %d, want 5", len(g))
	}
	want := []float64{0, 0.125, 0.25, 0.375, 0.5}
	for i := range g {
		if math.Abs(g[i]-want[i]) > 1e-12 {
			t.Fatalf("grid[%d] = %v, want %v", i, g[i], want[i])
		}
	}
}

func TestFIRPassbandFlat(t *testing.T) {
	b, err := design.DesignFIR(&design.DesignSpec{Kind: design.KindFIR, Order: 60, Cutoff: 0.3, Window: "hamming"})
	if err != nil {
		t.Fatal(err)
	}
	res, err := Compute(b, []float64{1}, []float64{0.001, 0.1, 0.2})
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range res.MagnitudeDB {
		if math.Abs(v) > 0.2 {
			t.Fatalf("passband ripple too large: %v dB", v)
		}
	}
}
