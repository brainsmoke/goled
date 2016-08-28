package led

import (
	"math"
//"fmt"
	"post6.net/goled/util/clip"
)

type LedStrip interface {

	Gamma() float64
	SetGamma(gamma float64)

	Brightness() float64
	SetBrightness(brightness float64)

	SetCutoff(lowCutoff, highCutoff uint)
	Cutoff() (uint, uint)

	LedSize() int

	MapRange(start, count, byteOffset int)

	LoadFrame(pixels [][3]byte, frameBuffer []byte)
}

type BaseStrip struct {
	gammaMap          [256]uint
	order             LedOrder
	depth             int
	brightness, gamma float64
	maxValue          uint
	mapping           []int
	lowCutoff, highCutoff uint
}

func (s *BaseStrip) calcGammaMap() {

	gammaFactor := float64(s.maxValue)/math.Pow(255., s.gamma)
	for i := range s.gammaMap {

		s.gammaMap[i] = clip.FloatToUintRange(math.Pow(float64(i)*s.brightness, s.gamma) * gammaFactor,
		                                      0, s.maxValue)
	}
//fmt.Printf("%#v\n", s.gammaMap);
}

func NewLedStrip(size int, order LedOrder, depth int, maxValue uint, gamma, brightness float64) *BaseStrip {

	if !order.Valid() {
		panic("bad ledOrder value")
	}

	if depth != 8 && depth != 16 {
		panic("bad depth value")
	}

	s := &BaseStrip{
		order: order,
		depth: depth,
		brightness: brightness,
		gamma: gamma,
		maxValue: maxValue,
	}

	s.calcGammaMap()

	s.mapping = make([]int, size)
	for i := range(s.mapping) {
		s.mapping[i] = -1
	}

	return s
}

func (s *BaseStrip) Gamma() float64 {
	return s.gamma
}

func (s *BaseStrip) SetGamma(gamma float64) {

	s.gamma = gamma
	s.calcGammaMap()
}

func (s *BaseStrip) Brightness() float64 {
	return s.brightness
}

func (s *BaseStrip) SetBrightness(brightness float64) {

	s.brightness = brightness
	s.calcGammaMap()
}


func (s *BaseStrip) SetCutoff(lowCutoff, highCutoff uint) {

	s.lowCutoff = lowCutoff
	s.highCutoff = highCutoff
}

func (s *BaseStrip) Cutoff() (uint, uint) {

	return s.lowCutoff, s.highCutoff
}

func (s *BaseStrip) LedSize() int {
	return len(s.order.IndexMap())*s.depth/8
}

func (s *BaseStrip) MapRange(start, count, byteOffset int) {

	ledSize := s.LedSize()

	for i := 0; i < count && start+i < len(s.mapping); i++ {
		s.mapping[start+i] = byteOffset+i*ledSize
	}
}


func min3(a, b, c uint) uint {
	x := a

	if b < x {
		x = b
	}
	if c < x {
		x = c
	}
	return x
}

func max3(a, b, c uint) uint {
	x := a

	if b > x {
		x = b
	}
	if c > x {
		x = c
	}
	return x
}

func (s *BaseStrip) RGBToRGB(pixel [3]byte) (r, g, b uint) {

	gMap := s.gammaMap
	r, g, b = gMap[pixel[0]], gMap[pixel[1]], gMap[pixel[2]]
	max := max3(r, g, b)
	if max < s.highCutoff {
		r, g, b = 0, 0, 0
	}
	if r < s.lowCutoff {
		r = 0
	}
	if g < s.lowCutoff {
		g = 0
	}
	if b < s.lowCutoff {
		b = 0
	}
	return r, g, b
}

func (s *BaseStrip) RGBToRGBW(pixel [3]byte) (r, g, b, w uint) {

	gMap := s.gammaMap
	r, g, b = gMap[pixel[0]], gMap[pixel[1]], gMap[pixel[2]]
	w = min3(r, g, b)
	max := max3(r, g, b)
	r, g, b = r-w, g-w, b-w
	if max < s.highCutoff {
		r, g, b, w = 0, 0, 0, 0
	}
	if r < s.lowCutoff {
		r = 0
	}
	if g < s.lowCutoff {
		g = 0
	}
	if b < s.lowCutoff {
		b = 0
	}
	if w < s.lowCutoff {
		w = 0
	}
	return r, g, b, w
}

func (s *BaseStrip) LoadFrame(pixels [][3]byte, frameBuffer []byte) {

	colorIndex := s.order.IndexMap()

	if s.depth == 8 {

		if len(colorIndex) == 3 {
			r_i, g_i, b_i := colorIndex[0], colorIndex[1], colorIndex[2]

			for i,bufOffset := range(s.mapping) {
				if s.mapping[i] < 0 {
					continue
				}
				r, g, b := s.RGBToRGB(pixels[i])
				frameBuffer[bufOffset+r_i] = byte(r)
				frameBuffer[bufOffset+g_i] = byte(g)
				frameBuffer[bufOffset+b_i] = byte(b)
			}
		} else if len(colorIndex) == 4 {
			r_i, g_i, b_i, w_i := colorIndex[0], colorIndex[1], colorIndex[2], colorIndex[3]

			for i,bufOffset := range(s.mapping) {
				if s.mapping[i] < 0 {
					continue
				}
				r, g, b, w := s.RGBToRGBW(pixels[i])
				frameBuffer[bufOffset+r_i] = byte(r)
				frameBuffer[bufOffset+g_i] = byte(g)
				frameBuffer[bufOffset+b_i] = byte(b)
				frameBuffer[bufOffset+w_i] = byte(w)
			}
		} else {
			panic("bad number of colors")
		}
	} else if s.depth == 16 {

		if len(colorIndex) == 3 {
			r_i, g_i, b_i := colorIndex[0]*2, colorIndex[1]*2, colorIndex[2]*2

			for i,bufOffset := range(s.mapping) {
				if s.mapping[i] < 0 {
					continue
				}
				r, g, b := s.RGBToRGB(pixels[i])
				frameBuffer[bufOffset+r_i] = byte(r)
				frameBuffer[bufOffset+r_i+1] = byte(r>>8)
				frameBuffer[bufOffset+g_i] = byte(g)
				frameBuffer[bufOffset+g_i+1] = byte(g>>8)
				frameBuffer[bufOffset+b_i] = byte(b)
				frameBuffer[bufOffset+b_i+1] = byte(b>>8)
			}
		} else if len(colorIndex) == 4 {
			r_i, g_i, b_i, w_i := colorIndex[0]*2, colorIndex[1]*2, colorIndex[2]*2, colorIndex[3]*2

			for i,bufOffset := range(s.mapping) {
				if s.mapping[i] < 0 {
					continue
				}
				r, g, b, w := s.RGBToRGBW(pixels[i])
				frameBuffer[bufOffset+r_i] = byte(r)
				frameBuffer[bufOffset+r_i+1] = byte(r>>8)
				frameBuffer[bufOffset+g_i] = byte(g)
				frameBuffer[bufOffset+g_i+1] = byte(g>>8)
				frameBuffer[bufOffset+b_i] = byte(b)
				frameBuffer[bufOffset+b_i+1] = byte(b>>8)
				frameBuffer[bufOffset+w_i] = byte(w)
				frameBuffer[bufOffset+w_i+1] = byte(w>>8)
			}
		} else {
			panic("bad number of colors")
		}

	} else {
		panic("bad number of colors")
	}
}
