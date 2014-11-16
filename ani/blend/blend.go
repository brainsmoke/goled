package blend

import (
	"math"
	"post6.net/goled/ani"
)

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Blend(dst, f1, f2 [][3]byte, d float64) {

	size := min(len(dst), min(len(f1), len(f2)))

	for i := 0; i < size; i++ {

		r1, g1, b1 := float64(f1[i][0]), float64(f1[i][1]), float64(f1[i][2])
		r2, g2, b2 := float64(f2[i][0]), float64(f2[i][1]), float64(f2[i][2])
		dst[i] = [3]byte{
			byte(r1*(1-d) + r2*d),
			byte(g1*(1-d) + g2*d),
			byte(b1*(1-d) + b2*d),
		}
	}
}

type BlendAni struct {
	a1, a2          ani.Animation
	phase, maxPhase int
	buf             [][3]byte
}

func NewBlendAni(a1, a2 ani.Animation, period int) (a *BlendAni) {

	a = new(BlendAni)
	a.a1, a.a2 = a1, a2

	a.phase, a.maxPhase = 0, period

	return a
}

func (a *BlendAni) Next() [][3]byte {

	b1, b2 := a.a1.Next(), a.a2.Next()
	if a.buf == nil {
		a.buf = make([][3]byte, min(len(b1), len(b2)))
	}
	a.phase = (a.phase + 1) % a.maxPhase
	Blend(a.buf, b1, b2, 0.5+0.5*-math.Cos(math.Pi*2*float64(a.phase)/float64(a.maxPhase)))
	return a.buf
}
