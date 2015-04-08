package model

import (
	"math"
)

type Model3D struct {

	Leds []Led3D
	Neighbours [][]int
	Groups map[string][]int
}

func copyNeighbours(list [][]int) [][]int {

	newList := make([][]int, len(list))

	for i, v := range list {
		newList[i] = append([]int(nil), v...)
	}

	return newList
}

func copyGroups(groups map[string][]int) map[string][]int {

	newGroups := make(map[string][]int)

	for k, v := range groups {
		newGroups[k] = append([]int(nil), v...)
	}

	return newGroups
}

func (m *Model3D) Copy() *Model3D {

	c := new(Model3D)
	c.Leds = append([]Led3D(nil), m.Leds...)
	c.Neighbours = copyNeighbours(m.Neighbours)
	c.Groups = copyGroups(m.Groups)

	return c
}

func (m *Model3D) Scale(factor float64) *Model3D {

	scaledModel := m.Copy()

	for i := range scaledModel.Leds {
		scaledModel.Leds[i].Position = scaledModel.Leds[i].Position.Mul(1. / factor)
	}

	return scaledModel
}

func (m *Model3D) UnitScale() *Model3D {

	max := 0.
	for _, p := range m.Leds {
		max = math.Max(max, p.Position.Magnitude())
	}

	return m.Scale(1./max)
}

func (m *Model3D) Smooth() *Model3D {

	smoothModel := m.Copy()

	for i := range smoothModel.Leds {

		p := smoothModel.Leds[i].Position
		smoothModel.Leds[i].Position = p.Normalize()

		if smoothModel.Leds[i].Inside {
			smoothModel.Leds[i].Normal = p.Mul(-1)
		} else {
			smoothModel.Leds[i].Normal = p
		}
	}

	return smoothModel
}
