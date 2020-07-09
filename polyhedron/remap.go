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

func rotateFace(orig Face, first int) Face {

	return Face {
		Normal : orig.Normal,
		Center : orig.Center,
		Polygon : append(append([]int{}, orig.Polygon[first:]...), orig.Polygon[:first]...),
		Neighbours : append(append([]int{}, orig.Neighbours[first:]...), orig.Neighbours[:first]...),
		Angles : append(append([]float64{}, orig.Angles[first:]...), orig.Angles[:first]...),
	}
}

func RemapSolid(solid Solid, first int, route RemapRoute) Solid {
	reorientRoute := make(RemapReorientRoute, len(route))
	for i := range route {
		reorientRoute[i][0] = route[i]
		reorientRoute[i][1] = -1
	}
	return RemapReorientSolid(solid, first, reorientRoute)
}

func RemapReorientSolid(solid Solid, first int, route RemapReorientRoute) Solid {

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
			egress, ingress := route[i][0], route[i][1]
			next := solid.Faces[current].Neighbours[egress]
			if ingress != -1 {
				for j,n := range solid.Faces[next].Neighbours {
					if n == current {
						firstEdge := (j - ingress + len(solid.Faces[next].Neighbours)) % len(solid.Faces[next].Neighbours)
						solid.Faces[next] = rotateFace(solid.Faces[next], firstEdge)
						break
					}
				}
			}
			current = next
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

