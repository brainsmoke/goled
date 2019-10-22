package polyhedron

import (
	"post6.net/goled/vector"
)

type RemapRoute []int
type RemapReorientRoute [][2]int

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func RemapSolid(solid Solid, first int, route RemapRoute) Solid {

	faces := make([]Face, min(len(solid.Faces), len(route)+1))
	points := []vector.Vector3{}

	face_mapping := make([]int, len(solid.Faces))
	point_mapping := make([]int, len(solid.Points))

	for i := range face_mapping {
		face_mapping[i] = -1
	}

	for i := range point_mapping {
		point_mapping[i] = -1
	}

	current := first

	for i := range faces {

		faces[i] = solid.Faces[current]
		if face_mapping[current] != -1 {
			panic("overlap in remapping");
		}
		face_mapping[current] = i

		if i < len(route) {
			current = solid.Faces[current].Neighbours[route[i]]
		}
	}

	for i := range faces {
		faces[i].Polygon = append([]int{}, faces[i].Polygon...)
		faces[i].Angles = append([]float64{}, faces[i].Angles...)

		faces[i].Neighbours = append([]int{}, faces[i].Neighbours...)
		for j, k := range faces[i].Neighbours {

			faces[i].Neighbours[j] = face_mapping[k]
		}
		for j, k := range faces[i].Polygon {

			if point_mapping[k] == -1 {
				point_mapping[k] = len(points)
				points = append(points, solid.Points[k])
			}

			faces[i].Polygon[j] = point_mapping[k]
		}
	}

	return Solid{ Points: points, Faces: faces }
}

