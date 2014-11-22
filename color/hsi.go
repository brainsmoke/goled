package color

import (
	"math"
	"post6.net/goled/util/clip"
)

var hsiTab [256]float64

func init() {
	for i := range hsiTab {
		hsiTab[i] = math.Cos(math.Pi*2*float64(i)/768) / math.Cos(math.Pi*2*float64(128-i)/768)
	}
}

/*
 * source formulas:
 * http://blog.saikoled.com/post/43693602826/why-every-led-light-should-be-using-hsi-colorspace
 */

func HSIToRGB(h, s, i float64) [3]byte {

	hInt := uint(768*h) % 786
	sClip := clip.FloatBetween(s, 0, 1)
	iClip := clip.FloatBetween(i, 0, 1)

	ix := hInt & 255
	a := clip.FloatToByte(255 * iClip * (1 + sClip*hsiTab[ix]) / 3)
	b := clip.FloatToByte(255 * iClip * (1 + sClip*(1-hsiTab[ix])) / 3)
	c := clip.FloatToByte(255 * iClip * (1 - sClip))
	if hInt < 256 {
		return [3]byte{a, b, c}
	} else if hInt < 512 {
		return [3]byte{b, c, a}
	} else {
		return [3]byte{c, a, b}
	}
}
