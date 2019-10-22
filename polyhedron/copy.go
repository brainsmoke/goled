package polyhedron

import (
	"post6.net/goled/vector"
)


func (s *Solid) Copy() Solid {

	faces := append([]Face{}, s.Faces...)
	points := append([]vector.Vector3{}, s.Points...)

	for i := range faces {
		faces[i].Polygon = append([]int{}, faces[i].Polygon...)
		faces[i].Angles = append([]float64{}, faces[i].Angles...)
		faces[i].Neighbours = append([]int{}, faces[i].Neighbours...)
	}

	return Solid{ Points: points, Faces: faces }
}

func Combine(s1, s2 Solid) Solid {

	faces := append(append([]Face{}, s1.Faces...), s2.Faces...)
	points := append(append([]vector.Vector3{}, s1.Points...), s2.Points...)

	s2off := len(s1.Points)
	for i := len(s1.Faces); i < len(faces); i++ {
		for j,v := range faces[i].Polygon {
			faces[i].Polygon[j] = v+s2off
		}
	}

	return Solid{ Points: points, Faces: faces }
}

