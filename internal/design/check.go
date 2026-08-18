package design

import "math"

func CheckCoeffs(b, a []float64) error {
	if len(b) == 0 || len(a) == 0 {
		return &Error{Code: ErrBadOrder, Message: "coefficients must be non-empty"}
	}
	if a[0] == 0 {
		return &Error{Code: ErrBadOrder, Message: "a[0] must be non-zero"}
	}
	for _, v := range b {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return &Error{Code: ErrBadOrder, Message: "b contains NaN or infinity"}
		}
	}
	for _, v := range a {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return &Error{Code: ErrBadOrder, Message: "a contains NaN or infinity"}
		}
	}
	return nil
}
