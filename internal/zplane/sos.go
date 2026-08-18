package zplane

type Section struct {
	B0, B1, B2 float64
	A0, A1, A2 float64
}

func SOS(zp *ZPResult) []Section {
	zeros := toComplex128(zp.Zeros)
	poles := toComplex128(zp.Poles)
	ns := maxInt(len(zeros), len(poles))
	sections := make([]Section, 0, (ns+1)/2)
	zUsed := make([]bool, len(zeros))
	pUsed := make([]bool, len(poles))
	for len(sections) < (ns+1)/2 {
		var zPair []complex128
		var pPair []complex128
		zPair = nextPair(zeros, zUsed, zPair)
		pPair = nextPair(poles, pUsed, pPair)
		if len(zPair) == 1 {
			zPair = append(zPair, complex(0, 0))
		}
		if len(pPair) == 1 {
			pPair = append(pPair, complex(0, 0))
		}
		if len(zPair) == 0 && len(pPair) == 0 {
			break
		}
		sections = append(sections, Section{
			B0: 1, B1: -(real(zPair[0]) + real(zPair[1])), B2: real(zPair[0]*zPair[1]),
			A0: 1, A1: -(real(pPair[0]) + real(pPair[1])), A2: real(pPair[0]*pPair[1]),
		})
	}
	return sections
}

func nextPair(roots []complex128, used []bool, pair []complex128) []complex128 {
	for i, r := range roots {
		if used[i] {
			continue
		}
		used[i] = true
		return append(pair, r)
	}
	return pair
}

func toComplex128(cs []Complex) []complex128 {
	out := make([]complex128, len(cs))
	for i, c := range cs {
		out[i] = complex(c.Re, c.Im)
	}
	return out
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
