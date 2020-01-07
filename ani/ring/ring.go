package ring

import (
	"math"
	"math/cmplx"
	"post6.net/goled/vector"
	"post6.net/goled/color"
	"post6.net/goled/model"
	"post6.net/goled/util/clip"
)

type RingWave struct {
	phaseMax, phase int
	colors          [12]*color.ColorPlay
	rot             []float64
	kind            []int
	buf             [][3]byte
	wave            []float64
}

func NewRingWave(model *model.Model3D) *RingWave {

	leds := model.Leds

	t := new(RingWave)
	t.buf = make([][3]byte, len(leds))

	t.rot = make([]float64, len(leds))
	t.kind = make([]int, len(leds))
	t.wave = make([]float64, 1024)

	for i := range t.colors {
		t.colors[i] = color.NewColorPlay(128*(3+i), 3)
	}

	var matrix [12]vector.Matrix3x3
	var done [12]bool

	for i, l := range leds {
		eye := l.Normal
		center := vector.Vector3{0,0,0}
		north := l.Position
		kind := model.Groups["sides"][i]
		t.kind[i] = kind
		if !done[kind] {
			done[kind] = true
			matrix[kind] = vector.RotationMatrix(eye, center, north)
		}
	}

	for i, l := range leds {
		v := matrix[t.kind[i]].Mul(l.Position)
		t.rot[i] = cmplx.Phase(complex(v.X, v.Y)) / math.Pi / 2
	}

	for i := range t.wave {
		t.wave[i] = (1+math.Sin(float64(i)/1024*2*math.Pi))/2
	}

	t.phaseMax = 3 * 5 * 7 * 8 * 9
	t.phase = 0

	return t
}

func (t *RingWave) Next() [][3]byte {

	var colors [12][3]byte
	phaseMul := [12]float64{ 3, 5, 7, 8, 9, 5, 3, 5, 7, 8, 9, 5 }

	for i := range colors {
		colors[i] = t.colors[i].NextColor()
	}

	phi := float64(t.phase) / float64(t.phaseMax)

	for i, rot := range t.rot {
		kind := t.kind[i]
		m := t.wave[int(1024*(rot-phi*phaseMul[kind]+10)) % 1024]
		t.buf[i] = [3]byte{
			clip.FloatToByte(float64(colors[kind][0]) * m),
			clip.FloatToByte(float64(colors[kind][1]) * m),
			clip.FloatToByte(float64(colors[kind][2]) * m),
		}
	}

	t.phase += 5
	//t.phase += 1
	t.phase %= t.phaseMax

	return t.buf[:]
}
