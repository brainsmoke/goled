package wobble

import (
	"math"
	"math/cmplx"
	"post6.net/goled/model"
)

const (
	All = iota
	Inside
	Outside
)

type Wobble struct {
	phaseMax, phase int
	yRot            []float64
	yPos            []float64
	use             []bool
	buf             [][3]byte
	wave            []byte
}

func NewWobble(leds []model.Led3D, whichLeds int) *Wobble {

	yRot := make([]float64, len(leds))
	yPos := make([]float64, len(leds))
	use := make([]bool, len(leds))
	buf := make([][3]byte, len(leds))
	wave := make([]byte, 1024)

	for i, l := range leds {
		v := l.Position
		yRot[i] = cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2
		yPos[i] = v.Y
		use[i] = whichLeds == All || (l.Inside == (whichLeds == Inside))
	}

	for i := range wave {
		wave[i] = byte(math.Pow((1+math.Sin(float64(i)/1024*2*math.Pi))/2, 2) * 255)
	}

	return &Wobble{phaseMax: 512 * 3, phase: 0, yRot: yRot, yPos: yPos, use: use, buf: buf, wave: wave}
}

func (t *Wobble) Next() [][3]byte {

	phi := float64(t.phase) / float64(t.phaseMax)

	for i := range t.buf {
		if t.use[i] {
			t.buf[i] = [3]byte{
				t.wave[int(1024*(2+phi*4+t.yRot[i]))%1024],
				t.wave[int(1024*(2+phi+t.yRot[i]))%1024],
				byte(math.Pow((1+math.Cos((t.yPos[i]+phi*8)*2*math.Pi)*math.Sin(phi*2*math.Pi/3))/2, 2) * 255),
			}
		} else {
			t.buf[i] = [3]byte{0, 0, 0}
		}
	}

	t.phase++
	t.phase %= t.phaseMax * 3

	return t.buf[:]
}
