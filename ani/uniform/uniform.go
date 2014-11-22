package uniform

import (
	"post6.net/goled/color"
)

type Uniform struct {
	colorPlay *color.ColorPlay
	buf       [][3]byte
	inside    bool
}

func newUniform(inside bool) *Uniform {

	return &Uniform{
		buf:       make([][3]byte, 300),
		colorPlay: color.NewColorPlay(1024, 2),
		inside:    inside,
	}
}

func NewUniform() *Uniform {

	return newUniform(false)
}

func NewUniformInside() *Uniform {

	return newUniform(true)
}

func (u *Uniform) Next() [][3]byte {

	color := u.colorPlay.NextColor()
	for i := range u.buf {
		if !u.inside || i%5 == 4 {
			u.buf[i] = color
		} else {
			u.buf[i] = [3]byte{0, 0, 0}
		}
	}

	return u.buf[:]
}
