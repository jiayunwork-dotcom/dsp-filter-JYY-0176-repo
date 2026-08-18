package response

type Stats struct {
	BandEdge3DB      float64 `json:"band_edge_3db"`
	FoundEdge        bool    `json:"found_edge"`
	PassbandRipple   float64 `json:"passband_ripple_db"`
	StopbandAtten    float64 `json:"stopband_attenuation_db"`
	GroupDelayAvg    float64 `json:"group_delay_avg"`
	DCGainDB         float64 `json:"dc_gain_db"`
	MaxGainDB        float64 `json:"max_gain_db"`
}

func StatsOf(res *Result, passbandUpTo, stopbandFrom int) Stats {
	var st Stats
	st.DCGainDB = res.MagnitudeDB[0]
	maxV := res.MagnitudeDB[0]
	for _, v := range res.MagnitudeDB {
		if v > maxV {
			maxV = v
		}
	}
	st.MaxGainDB = maxV
	if edge, ok := BandEdge(res.Freq, res.MagnitudeDB, -3.01); ok {
		st.BandEdge3DB = edge
		st.FoundEdge = true
	}
	st.PassbandRipple = PassbandRipple(res.MagnitudeDB, passbandUpTo)
	st.StopbandAtten = StopbandAttenuation(res.MagnitudeDB, stopbandFrom)
	sum := 0.0
	count := 0
	for _, v := range res.GroupDelay {
		if v == v && v < 1e9 {
			sum += v
			count++
		}
	}
	if count > 0 {
		st.GroupDelayAvg = sum / float64(count)
	}
	return st
}
