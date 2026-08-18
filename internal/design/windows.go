package design

import "math"

type Window int

const (
	WindowRect Window = iota
	WindowHann
	WindowHamming
	WindowBlackman
	WindowBartlett
)

func ParseWindow(name string) (Window, bool) {
	switch name {
	case "rect", "rectangular":
		return WindowRect, true
	case "hann", "hanning", "Hann":
		return WindowHann, true
	case "hamming":
		return WindowHamming, true
	case "blackman":
		return WindowBlackman, true
	case "bartlett", "triangular":
		return WindowBartlett, true
	}
	return 0, false
}

func (w Window) String() string {
	switch w {
	case WindowRect:
		return "rect"
	case WindowHann:
		return "hann"
	case WindowHamming:
		return "hamming"
	case WindowBlackman:
		return "blackman"
	case WindowBartlett:
		return "bartlett"
	}
	return "?"
}

func (w Window) Sample(n, total int) float64 {
	den := float64(total) - 1
	switch w {
	case WindowRect:
		return 1
	case WindowHann:
		return 0.5 - 0.5*math.Cos(2*math.Pi*float64(n)/den)
	case WindowHamming:
		return 0.54 - 0.46*math.Cos(2*math.Pi*float64(n)/den)
	case WindowBlackman:
		return 0.42 - 0.5*math.Cos(2*math.Pi*float64(n)/den) +
			0.08*math.Cos(4*math.Pi*float64(n)/den)
	case WindowBartlett:
		mid := float64(total-1) / 2
		return 1 - math.Abs(float64(n)-mid)/mid
	}
	return 1
}
