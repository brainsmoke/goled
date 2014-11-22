package color

import (
	"math"
	"math/rand"
	"post6.net/goled/util/clip"
)

type ColorPlay struct {
	phaseMax, phase int
	colors          [][3]byte
	wave            []int
}

func NewColorPlay(phaseMax, nColors int) *ColorPlay {

	wave := make([]int, phaseMax)
	for i := range wave {
		wave[i] = int(127 - math.Cos(math.Pi*2*float64(i)/float64(phaseMax))*127)
	}
	colors := make([][3]byte, nColors)
	return &ColorPlay{phaseMax: phaseMax, phase: 0, wave: wave, colors: colors}
}

func (c *ColorPlay) NextColor() [3]byte {

	r, g, b := 0, 0, 0

	for i := range c.colors {
		p := (c.phase + i*c.phaseMax/len(c.colors)) % c.phaseMax
		if p == 0 {
			c.colors[i] = HSIToRGB(rand.Float64(), 1, .5+rand.Float64()/2)
			//			c.colors[i] = [3]byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))}
		}
		r += int(c.colors[i][0]) * c.wave[p] / 256
		g += int(c.colors[i][1]) * c.wave[p] / 256
		b += int(c.colors[i][2]) * c.wave[p] / 256
	}

	c.phase++
	c.phase %= c.phaseMax

	return [3]byte{clip.IntToByte(r), clip.IntToByte(g), clip.IntToByte(b)}
}
