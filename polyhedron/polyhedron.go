package polyhedron

import (
	"math"
	"math/cmplx"
	"post6.net/goled/vector"
//		"fmt"
)

type Face struct {
	Normal, Center vector.Vector3
	Polygon        []int
	Neighbours     []int
	Angles         []float64
}

type Solid struct {
	Points []vector.Vector3
	Faces []Face
}

func neighbourDistance(points []vector.Vector3, index int) float64 {

	minD := -1.

	for i := range points {
		if i != index {
			d := points[index].Distance(points[i])
			if minD < 0 || d < minD {
				minD = d
			}
		}
	}

	return minD
}

func neighbours(points []vector.Vector3, index int) []int {

	list := []int{}
	d := neighbourDistance(points, index)

	for i, v := range points {
		if math.Abs(points[index].Distance(v)-d) < .001 {

			list = append(list, i)
		}
	}

	return list
}

func ccwNeighbours(points []vector.Vector3, index int) []int {

	list := neighbours(points, index)
	p := points[index]
	pNormal := points[index].Normalize()
	phases := make([]float64, len(list))
	vA := points[list[0]].Sub(p).Normalize()

	for i, pi := range list {
		q := points[pi]
		vB := q.Sub(p).Normalize()

		sinTheta := vA.CrossProduct(vB).ScalarProduct(pNormal)
		cosTheta := vA.ScalarProduct(vB)
		phases[i] = cmplx.Phase(complex(cosTheta, sinTheta))
	}

	for i := 0; i < len(list)-1; i++ {
		for j := i + 1; j < len(list); j++ {
			if phases[j] < phases[i] {
				list[j], list[i] = list[i], list[j]
				phases[j], phases[i] = phases[i], phases[j]
			}
		}
	}

	return list
}

func findNextPoint(points []vector.Vector3, p0, p1, p2 int) int {

	vA := points[p0].Sub(points[p1])
	vB := points[p1].Sub(points[p2])
	AxB := vA.CrossProduct(vB)

	for _, i := range neighbours(points, p2) {
		vC := points[p2].Sub(points[i])
		BxC := vB.CrossProduct(vC)
		if AxB.Distance(BxC) < 0.001 {
			return i
		}
	}

	panic("meh, not a valid platonic/archimedean")
}

func findRegularPolygon(points []vector.Vector3, last, first, second int) []int {

	a, b, c := last, first, second
	list := []int{b, c}

	for c != last {
		a, b, c = b, c, findNextPoint(points, a, b, c)
		list = append(list, c)
	}

	return list
}

func ScaleToPlaneOnNormal(point, normal vector.Vector3) vector.Vector3 {
	return point.Mul(normal.ScalarProduct(normal) /
		point.ScalarProduct(normal))
}

func vertexInfo(points []vector.Vector3, index int) (neighbours []int, faces [][]int) {

	list := ccwNeighbours(points, index)
	size := len(list)

	polygons := make([][]int, size)
	for i := range list {
		last, first, second := list[i], index, list[(i+size-1)%size]
		polygons[i] = findRegularPolygon(points, last, first, second)
	}

	first := 0

	/* sort list of polygons based on vertex figure notation (smallest first) */
	for rot := range polygons {
		for i := range polygons {
			if len(polygons[(i+rot)%size]) < len(polygons[(i+first)%size]) {
				first = rot
				break
			}
			if len(polygons[(i+rot)%size]) > len(polygons[(i+first)%size]) {
				break
			}
		}
	}

	return append(list[first:], list[:first]...),
	       append(polygons[first:], polygons[:first]...)
}


func faceList(points []vector.Vector3) [][]int {

	face_list := [][]int{}

	for i := range points {
		_, polygons := vertexInfo(points, i)
		for _, p := range polygons {
			if smallestElementIndex(p) == 0 {
				face_list = append(face_list, p)
			}
		}
	}
	return face_list
}

func intArrayEquals(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i,v := range(a) {
		if b[i] != v {
			return false
		}
	}

	return true
}


func findFaceIndex(face_list [][]int, face[]int) int {
	e0 := smallestElementIndex(face)
	normalized_face := append(face[e0:], face[:e0]...)
	for i := range(face_list) {
		if intArrayEquals(face_list[i], normalized_face) {
			return i
		}
	}

	panic("face not found")
	//return -1
}

func CatalanDualFace(points []vector.Vector3, face_list [][]int, index int) Face {

	neighbours, polygons := vertexInfo(points, index)

	var f Face

	f.Center = points[index].Normalize()
	f.Normal = f.Center
	f.Neighbours = neighbours

	f.Angles = make([]float64, len(f.Neighbours))
	for i, neighbour := range f.Neighbours {
		d := f.Normal.Distance(points[neighbour].Normalize())
		f.Angles[i] = math.Acos(1-d*d)
	}

	f.Polygon = make([]int, len(polygons))

	for i, poly := range polygons {
		f.Polygon[i] = findFaceIndex(face_list, poly)
	}

	return f
}

func CatalanDual(points []vector.Vector3) Solid {
	/* in: vertices of an archimedean solid , out: its catalan dual */

	var solid Solid

	face_list := faceList(points)
	solid.Points = make([]vector.Vector3, len(face_list)) /* points become faces, faces become points */

	for i := range(face_list) {
		v := vector.Vector3{0,0,0}
		for _, ix := range face_list[i] {
			v = v.Add(points[ix])
		}
		plane_normal := points[face_list[i][0]].Normalize()
		solid.Points[i] = ScaleToPlaneOnNormal(v, plane_normal) /* scale to inscribed sphere */
	}

	solid.Faces = make([]Face, len(points))

	for i := range points {

		solid.Faces[i] = CatalanDualFace(points, face_list, i)
	}

	return solid
}

func smallestElementIndex(list []int) int {

	res := 0
	for i := range list {
		if list[res] > list[i] {
			res = i
		}
	}
	return res
}

func findPolygonWithEdge(pList [][]int, a, b int) int {

	for i := range pList {
		for j := range pList[i] {

			if a == pList[i][j] && b == pList[i][(j+1)%len(pList[i])] {
				return i
			}
		}
	}

	panic("meh, not a valid polyhedron")
}


func Archimedean(points []vector.Vector3) Solid {

	var solid Solid
	solid.Points = make([]vector.Vector3, len(points))

	for i := range points {
		solid.Points[i] = points[i].Normalize()
	}

	face_list := [][]int{}

	for i := range points {
		_, polygons := vertexInfo(points, i)
		for _, p := range polygons {
			if smallestElementIndex(p) == 0 {
				face_list = append(face_list, p)
			}
		}
	}

	f := make([]Face, len(face_list))

	for i := range face_list {

		f[i].Polygon = append([]int{}, face_list[i]...)
		f[i].Neighbours = make([]int, len(face_list[i]))

		v := vector.Vector3{0,0,0}

		for _, ix := range f[i].Polygon {
			v = v.Add(solid.Points[ix])
		}

		f[i].Center = v.Mul(1/float64(len(f[i].Polygon)))
		f[i].Normal = v.Normalize()

		for j := range face_list[i] {

			p, pNext := face_list[i][j], face_list[i][(j+1)%len(face_list[i])]

			f[i].Neighbours[j] = findPolygonWithEdge(face_list, pNext, p)
		}
	}

	for i := range f {
		f[i].Angles = make([]float64, len(f[i].Neighbours))
		for j, n := range f[i].Neighbours {
			d := f[i].Normal.Distance(f[n].Normal)
			f[i].Angles[j] = math.Acos(1-d*d)
		}
	}

	solid.Faces = f

	return solid
}


