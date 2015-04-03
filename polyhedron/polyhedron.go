package polyhedron

import (
	"math"
	"math/cmplx"
	"post6.net/goled/vector"
//		"fmt"
)

type Face struct {
	Normal, Center vector.Vector3
	Polygon        []vector.Vector3
	Neighbours     []int
	Angles         []float64
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

	panic("meh, not a valid polyhedron")
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

func CatalanDualFace(points []vector.Vector3, index int) Face {

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

	f.Polygon = make([]vector.Vector3, len(polygons))
	for i, poly := range polygons {
		x, y, z := 0., 0., 0.
		for _, j := range poly {
			x, y, z = x+points[j].X, y+points[j].Y, z+points[j].Z
		}
		f.Polygon[i] = ScaleToPlaneOnNormal(vector.Vector3{x, y, z}, f.Center)
	}

	return f
}

func CatalanDualFaces(points []vector.Vector3) []Face {

	f := make([]Face, len(points))

	for i := range points {

		f[i] = CatalanDualFace(points, i)
	}

	return f
}

func smallestElement(list []int) int {

	res := list[0]
	for i := range list {
		if res > list[i] {
			res = list[i]
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

func ArchimedeanFaces(points []vector.Vector3) []Face {

	face_list := [][]int{}

	for i := range points {
		_, polygons := vertexInfo(points, i)
		for _, p := range polygons {
			if p[0] == smallestElement(p) {
				face_list = append(face_list, p)
			}
		}
	}

	f := make([]Face, len(face_list))

	for i := range face_list {

		f[i].Polygon = make([]vector.Vector3, len(face_list[i]))
		f[i].Neighbours = make([]int, len(face_list[i]))

		for j := range face_list[i] {
			f[i].Polygon[j] = points[face_list[i][j]].Normalize()
		}

		f[i].Center = vector.Average(f[i].Polygon)
		f[i].Normal = f[i].Center.Normalize()

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

	return f
}


