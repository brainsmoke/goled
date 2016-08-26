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
	Binary
	Striped
)

const (
	All = iota
	Inside
	Outside
)

type Gradient struct {
	phaseMax, phase int
	rot             []Vector3
	buf             [][3]byte
	wave            []byte
	useMap          []bool
}

func NewSpiral(leds []model.Led3D, gradientType, whichLeds int) *Gradient {

	rot := make([]Vector3, len(leds))
	buf := make([][3]byte, len(leds))
	useMap := make([]bool, len(leds))
	wave := make([]byte, 1024)

	for i, l := range leds {
		v := l.Position.Normalize()
		rot[i] = Vector3{
			(1.-v.X/2.+cmplx.Phase(complex(v.Y, v.Z)) / math.Pi / 2),
			(1.-v.Y/2.+cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2),
			(1.-v.Z/2.+cmplx.Phase(complex(v.X, v.Y)) / math.Pi / 2),
		}

		rot[i].X -= math.Floor(rot[i].X)
		rot[i].Y -= math.Floor(rot[i].Y)
		rot[i].Z -= math.Floor(rot[i].Z)
		useMap[i] = (whichLeds == All) || (leds[i].Inside == (whichLeds == Inside) )
	}

	for i := range wave {
		if gradientType == Smooth {
			wave[i] = byte(math.Pow((1+math.Sin(float64(i)/1024*2*math.Pi))/2, 2) * 255)
		} else if gradientType == Hard {
			wave[i] = byte(i / 4)
		} else if gradientType == Striped {
			if i & 0x1c0 == 0x1c0 {
				wave[i] = 255
			} else {
				wave[i] = 0
			}
		} else {
			wave[i] = byte( ((i>>9)&1) * 255 )
		}
	}

	return &Gradient{phaseMax: 4 * 7 * 8 * 9, phase: 0, rot: rot, buf: buf, wave: wave, useMap: useMap}
}

func NewGradient(leds []model.Led3D, gradientType, whichLeds int) *Gradient {

	rot := make([]Vector3, len(leds))
	buf := make([][3]byte, len(leds))
	useMap := make([]bool, len(leds))
	wave := make([]byte, 1024)

	for i, l := range leds {
		v := l.Position
		rot[i] = Vector3{
			cmplx.Phase(complex(v.Y, v.Z)) / math.Pi / 2,
			cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2,
			cmplx.Phase(complex(v.X, v.Y)) / math.Pi / 2,
		}
		useMap[i] = (whichLeds == All) || (leds[i].Inside == (whichLeds == Inside) )
	}

	for i := range wave {
		if gradientType == Smooth {
			wave[i] = byte(math.Pow((1+math.Sin(float64(i)/1024*2*math.Pi))/2, 2) * 255)
		} else if gradientType == Hard {
			wave[i] = byte(i / 4)
		} else if gradientType == Striped {
			if i & 0x1c0 == 0x1c0 {
				wave[i] = 255
			} else {
				wave[i] = 0
			}
		} else {
			wave[i] = byte( ((i>>9)&1) * 255 )
		}
	}

	return &Gradient{phaseMax: 4 * 7 * 8 * 9, phase: 0, rot: rot, buf: buf, wave: wave, useMap: useMap}
}

func (t *Gradient) Next() [][3]byte {

	phi := float64(t.phase) / float64(t.phaseMax)

	for i, rot := range t.rot {
		if t.useMap[i] {
			t.buf[i] = [3]byte{
				t.wave[int(1024*(rot.Z-phi*7+8))%1024],
				t.wave[int(1024*(rot.X-phi*8+9))%1024],
				t.wave[int(1024*(rot.Y-phi*9+10))%1024],
			}
		}
	}

	t.phase++
	t.phase %= t.phaseMax

	return t.buf[:]
}
