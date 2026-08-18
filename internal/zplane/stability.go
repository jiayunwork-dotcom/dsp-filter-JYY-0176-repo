package zplane

func Stable(zp *ZPResult) bool {
	return zp.Stable
}

func MaxPoleRadius(poles []Complex) float64 {
	max := 0.0
	for _, p := range poles {
		r := p.Re*p.Re + p.Im*p.Im
		if r > max {
			max = r
		}
	}
	return sqrtFloat(max)
}

func sqrtFloat(v float64) float64 {
	if v <= 0 {
		return 0
	}
	x := v
	for i := 0; i < 24; i++ {
		x = (x + v/x) / 2
	}
	return x
}
