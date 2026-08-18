package design

import "testing"

func TestDesignFIRHammingSymmetric(t *testing.T) {
	b, err := DesignFIR(&DesignSpec{Kind: KindFIR, Order: 30, Cutoff: 0.2, Window: "hamming"})
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 31 {
		t.Fatalf("order 30 must yield 31 taps, got %d", len(b))
	}
	if !SymmetricCheck(b, 1e-12) {
		t.Fatalf("type-I fir must be even symmetric")
	}
}

func TestDesignFIRHammingSymmetricHann(t *testing.T) {
	b, err := DesignFIR(&DesignSpec{Kind: KindFIR, Order: 40, Cutoff: 0.15, Window: "hann"})
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 41 {
		t.Fatalf("order 40 must yield 41 taps, got %d", len(b))
	}
	if !SymmetricCheck(b, 1e-12) {
		t.Fatalf("type-I fir must be even symmetric")
	}
}

func TestDesignFIRUnknownWindow(t *testing.T) {
	_, err := DesignFIR(&DesignSpec{Kind: KindFIR, Order: 30, Cutoff: 0.2, Window: "kaiser"})
	if !IsError(err, ErrUnknownWindow) {
		t.Fatalf("expected unknown window error, got %v", err)
	}
}

func TestDesignFIRCutoffOutOfRange(t *testing.T) {
	for _, fc := range []float64{0, -0.1, 0.5, 0.9} {
		_, err := DesignFIR(&DesignSpec{Kind: KindFIR, Order: 30, Cutoff: fc, Window: "hamming"})
		if !IsError(err, ErrCutoffRange) {
			t.Fatalf("cutoff %v must be rejected, got %v", fc, err)
		}
	}
}

func TestDesignFIRZeroOrder(t *testing.T) {
	_, err := DesignFIR(&DesignSpec{Kind: KindFIR, Order: 0, Cutoff: 0.2, Window: "hamming"})
	if !IsError(err, ErrBadOrder) {
		t.Fatalf("order 0 must be rejected, got %v", err)
	}
	_, err = DesignFIR(&DesignSpec{Kind: KindFIR, Order: 15, Cutoff: 0.2, Window: "hamming"})
	if !IsError(err, ErrBadOrder) {
		t.Fatalf("odd fir order must be rejected, got %v", err)
	}
}

func TestDesignIIRA0Normalized(t *testing.T) {
	for _, order := range []int{1, 2, 4, 8} {
		b, a, err := DesignIIR(&DesignSpec{Kind: KindIIR, Order: order, Cutoff: 0.2})
		if err != nil {
			t.Fatal(err)
		}
		if a[0] != 1 {
			t.Fatalf("a[0] must be 1, got %v", a[0])
		}
		if len(b) != len(a) {
			t.Fatalf("b and a lengths differ: %d vs %d", len(b), len(a))
		}
	}
}

func TestDesignIIROrderLimits(t *testing.T) {
	_, _, err := DesignIIR(&DesignSpec{Kind: KindIIR, Order: 0, Cutoff: 0.2})
	if !IsError(err, ErrBadOrder) {
		t.Fatalf("order 0 must be rejected, got %v", err)
	}
	_, _, err = DesignIIR(&DesignSpec{Kind: KindIIR, Order: 17, Cutoff: 0.2})
	if !IsError(err, ErrBadOrder) {
		t.Fatalf("order 17 must be rejected, got %v", err)
	}
}

func TestWindowHannAlias(t *testing.T) {
	a, ok1 := ParseWindow("hann")
	b, ok2 := ParseWindow("hanning")
	if !ok1 || !ok2 || a != b {
		t.Fatalf("hann and hanning must be the same window")
	}
}

func TestDesignIICutoffBelowMax(t *testing.T) {
	_, _, err := DesignIIR(designSpecIIR(2, 0.5))
	if !IsError(err, ErrCutoffRange) {
		t.Fatalf("cutoff 0.5 must be rejected, got %v", err)
	}
}

func designSpecIIR(order int, fc float64) *DesignSpec {
	return &DesignSpec{Kind: KindIIR, Order: order, Cutoff: fc}
}
