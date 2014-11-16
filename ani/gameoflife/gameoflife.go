package gameoflife

import (
	"math/rand"
	"post6.net/goled/color"
	"post6.net/goled/model"
)

type GameOfLife struct {
	colorPlay       *color.ColorPlay
	neighbours      [][4]int
	state, newState []int
	phase, phaseMax int
	iter, iterMax   int
	buf             [][3]byte
}

func NewGameOfLife() (g *GameOfLife) {

	g = new(GameOfLife)
	g.phaseMax = 32
	g.iterMax = 64

	g.neighbours = model.LedballLedNeighbours()

	g.state = make([]int, len(g.neighbours))
	g.buf = make([][3]byte, len(g.neighbours))
	g.newState = make([]int, len(g.neighbours))
	g.colorPlay = color.NewColorPlay(512, 3)

	return g
}

func (g *GameOfLife) init() {

	for i, _ := range g.state {
		g.state[i] = rand.Intn(2)
	}

}

func (g *GameOfLife) nextIteration() {

	if g.iter == 0 {
		g.init()
	}

	for i, _ := range g.state {
		sum := 0

		for _, n := range g.neighbours[i] {
			if g.state[n] == 1 {
				sum++
			}
		}

		if sum == 2 {
			g.newState[i] = 1
		} else if sum == 3 {
			g.newState[i] = g.state[i]
		} else {
			g.newState[i] = 0
		}
	}

	copy(g.state, g.newState)

	g.iter++
	g.iter %= g.iterMax
}

func (g *GameOfLife) Next() [][3]byte {

	color := g.colorPlay.NextColor()

	for i, x := range g.state {

		g.buf[i][0] = byte(float64(g.buf[i][0])*63/64 + float64(x*int(color[0]))/64)
		g.buf[i][1] = byte(float64(g.buf[i][1])*63/64 + float64(x*int(color[1]))/64)
		g.buf[i][2] = byte(float64(g.buf[i][2])*63/64 + float64(x*int(color[2]))/64)
	}

	if g.phase == 0 {
		g.nextIteration()
	}

	g.phase++
	g.phase %= g.phaseMax

	return g.buf[:]
}
