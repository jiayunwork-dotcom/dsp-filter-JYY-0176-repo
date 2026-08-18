package response

import "math"

func Grid(n int) []float64 {
	if n < 2 {
		return nil
	}
	out := make([]float64, n)
	for i := range out {
		out[i] = 0.5 * float64(i) / float64(n-1)
	}
	return out
}

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Message }

const (
	ErrEmptyGrid   = "empty_frequency_grid"
	ErrBadFrequency = "bad_frequency"
)

func ValidateFreq(freq []float64) error {
	if len(freq) == 0 {
		return &Error{Code: ErrEmptyGrid, Message: "frequency grid is empty"}
	}
	for _, f := range freq {
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return &Error{Code: ErrBadFrequency, Message: "frequency grid contains NaN or infinity"}
		}
		if f < 0 || f > 0.5 {
			return &Error{Code: ErrBadFrequency, Message: "normalized frequency must be in [0, 0.5]"}
		}
	}
	return nil
}

func IsError(err error, code string) bool {
	e, ok := err.(*Error)
	return ok && e.Code == code
}
