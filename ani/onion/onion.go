package onion

import (
	"math"
	"math/cmplx"
	"post6.net/goled/model"
)

var green = [3]byte{50, 186, 50}
var purple = [3]byte{135, 0, 135}
var white = [3]byte{255, 224, 98}

type Onion struct {
	phaseMax, phase, curColor, colorPhase int
	yRot                                  []int
	buf                                   [][3]byte
	colors                                [][3]byte
	colorMax                              []int
}

func NewOnion(leds []model.Led3D) *Onion {

	phaseMax := 600

	yRot := make([]int, len(leds))
	buf := make([][3]byte, len(leds))

	colors := [][3]byte{green, purple, white}
	colorMax := []int{480, 480, 180}

	for i, l := range leds {
		v := l.Position
		if l.Inside {
			yRot[i] = -1
			buf[i] = white
		} else {
			yRot[i] = int(float64(phaseMax) * (.5 + (cmplx.Phase(complex(v.Z, v.X)) / math.Pi / 2)))
		}
	}

	return &Onion{phaseMax: phaseMax, colorMax: colorMax, colors: colors, yRot: yRot, buf: buf}
}

func (o *Onion) Next() [][3]byte {

	for i := range o.buf {
		if o.yRot[i] == o.phase {
			o.buf[i] = o.colors[o.curColor]
		}
	}

	o.phase++
	o.phase %= o.phaseMax

	o.colorPhase++
	if o.colorPhase == o.colorMax[o.curColor] {
		o.colorPhase = 0
		o.curColor++
		o.curColor %= len(o.colors)
	}

	return o.buf[:]
}
