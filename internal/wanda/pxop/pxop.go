package pxop

import (
	"gocv.io/x/gocv"
)

type Vecb []uint8

func VecbFrom4(v gocv.Scalar) Vecb {
	return Vecb{uint8(v.Val1), uint8(v.Val2), uint8(v.Val3)}
}

func GetVecbAt(m gocv.Mat, row int, col int) Vecb {
	ch := m.Channels()
	v := make(Vecb, ch)

	for c := 0; c < ch; c++ {
		v[c] = m.GetUCharAt(row, col*ch+c)
	}

	return v
}

func IsVecbInBetween(c, upper, lower Vecb) bool {
	for i, v := range c {
		if v < lower[i] || v > upper[i] {
			return false
		}
	}

	return true
}
