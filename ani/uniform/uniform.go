package uniform

import (
	"post6.net/goled/color"
	"post6.net/goled/model"
)

type Uniform struct {
	colorPlay *color.ColorPlay
	buf       [][3]byte
	use       []bool
}

func newUniform(leds []model.Led3D, insideOnly bool) *Uniform {

	use := make([]bool, len(leds))

	for i := range use {
		use[i] = leds[i].Inside || !insideOnly
	}

	return &Uniform{
		buf:       make([][3]byte, len(leds)),
		use:       use,
		colorPlay: color.NewColorPlay(1024, 2),
	}
}

func NewUniform(leds []model.Led3D) *Uniform {

	return newUniform(leds, false)
}

func NewUniformInside(leds []model.Led3D) *Uniform {

	return newUniform(leds, true)
}

func (u *Uniform) Next() [][3]byte {

	color := u.colorPlay.NextColor()
	for i := range u.buf {
		if u.use[i] {
			u.buf[i] = color
		} else {
			u.buf[i] = [3]byte{0, 0, 0}
		}
	}

	return u.buf[:]
}
