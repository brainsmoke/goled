package polyhedron

import (
//	"post6.net/goled/vector"
)

func (s *Solid) Scale(factor float64) {

	for i := range s.Points {
		s.Points[i] = s.Points[i].Mul(factor)
	}

	for i := range s.Faces {
		s.Faces[i].Center = s.Faces[i].Center.Mul(factor)
	}
}

