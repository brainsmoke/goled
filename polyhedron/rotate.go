package polyhedron

import (
	"post6.net/goled/vector"
)

func Rotate(faces []Face, eye, center, north vector.Vector3) []Face {

	newFaces := make([]Face, len(faces))
	m := vector.RotationMatrix(eye, center, north)

	for i, f := range faces {
		newFaces[i] = f

		newFaces[i].Neighbours = append([]int(nil), f.Neighbours...)
		newFaces[i].Angles = append([]float64(nil), f.Angles...)

		newFaces[i].Center = m.Mul(f.Center)
		newFaces[i].Normal = m.Mul(f.Normal)
		newFaces[i].Polygon = make([]vector.Vector3, len(f.Polygon))

		for j := range f.Polygon {

			newFaces[i].Polygon[j] = m.Mul(f.Polygon[j])
		}
	}

	return newFaces
}

