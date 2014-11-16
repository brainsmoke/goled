package topo

import (
	"math"
	"post6.net/goled/model"
	"post6.net/goled/util/clip"
)

type Topo struct {
	phaseMax, phase int

	groups [5][300]int
	count  [5]int
	wave   []int
	buf    [][3]byte
}

func NewTopo() *Topo {

	t := new(Topo)

	t.phaseMax = 10000

	t.buf = make([][3]byte, 300)

	leds := model.LedballRaw()
	faces := model.LedballFaces()

	icosaedron := make([]int, 60)
	dodecahedron := make([]int, 60)

	for i := 0; i < 60; i++ {
		icosaedron[i] = -1
		dodecahedron[i] = -1
	}

	icoCount, dodeCount := 0, 0
	for i := 0; i < 60; i++ {
		if icosaedron[i] == -1 {

			for j := i; icosaedron[j] == -1; j = faces[j].Neighbours[model.TopLeft] {
				icosaedron[j] = icoCount
			}
			icoCount++
		}
		if dodecahedron[i] == -1 {

			for j := i; dodecahedron[j] == -1; j = faces[j].Neighbours[model.BottomLeft] {
				dodecahedron[j] = dodeCount
			}
			dodeCount++
		}
	}

	rhombCount := icoCount + dodeCount
	for i, l := range leds {

		f := l.Face
		t.groups[0][i] = icosaedron[f]
		t.groups[1][i] = f
		t.groups[2][i] = dodecahedron[f]
		if i%5 == 0 {
			t.groups[3][i] = icosaedron[f]
		} else if i%5 == 2 {
			t.groups[3][i] = icoCount + dodecahedron[f]
		} else if i%5 == 1 && t.groups[3][i] == 0 {

			topRight := faces[f].Neighbours[model.TopRight]
			bottomRight := faces[f].Neighbours[model.BottomRight]
			opposite := faces[topRight].Neighbours[model.BottomLeft]

			t.groups[3][i] = rhombCount
			t.groups[3][5*topRight+3] = rhombCount
			t.groups[3][5*opposite+1] = rhombCount
			t.groups[3][5*bottomRight+3] = rhombCount

			rhombCount++
		}
		t.groups[4][i] = i % 5
	}

	t.wave = make([]int, t.phaseMax)

	t.count[0] = icoCount
	t.count[1] = len(faces)
	t.count[2] = dodeCount
	t.count[4] = 5
	t.count[3] = rhombCount

	part := t.phaseMax / len(t.count)
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

	for i := range t.count {
		p := (t.phase + i*t.phaseMax/len(t.count)) % t.phaseMax
		phase := float64(p*5) / float64(t.phaseMax)
		mul := float64(t.wave[p])
		for j := range t.buf {
			phi := phase + float64(t.groups[i][j]*(t.count[i]/2+1))/float64(t.count[i])
			rPhi, gPhi, bPhi := phi*2*math.Pi, (phi+1./3.)*2*math.Pi, (phi+2./3.)*2*math.Pi
			t.buf[j][0] = clip.FloatToByte(float64(t.buf[j][0]) + math.Sin(rPhi*2)*mul)
			t.buf[j][1] = clip.FloatToByte(float64(t.buf[j][1]) + math.Sin(gPhi*3)*mul)
			t.buf[j][2] = clip.FloatToByte(float64(t.buf[j][2]) + math.Sin(bPhi)*mul)
		}
	}

	t.phase++
	t.phase %= t.phaseMax

	return t.buf[:]
}
