package design

import (
	"fmt"
	"math"
)

type Filter struct {
	Kind Kind      `json:"kind"`
	B    []float64 `json:"b"`
	A    []float64 `json:"a"`
}

func Design(spec *DesignSpec) (*Filter, error) {
	switch spec.Kind {
	case KindFIR:
		b, err := DesignFIR(spec)
		if err != nil {
			return nil, fmt.Errorf("design: %v", err)
		}
		return &Filter{Kind: KindFIR, B: b, A: []float64{1}}, nil
	case KindIIR:
		b, a, err := DesignIIR(spec)
		if err != nil {
			return nil, fmt.Errorf("design: %v", err)
		}
		return &Filter{Kind: KindIIR, B: b, A: a}, nil
	}
	return nil, &Error{Code: ErrBadKind, Message: "kind must be fir or iir"}
}

func (f *Filter) Order() int {
	if len(f.B) == 0 {
		return 0
	}
	return len(f.B) - 1
}

func (f *Filter) StableByDesign() bool {
	if len(f.A) == 0 {
		return true
	}
	// quick check on a[0] sign convention; full stability lives in zplane
	return math.Abs(f.A[0]-1) < 1e-9
}
