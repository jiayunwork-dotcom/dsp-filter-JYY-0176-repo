package zplane

import "math"

type UnitPoint struct {
	Re float64 `json:"re"`
	Im float64 `json:"im"`
}

func UnitCirclePoints(n int) []UnitPoint {
	out := make([]UnitPoint, n)
	for i := 0; i < n; i++ {
		angle := 2 * math.Pi * float64(i) / float64(n)
		out[i] = UnitPoint{Re: math.Cos(angle), Im: math.Sin(angle)}
	}
	return out
}

type Polar struct {
	Radius float64 `json:"radius"`
	Angle  float64 `json:"angle"`
}

func ToPolar(p Complex) Polar {
	return Polar{Radius: math.Hypot(p.Re, p.Im), Angle: math.Atan2(p.Im, p.Re)}
}

func RootsNearUnitCircle(poles []Complex, eps float64) []Complex {
	out := make([]Complex, 0)
	for _, p := range poles {
		r := math.Hypot(p.Re, p.Im)
		if r > 1-eps && r < 1+eps {
			out = append(out, p)
		}
	}
	return out
}
