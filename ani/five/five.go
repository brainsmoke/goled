package five

import (
	"post6.net/goled/color"
)

type Five struct {
	colors [5]*color.ColorPlay
	buf    [][3]byte
}

func NewFive() *Five {

	t := new(Five)
	t.buf = make([][3]byte, 300)

	for i := range t.colors {
		t.colors[i] = color.NewColorPlay(64*(3+i), 1)
	}

	return t
}

func (t *Five) Next() [][3]byte {

	var colors [5][3]byte

	for i := range colors{
		colors[i] = t.colors[i].NextColor()
	}

	for i := range t.buf {
		
		t.buf[i] = colors[i%5]
	}

	return t.buf[:]
}
