package shadowplay

import (
	"math"
	"math/rand"
	"post6.net/goled/util/clip"
)

type ShadowPlay struct {
	phaseMax, phase int
	colors          [][3]byte
	indices         []int
	wave            []int
	buf             [][3]byte
}

func NewShadowPlay(phaseMax, nColors int) *ShadowPlay {

	buf := make([][3]byte, 300)

	wave := make([]int, phaseMax)
	colors := make([][3]byte, nColors)
	indices := make([]int, nColors)

	for i := range wave {
		wave[i] = int(127 - math.Cos(math.Pi*2*float64(i)/float64(phaseMax))*127)
	}

	return &ShadowPlay{phaseMax: phaseMax, phase: 0, buf: buf, wave: wave, colors: colors, indices: indices}
}

func (t *ShadowPlay) Next() [][3]byte {

	for i := range t.buf {
		t.buf[i] = [3]byte{0, 0, 0}
	}

	for i := range t.colors {
		p := (t.phase + i*t.phaseMax/len(t.colors)) % t.phaseMax
		if p == 0 {
			t.indices[i] = rand.Intn(60)*5 + 4
			t.colors[i] = [3]byte{byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256))}
		}
		t.buf[t.indices[i]][0] = clip.IntToByte(int(t.buf[t.indices[i]][0]) + int(t.colors[i][0])*t.wave[p]/256)
		t.buf[t.indices[i]][1] = clip.IntToByte(int(t.buf[t.indices[i]][1]) + int(t.colors[i][1])*t.wave[p]/256)
		t.buf[t.indices[i]][2] = clip.IntToByte(int(t.buf[t.indices[i]][2]) + int(t.colors[i][2])*t.wave[p]/256)
	}

	t.phase++
	t.phase %= t.phaseMax

	return t.buf[:]
}
