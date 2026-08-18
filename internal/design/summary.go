package design

type Summary struct {
	Kind       string  `json:"kind"`
	Order      int     `json:"order"`
	Taps       int     `json:"taps"`
	DCOGain    float64 `json:"dc_gain"`
	GroupDelay float64 `json:"group_delay_samples"`
	Stable     bool    `json:"stable"`
	Window     string  `json:"window,omitempty"`
	Cutoff     float64 `json:"cutoff"`
}

func Summarize(f *Filter, cutoff float64, window string) Summary {
	s := Summary{
		Kind:    string(f.Kind),
		Order:   f.Order(),
		Taps:    len(f.B),
		DCOGain: dcGain(f.B, f.A),
		Cutoff:  cutoff,
		Window:  window,
	}
	if f.Kind == KindFIR {
		s.GroupDelay = float64(f.Order()) / 2
		s.Stable = true
	} else {
		s.Stable = JuryStable(f.A)
	}
	return s
}
