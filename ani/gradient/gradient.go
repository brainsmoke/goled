package gradient

import (
	"math"
	"math/cmplx"
	"post6.net/goled/model"
	. "post6.net/goled/vector"
)

const (
	Smooth = iota
	Hard
)

type Gradient struct {
	phaseMax, phase int
	rot             []Vector3
	buf             [][3]byte
	wave            []byte
}

func NewGradient(leds []model.Led3D, gradientType int) *Gradient {

	rot := make([]Vector3, len(leds))
	buf := make([][3]byte, len(leds))
	wave := make([]byte, 1024)

	for i, l := range leds {
		v := l.Position
		rot[i] = Vector3{
			cmplx.Phase(complex(v.Y, v.Z)) / math.Pi / 2,
			cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2,
			cmplx.Phase(complex(v.X, v.Y)) / math.Pi / 2,
		}
	}

	for i := range wave {
		if gradientType == Smooth {
			wave[i] = byte(math.Pow((1+math.Sin(float64(i)/1024*2*math.Pi))/2, 2) * 255)
		} else {
			wave[i] = byte(i / 4)
		}
	}

	return &Gradient{phaseMax: 4 * 7 * 8 * 9, phase: 0, rot: rot, buf: buf, wave: wave}
}

func (t *Gradient) Next() [][3]byte {

	phi := float64(t.phase) / float64(t.phaseMax)

	for i, rot := range t.rot {
		t.buf[i] = [3]byte{
			t.wave[int(1024*(rot.Z-phi*7+8))%1024],
			t.wave[int(1024*(rot.X-phi*8+9))%1024],
			t.wave[int(1024*(rot.Y-phi*9+10))%1024],
		}
	}

	t.phase++
	t.phase %= t.phaseMax

	return t.buf[:]
}
