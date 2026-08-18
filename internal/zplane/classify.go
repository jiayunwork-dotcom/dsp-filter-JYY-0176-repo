package zplane

import "math"

type PoleClass struct {
	Inside  int
	Outside int
	OnUnit  int
	Real    int
	ConjugatePairs int
}

func Classify(poles []complex128, eps float64) PoleClass {
	var c PoleClass
	for _, p := range poles {
		r := absC(p)
		switch {
		case r > 1+eps:
			c.Outside++
		case r < 1-eps:
			c.Inside++
		default:
			c.OnUnit++
		}
		if math.Abs(imag(p)) < eps {
			c.Real++
		}
	}
	c.ConjugatePairs = countConjugatePairs(poles, eps)
	return c
}

func countConjugatePairs(poles []complex128, eps float64) int {
	used := make([]bool, len(poles))
	pairs := 0
	for i := range poles {
		if used[i] || math.Abs(imag(poles[i])) < eps {
			continue
		}
		for j := i + 1; j < len(poles); j++ {
			if used[j] {
				continue
			}
			if math.Abs(real(poles[i])-real(poles[j])) < eps &&
				math.Abs(imag(poles[i])+imag(poles[j])) < eps {
				used[i] = true
				used[j] = true
				pairs++
				break
			}
		}
	}
	return pairs
}

func Margin(poles []complex128) float64 {
	min := 1e9
	found := false
	for _, p := range poles {
		r := absC(p)
		if !found || r < min {
			min = r
			found = true
		}
	}
	if !found {
		return 1e9
	}
	return 1 - min
}
