package snake

import (
	"math/rand"
	"post6.net/goled/color"
	"post6.net/goled/model"
)

type Snake struct {
	colorPlay       *color.ColorPlay
	neighbours      [][4]int
	state, newState []int
	phase, phaseMax int
	buf             [][3]byte
	snake           []int
	dir, len, pos   int
}

func NewSnake() (s *Snake) {

	s = new(Snake)
	s.phaseMax = 20
	s.dir = 0
	s.snake = make([]int, 9)
	for i := range s.snake {
		s.snake[i] = -1
	}
	s.snake[0] = 4

	s.neighbours = model.LedballLedNeighbours()

	s.state = make([]int, len(s.neighbours))
	s.buf = make([][3]byte, len(s.neighbours))
	s.newState = make([]int, len(s.neighbours))
	s.colorPlay = color.NewColorPlay(256, 3)

	return s
}

func (s *Snake) nextIteration() {

	var r int

	if s.state[s.neighbours[s.snake[0]][s.dir]] != 0 {
		r = rand.Intn(5)
	} else {
		r = rand.Intn(32)
	}
	if r < 2 {
		s.dir++
	} else if r < 4 {
		s.dir += 3
	}

	s.dir %= 4
	copy(s.snake[1:len(s.snake)], s.snake[0:len(s.snake)-1])

	if r == 4 {
		if s.snake[1]%5 == 4 {
			s.snake[0] = (s.snake[1] / 5) * 5
		} else {
			s.snake[0] = (s.snake[1]/5)*5 + 4
		}
	} else {
		s.snake[0] = s.neighbours[s.snake[1]][s.dir]
	}

	for i := range s.state {
		s.state[i] = 0
	}
	for _, i := range s.snake {
		if i != -1 {
			s.state[i] = 1
		}
	}
}

func (s *Snake) Next() [][3]byte {

	color := s.colorPlay.NextColor()

	for i, x := range s.state {

		s.buf[i][0] = byte(float64(s.buf[i][0])*63/64 + float64(x*int(color[0]))/64)
		s.buf[i][1] = byte(float64(s.buf[i][1])*63/64 + float64(x*int(color[1]))/64)
		s.buf[i][2] = byte(float64(s.buf[i][2])*63/64 + float64(x*int(color[2]))/64)
	}

	if s.phase == 0 {
		s.nextIteration()
	}

	s.phase++
	s.phase %= s.phaseMax

	return s.buf[:]
}
