package design

import "fmt"

type Kind string

const (
	KindFIR Kind = "fir"
	KindIIR Kind = "iir"
)

type DesignSpec struct {
	Kind   Kind
	Order  int
	Cutoff float64
	Window string
}

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string { return e.Message }

const (
	ErrUnknownWindow = "unknown_window"
	ErrCutoffRange   = "cutoff_out_of_range"
	ErrBadOrder      = "bad_order"
	ErrBadKind       = "bad_kind"
)

func IsError(err error, code string) bool {
	e, ok := err.(*Error)
	return ok && e.Code == code
}

const (
	MaxIIROrder = 16
	MinIIROrder = 1
)

func Validate(spec *DesignSpec) error {
	if spec.Kind != KindFIR && spec.Kind != KindIIR {
		return &Error{Code: ErrBadKind, Message: "kind must be fir or iir"}
	}
	if spec.Cutoff <= 0 || spec.Cutoff >= 0.5 {
		return &Error{Code: ErrCutoffRange, Message: fmt.Sprintf("normalized cutoff %v must be in (0, 0.5)", spec.Cutoff)}
	}
	if spec.Order <= 0 {
		return &Error{Code: ErrBadOrder, Message: "order must be positive"}
	}
	switch spec.Kind {
	case KindFIR:
		if spec.Order%2 != 0 {
			return &Error{Code: ErrBadOrder, Message: "fir order must be even for type-I linear phase"}
		}
		if _, ok := ParseWindow(spec.Window); !ok {
			return &Error{Code: ErrUnknownWindow, Message: "unknown window " + spec.Window}
		}
	case KindIIR:
		if spec.Order < MinIIROrder || spec.Order > MaxIIROrder {
			return &Error{Code: ErrBadOrder, Message: fmt.Sprintf("iir order must be in [%d, %d]", MinIIROrder, MaxIIROrder)}
		}
	}
	return nil
}
