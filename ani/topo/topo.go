package topo

import (
	"math"
	"post6.net/goled/model"
	"post6.net/goled/util/clip"
)

type Topo struct {
	phaseMax, phase int

	groups [][]int
	count  []int
	wave   []int
	buf    [][3]byte
}

func NewTopo(model *model.Model3D) *Topo {

	t := new(Topo)

	t.phaseMax = 10000

	t.buf = make([][3]byte, len(model.Leds))
	t.groups = [][]int(nil)
	t.count = []int(nil)

	for _, v := range model.Groups {
		max := 0

		for _, num := range v {
			if num > max {
				max = num
			}
		}

		t.groups = append(t.groups, v)
		t.count = append(t.count, max+1)
	}

	t.wave = make([]int, t.phaseMax)

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
