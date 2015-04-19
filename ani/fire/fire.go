package fire

import (
	"math"
	"math/rand"
	"post6.net/goled/model"
	"post6.net/goled/projection"
	"post6.net/goled/util/clip"
	"post6.net/goled/vector"
)

const (
	All = iota
	Inside
	Outside
)

type Fire2d struct {
	width, height int
	colorTab      [2048][3]byte
	caMap         [2048]uint
	cells         []uint
	buf           [][3]byte
}

type Fire struct {
	fire *Fire2d

	old, new, mix [][3]byte

	buf             [][3]byte
	accum           [][3]float64
	max             []float64
	sinY            []float64
	vmap            []int
	phase, phaseMax int
	whichLeds       int
}

func uintMin(a, b uint) uint {

	if a < b {
		return a
	} else {
		return b
	}
}

func (a *Fire2d) initColorTab() {

	for i := range a.colorTab {
		a.colorTab[i] = [3]byte{0, 0, 0}
	}

	for i := 0; i < 256; i++ {

		r := clip.FloatToByte(float64(i) * 3)
		g := clip.FloatToByte(math.Pow(float64(i)/256, 1.5) * 1024)
		b := clip.FloatToByte(float64(i*i) / 256)
		a.colorTab[i] = [3]byte{r, g, b}
	}
}

func (a *Fire2d) initCaMap() {

	for i := range a.caMap {
		a.caMap[i] = uintMin(2047, uint(math.Pow(float64(i)/256, 1.25)/3.7*256.))
	}
}

func (a *Fire2d) initFire(w, h int) {

	a.width, a.height = w, h
	a.cells = make([]uint, w*(h+1))
	for i := range a.cells {
		a.cells[i] = 0
	}
}

func NewFire2d(width, height int) (a *Fire2d) {

	a = new(Fire2d)

	a.initColorTab()
	a.initCaMap()
	a.initFire(width, height)

	a.buf = make([][3]byte, width*height)

	return a
}

func (a *Fire2d) Next() [][3]byte {

	w, h := a.width, a.height

	for y := 0; y < h; y++ {

		sum := a.cells[(y+2)*w-1] + a.cells[(y+1)*w] + a.cells[(y+1)*w+1]
		a.cells[y*w] = a.caMap[uintMin(2047, sum)]

		for x := 1; x < w-1; x++ {

			sum = a.cells[(y+1)*w+x-1] + a.cells[(y+1)*w+x] + a.cells[(y+1)*w+x+1]
			a.cells[y*w+x] = a.caMap[uintMin(2047, sum)]

		}

		sum = a.cells[(y+2)*w-2] + a.cells[(y+2)*w-1] + a.cells[(y+1)*w]
		a.cells[y*w+w-1] = a.caMap[uintMin(2047, sum)]
	}

	for x := 0; x < w; x++ {
		a.cells[h*w+x] = uint(rand.Intn(2)) * 221
	}

	for i := range a.buf {
		a.buf[i] = a.colorTab[a.cells[i]]
	}

	return a.buf[:]
}

func newFire(leds []model.Led3D, whichLeds int) (a *Fire) {

	a = new(Fire)

	a.fire = NewFire2d(42, 20)

	points := make([]vector.Vector3, len(leds))

	for i, led := range leds {
		points[i] = led.Position
	}

	a.vmap, a.sinY = projection.Voronoi(a.fire.width, a.fire.height, points)

	a.buf = make([][3]byte, len(leds))
	a.accum = make([][3]float64, len(leds))
	a.max = make([]float64, len(leds))

	for y := 0; y < a.fire.height; y++ {
		for x := 0; x < a.fire.width; x++ {
			a.max[a.vmap[y*a.fire.width+x]] += a.sinY[y]
		}
	}

	a.mix = a.fire.Next()
	a.phase, a.phaseMax = 0, 4
	a.whichLeds = whichLeds

	return a
}

func NewFire(leds []model.Led3D) *Fire {
	return newFire(leds, All)
}

func NewInnerFire(leds []model.Led3D) *Fire {
	return newFire(leds, Inside)
}

func (a *Fire) Next() [][3]byte {

	if a.phase == 0 {

		a.mix = a.fire.Next()
	}

	a.phase++
	a.phase %= a.phaseMax

	sinY := a.sinY

	i := 0
	for y := 0; y < a.fire.height; y++ {
		density := sinY[y]
		for x := 0; x < a.fire.width; x++ {

			led := a.vmap[i]
			a.accum[led][0] += density * float64(a.mix[i][0])
			a.accum[led][1] += density * float64(a.mix[i][1])
			a.accum[led][2] += density * float64(a.mix[i][2])
			i++
		}
	}

	for i := range a.buf {

		if a.whichLeds == All || (i%5 == 4) == (a.whichLeds == Inside) {
			if a.max[i] == 0 {
				a.buf[i] = [3]byte{0, 0, 255}
			} else {
				a.buf[i][0] = clip.FloatToByte(a.accum[i][0] / a.max[i] / 8)
				a.buf[i][1] = clip.FloatToByte(a.accum[i][1] / a.max[i] / 8)
				a.buf[i][2] = clip.FloatToByte(a.accum[i][2] / a.max[i] / 8)
			}
			a.accum[i][0] = a.accum[i][0] * 7 / 8
			a.accum[i][1] = a.accum[i][1] * 7 / 8
			a.accum[i][2] = a.accum[i][2] * 7 / 8
		}
	}
	return a.buf[:]
}
