package polyhedron

import (
	"post6.net/goled/vector"
	"fmt"
)

type RemapRoute []int

type RemapStep struct {

	Source, Egress, Ingress int
}

type RemapReorientRoute []RemapStep

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
		reorientRoute[i] = RemapStep{ Source: -1, Egress: route[i], Ingress: -1 }
	}
	return RemapReorientSolid(solid, first, reorientRoute)
}

func RemapReorientSolid(solid Solid, first int, route RemapReorientRoute) Solid {

	faces := make([]Face, min(len(solid.Faces), len(route)+1))
	oldFaces := append([]Face(nil), solid.Faces...)
	points := []vector.Vector3{}

	face_inv_mapping := make([]int, len(solid.Faces))
	face_mapping := make([]int, len(solid.Faces))
	point_mapping := make([]int, len(solid.Points))

	for i := range face_mapping {
		face_mapping[i] = -1
	}
	for i := range face_inv_mapping {
		face_inv_mapping[i] = -1
	}
	for i := range point_mapping {
		point_mapping[i] = -1
	}

	current := first

	for i := range faces {

		faces[i] = oldFaces[current]
		if face_inv_mapping[current] != -1 {
			fmt.Print("%d, %d\n", face_inv_mapping[current], i)
			panic("overlap in remapping");
		}
		face_inv_mapping[current] = i
		face_mapping[i] = current

		if i < len(route) {
			if route[i].Source > i {
				panic("bad remap route")
			}
			if route[i].Source > -1 {
				current = face_mapping[route[i].Source]
			}
			next := oldFaces[current].Neighbours[route[i].Egress]
			if route[i].Ingress != -1 {
				for j,n := range oldFaces[next].Neighbours {
					if n == current {
						firstEdge := (j - route[i].Ingress + len(oldFaces[next].Neighbours)) % len(oldFaces[next].Neighbours)
						oldFaces[next] = rotateFace(oldFaces[next], firstEdge)
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

			faces[i].Neighbours[j] = face_inv_mapping[k]
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

