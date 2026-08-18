package filter

type State struct {
	X []float64
	Y []float64
}

func (f *DirectFormI) State() State {
	return State{X: append([]float64(nil), f.x...), Y: append([]float64(nil), f.y...)}
}

func (f *DirectFormI) SetState(s State) {
	copy(f.x, s.X)
	copy(f.y, s.Y)
}

func (f *DirectFormI) Latency() int {
	if len(f.A) > 1 {
		return len(f.A) - 1
	}
	return len(f.B) - 1
}
