package polyhedron

import (
	"post6.net/goled/vector"
)

type RemapRoute []int
type RemapReorientRoute [][2]int

func RemapFaces(faces []Face, first int, route RemapRoute) []Face {

	newFaces := make([]Face, len(faces))
	mapping := make([]int, len(faces))

	current := first

	for i := range newFaces {

		newFaces[i] = faces[current]
		mapping[current] = i
		current = faces[current].Neighbours[route[i]]
	}

	for i := range newFaces {
		newFaces[i].Polygon = append([]vector.Vector3(nil), newFaces[i].Polygon...)
		newFaces[i].Angles = append([]float64(nil), newFaces[i].Angles...)

		newFaces[i].Neighbours = append([]int(nil), newFaces[i].Neighbours...)
		for j, k := range newFaces[i].Neighbours {

			newFaces[i].Neighbours[j] = mapping[k]
		}
	}

	return newFaces
}

