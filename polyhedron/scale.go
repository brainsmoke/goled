package polyhedron

import (
	"post6.net/goled/vector"
)

func Scale(faces []Face, factor float64) []Face {

	newFaces := make([]Face, len(faces))

	for i, f := range faces {
		newFaces[i] = f

		newFaces[i].Neighbours = append([]int(nil), f.Neighbours...)
		newFaces[i].Angles = append([]float64(nil), f.Angles...)

		newFaces[i].Center = f.Center.Mul(factor)
		newFaces[i].Polygon = make([]vector.Vector3, len(f.Polygon))

		for j := range f.Polygon {

			newFaces[i].Polygon[j] = f.Polygon[j].Mul(factor)
		}
	}

	return newFaces
}

