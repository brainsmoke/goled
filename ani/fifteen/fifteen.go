package fifteen

import (
	"math"
	"math/cmplx"
	"post6.net/goled/color"
	"post6.net/goled/model"
	"post6.net/goled/util/clip"
)

type Fifteen struct {
	colors [15]*color.ColorPlay
	buf    [][3]byte
}

func NewFifteen(leds []model.Led3D) *Fifteen {

	t := new(Fifteen)
	t.buf = make([][3]byte, len(leds))

	for i := range t.colors {
		t.colors[i] = color.NewColorPlay(128*(3+i), 1)
	}

	return t
}

func (t *Fifteen) Next() [][3]byte {

	var colors [15][3]byte

	for i := range colors {
		colors[i] = t.colors[i].NextColor()
	}

	for i := range t.buf {

		t.buf[i] = colors[i%15]
	}

	return t.buf[:]
}

type FifteenWave struct {
	phaseMax, phase int
	colors          [15]*color.ColorPlay
	rot             [][15]float64
	buf             [][3]byte
	wave            []float64
}

func NewFifteenWave(leds []model.Led3D) *FifteenWave {

	t := new(FifteenWave)
	t.buf = make([][3]byte, len(leds))

	for i := range t.colors {
		t.colors[i] = color.NewColorPlay(128*(3+i), 3)
	}

	t.rot = make([][15]float64, len(leds))
	t.wave = make([]float64, 1024)

	for i := range t.colors {
		t.colors[i] = color.NewColorPlay(128*(3+i), 1)
	}

	for i, l := range leds {
		v := l.Position
		t.rot[i] = [15]float64{
			cmplx.Phase(complex(v.Y, v.Z)) / math.Pi / 2,
			cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2,
			cmplx.Phase(complex(v.X, v.Y)) / math.Pi / 2,
			cmplx.Phase(complex(v.Z, v.Y)) / math.Pi / 2,
			cmplx.Phase(complex(v.Y, v.X)) / math.Pi / 2,
			cmplx.Phase(complex(v.X, v.Z)) / math.Pi / 2,
			cmplx.Phase(complex(v.Y, v.Z)) / math.Pi / 2,
			cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2,
			cmplx.Phase(complex(v.X, v.Y)) / math.Pi / 2,
			cmplx.Phase(complex(v.Z, v.Y)) / math.Pi / 2,
			cmplx.Phase(complex(v.Y, v.X)) / math.Pi / 2,
			cmplx.Phase(complex(v.X, v.Z)) / math.Pi / 2,
			cmplx.Phase(complex(v.Y, v.Z)) / math.Pi / 2,
			cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2,
			cmplx.Phase(complex(v.X, v.Y)) / math.Pi / 2,
		}
	}

	for i := range t.wave {
		t.wave[i] = (1+math.Sin(float64(i)/1024*2*math.Pi))/2
	}

	t.phaseMax = 3 * 5 * 7 * 8 * 9
	t.phase = 0

	return t
}

func (t *FifteenWave) Next() [][3]byte {

	var colors [15][3]byte
	phaseMul := [15]float64{ 3, 5, 7, 8, 9, 9, 3, 5, 7, 8, 8, 9, 3, 5, 7,  }

	for i := range colors {
		colors[i] = t.colors[i].NextColor()
	}

	phi := float64(t.phase) / float64(t.phaseMax)

	for i, rot := range t.rot {
		kind := i%15
		m := t.wave[int(1024*(rot[kind]-phi*phaseMul[kind]+10)) % 1024]
		t.buf[i] = [3]byte{
			clip.FloatToByte(float64(colors[kind][0]) * m),
			clip.FloatToByte(float64(colors[kind][1]) * m),
			clip.FloatToByte(float64(colors[kind][2]) * m),
		}
	}

	t.phase += 3
	t.phase %= t.phaseMax

	return t.buf[:]
}
