package poly

import (
	"post6.net/goled/model"
	"post6.net/goled/polyhedron"
)

type FacePosition struct {

	X, Y float64
	Inside bool
}

func populateFaceLeds(s polyhedron.Solid, index int, ledPositions []FacePosition) []model.Led3D {

	f := s.Faces[index]
	top, center, normal := s.Points[f.Polygon[0]], f.Center, f.Normal
	leds := make([]model.Led3D, len(ledPositions))

	for i, p := range ledPositions {

		vZ := normal
		vY := top.Sub(center).Normalize()
		vX := vY.CrossProduct(vZ).Normalize()
		pos := center.Add(vX.Mul(p.X)).Add(vY.Mul(p.Y))
		normal := f.Normal
		if p.Inside {
			normal = normal.Mul(-1)
		}
		leds[i] = model.Led3D{pos, normal, index, p.Inside}
//print (p.X, p.Y, "\n")
//print (center.X, center.Y, center.Z, " ", pos.X, pos.Y, pos.Z, "\n")
	}

	return leds
}


func PopulateLeds(s polyhedron.Solid, ledPositions [][]FacePosition) []model.Led3D {

	leds := []model.Led3D(nil)

	for i := range s.Faces {
		leds = append(leds, populateFaceLeds(s, i, ledPositions[i])...)
	}

	return leds
}

