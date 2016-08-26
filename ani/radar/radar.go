package radar

import (
	"math"
	"math/cmplx"
	"post6.net/goled/model"
	"post6.net/goled/vector"
)

type Radar struct {
	phaseMax, phase int
	rot             []vector.Vector3
	buf             [][3]byte
	wave            []byte
	inside          []bool
}

func NewRadar(leds []model.Led3D) *Radar {

	rot := make([]vector.Vector3, len(leds))
	buf := make([][3]byte, len(leds))
	inside := make([]bool, len(leds))
	wave := make([]byte, 1024)

	for i, l := range leds {
		v := l.Position
		rot[i] = vector.Vector3{
			X: cmplx.Phase(complex(v.Y, v.Z)) / math.Pi / 2,
			Y: cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2,
			Z: cmplx.Phase(complex(v.X, v.Y)) / math.Pi / 2,
		}
		inside[i] = leds[i].Inside
	}

	for i := range wave {
		wave[i] = byte(math.Pow((1+math.Sin(float64(i)/1024*2*math.Pi))/2, 2) * 255)
	}

	return &Radar{phaseMax: 20 * 7 * 8 * 9, phase: 0, rot: rot, buf: buf, wave: wave, inside: inside}
}

func (r *Radar) Next() [][3]byte {

	phi := float64(r.phase) / float64(r.phaseMax)

	for i, rot := range r.rot {

		if !r.inside[i] {
			r.buf[i] = [3]byte{
				r.wave[int(1024*(rot.Z-phi*7+8))%1024],
				r.wave[int(1024*(rot.X-phi*8+9))%1024],
				0,
			}
		} else {

			r.buf[i] = [3]byte{
				0,
				0,
				r.wave[int(1024*(phi*9+10))%1024],
			}
		}
	}

	r.phase++
	r.phase %= r.phaseMax

	return r.buf[:]
}
