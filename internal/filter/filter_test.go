package filter

import (
	"math"
	"testing"

	"dsp-filter/internal/design"
)

func TestFIRImpulseResponseMatchesCoeffs(t *testing.T) {
	b, err := design.DesignFIR(&design.DesignSpec{Kind: design.KindFIR, Order: 4, Cutoff: 0.2, Window: "rect"})
	if err != nil {
		t.Fatal(err)
	}
	f := New(b, []float64{1})
	h := ImpulseResponse(f, len(b))
	for i := range b {
		if math.Abs(h[i]-b[i]) > 1e-12 {
			t.Fatalf("fir impulse response[%d] = %v, want %v", i, h[i], b[i])
		}
	}
}

func TestFIRStepSettlesToDCGain(t *testing.T) {
	b, err := design.DesignFIR(&design.DesignSpec{Kind: design.KindFIR, Order: 30, Cutoff: 0.2, Window: "hamming"})
	if err != nil {
		t.Fatal(err)
	}
	f := New(b, []float64{1})
	steps := StepResponse(f, 64)
	gain := SteadyGain(b, []float64{1})
	last := steps[len(steps)-1]
	if math.Abs(last-gain) > 1e-6 {
		t.Fatalf("fir step steady state %v must equal dc gain %v", last, gain)
	}
}

func TestIIRStepSettlesToDCGain(t *testing.T) {
	b, a, err := design.DesignIIR(&design.DesignSpec{Kind: design.KindIIR, Order: 4, Cutoff: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	f := New(b, a)
	steps := StepResponse(f, 512)
	gain := SteadyGain(b, a)
	last := steps[len(steps)-1]
	if math.Abs(last-gain) > 1e-3 {
		t.Fatalf("iir step steady state %v must equal dc gain %v", last, gain)
	}
}

func TestIIREnergyDecays(t *testing.T) {
	b, a, err := design.DesignIIR(&design.DesignSpec{Kind: design.KindIIR, Order: 4, Cutoff: 0.2})
	if err != nil {
		t.Fatal(err)
	}
	f := New(b, a)
	h := ImpulseResponse(f, 2048)
	if TailEnergyRatio(h, 128) > 0.01 {
		t.Fatalf("stable iir impulse tail must decay, ratio %v", TailEnergyRatio(h, 128))
	}
}

func TestResetRestoresState(t *testing.T) {
	b, _ := design.DesignFIR(&design.DesignSpec{Kind: design.KindFIR, Order: 4, Cutoff: 0.2, Window: "rect"})
	f := New(b, []float64{1})
	f.Step(1)
	f.Step(2)
	f.Reset()
	if out := f.Step(1); math.Abs(out-b[0]) > 1e-12 {
		t.Fatalf("after reset step(1) must equal b[0], got %v", out)
	}
}

func TestProcessBlockMatchesStep(t *testing.T) {
	b, _ := design.DesignFIR(&design.DesignSpec{Kind: design.KindFIR, Order: 4, Cutoff: 0.2, Window: "rect"})
	inputs := []float64{1, 2, 3, 4, 5}
	f1 := New(b, []float64{1})
	f2 := New(b, []float64{1})
	blockOut := f1.Process(inputs)
	for i, x := range inputs {
		if s := f2.Step(x); math.Abs(s-blockOut[i]) > 1e-12 {
			t.Fatalf("block vs step mismatch at %d", i)
		}
	}
}
