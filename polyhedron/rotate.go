package polyhedron

import (
	"post6.net/goled/vector"
)

func (s *Solid) Rotate(eye, center, north vector.Vector3) {

	m := vector.RotationMatrix(eye, center, north)

	for i := range s.Points {
		s.Points[i] = m.Mul(s.Points[i])
	}

	for i := range s.Faces {
		s.Faces[i].Center = m.Mul(s.Faces[i].Center)
		s.Faces[i].Normal = m.Mul(s.Faces[i].Normal)
	}
}

