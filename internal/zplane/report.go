package zplane

type RootReport struct {
	Value   Complex `json:"value"`
	Radius  float64 `json:"radius"`
	Angle   float64 `json:"angle"`
	Outside bool    `json:"outside"`
}

func ReportRoots(poles []complex128) []RootReport {
	out := make([]RootReport, len(poles))
	for i, p := range poles {
		cp := Complex{Re: real(p), Im: imag(p)}
		pol := ToPolar(cp)
		out[i] = RootReport{
			Value:   cp,
			Radius:  pol.Radius,
			Angle:   pol.Angle,
			Outside: pol.Radius > 1+StabilityEps,
		}
	}
	return out
}

func ReportRootsComplex(poles []Complex) []RootReport {
	out := make([]RootReport, len(poles))
	for i, p := range poles {
		pol := ToPolar(p)
		out[i] = RootReport{Value: p, Radius: pol.Radius, Angle: pol.Angle, Outside: pol.Radius > 1+StabilityEps}
	}
	return out
}
