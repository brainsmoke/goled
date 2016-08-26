package topo

import (
	//"math"
	"post6.net/goled/model"
	"post6.net/goled/color"
	"post6.net/goled/util/clip"
)

type Topo struct {
	phaseMax, phase int

	groups [][]int
	colors [][]*color.ColorPlay
	wave   []int
	buf    [][3]byte
}

func NewTopo(model *model.Model3D) *Topo {

	t := new(Topo)

	t.phaseMax = 10000

	t.buf = make([][3]byte, len(model.Leds))
	t.groups = [][]int(nil)
	t.colors = [][]*color.ColorPlay(nil)

	for _, v := range model.Groups {
		max := 0

		for _, num := range v {
			if num > max {
				max = num
			}
		}

		t.groups = append(t.groups, v)

		c := make([]*color.ColorPlay, max+1)
		for i := range c {

			c[i] = color.NewColorPlay(128+i*2, 1)
		}

		t.colors = append(t.colors, c)
	}

	t.wave = make([]int, t.phaseMax)

	part := t.phaseMax / len(t.groups)
	slope := part / 5
	for i := 0; i < slope; i++ {
		t.wave[i] = i * 255 / slope
		t.wave[part+i] = (slope - i - 1) * 255 / slope
	}
	for i := slope; i < part; i++ {
		t.wave[i] = 255
	}

	return t
}

func (t *Topo) Next() [][3]byte {

	for j := range t.buf {
		t.buf[j] = [3]byte{0, 0, 0}
	}

	for i := range t.colors {

		for j := range t.colors[i] {
			t.colors[i][j].NextColor()
		}
	}

	for i := range t.buf {
		var r, g, b int
		for j := range t.groups {
			n := t.groups[j][i]
			var color [3]byte = t.colors[j][n].Color()

			mul := t.wave[ (t.phase + j*t.phaseMax/len(t.groups)) % t.phaseMax ]

			r += int(color[0])*mul
			g += int(color[1])*mul
			b += int(color[2])*mul
		}
		t.buf[i][0] = clip.IntToByte(r/255)
		t.buf[i][1] = clip.IntToByte(g/255)
		t.buf[i][2] = clip.IntToByte(b/255)
	}

	t.phase++
	t.phase %= t.phaseMax

	return t.buf[:]
}
