package led

import (
	"errors"
	"math"
	"post6.net/goled/util/clip"
	"strings"
)

const (
	RGB = iota
	GRB
	RBG
	BRG
	GBR
	BGR
)

var ledString = []string{
	"RGB",
	"GRB",
	"RBG",
	"BRG",
	"GBR",
	"BGR",
}

var ledMap = [...][3]int{

	{0, 1, 2},
	{1, 0, 2},
	{0, 2, 1},
	{2, 0, 1},
	{1, 2, 0},
	{2, 1, 0},
}

type LedOrder int

func (o *LedOrder) Set(value string) error {
	for i, s := range ledString {
		if strings.ToUpper(value) == s {
			*o = LedOrder(i)
			return nil
		}
	}
	return errors.New("unknown ledorder value")
}

func (o *LedOrder) String() string {
	return ledString[int(*o)]
}

type LedStrip struct {
	gammaMap          [256]byte
	brightness, gamma float64
	ledOrder          [3]int
	buf               []byte
}

func (s *LedStrip) calcGammaMap() {

	maxGamma := math.Pow(255., s.gamma)
	for i := range s.gammaMap {

		s.gammaMap[i] = clip.FloatToByte((1 + 2*math.Pow(float64(i)*s.brightness, s.gamma)*255/maxGamma) / 2)
	}
}

func NewLedStrip(size, ledOrder LedOrder, gamma, brightness float64) (s *LedStrip) {

	if int(ledOrder) >= len(ledMap) {
		panic("bad ledOrder value")
	}

	s = new(LedStrip)
	s.buf = make([]byte, size*3)

	s.ledOrder = ledMap[int(ledOrder)]
	s.gamma = gamma
	s.brightness = brightness
	s.calcGammaMap()

	return
}

func (s *LedStrip) SetGamma(gamma float64) {

	s.gamma = gamma
	s.calcGammaMap()
}

func (s *LedStrip) SetBrightness(brightness float64) {

	s.brightness = brightness
	s.calcGammaMap()
}

func (s *LedStrip) Data() []byte {

	return s.buf
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (s *LedStrip) LoadFrame(pixels [][3]byte) []byte {

	n := min(len(pixels), len(s.buf)/3)

	r_i, g_i, b_i := s.ledOrder[0], s.ledOrder[1], s.ledOrder[2]

	for i := 0; i < n; i++ {

		s.buf[i*3+r_i] = s.gammaMap[pixels[i][0]]
		s.buf[i*3+g_i] = s.gammaMap[pixels[i][1]]
		s.buf[i*3+b_i] = s.gammaMap[pixels[i][2]]
	}

	return s.buf
}
