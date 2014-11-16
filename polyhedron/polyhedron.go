package polyhedron

import (
	"math"
	"math/cmplx"
	. "post6.net/goled/vector"
	//	"fmt"
)

type Face struct {
	Normal, Center Vector3
	Polygon        []Vector3
	Neighbours     []int
}

func neighbours(points []Vector3, index int) []int {

	list := []int{}

	for i, v := range points {
		if math.Abs(points[index].Distance(v)-2) < .001 {

			list = append(list, i)
		}
	}

	return list
}

func ccwNeighbours(points []Vector3, index int) []int {

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

func findNextPoint(points []Vector3, p0, p1, p2 int) int {

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

func findRegularPolygon(points []Vector3, last, first, second int) []int {

	a, b, c := last, first, second
	list := []int{b, c}

	for c != last {
		a, b, c = b, c, findNextPoint(points, a, b, c)
		list = append(list, c)
	}

	return list
}

func ScaleToPlaneOnNormal(point, normal Vector3) Vector3 {
	return point.Mul(normal.ScalarProduct(normal) /
		point.ScalarProduct(normal))
}

func CatalanDualFace(points []Vector3, index int) Face {

	list := ccwNeighbours(points, index)
	size := len(list)

	polygons := make([][]int, size)
	for i := range list {
		last, first, second := list[(i+size-1)%size], index, list[i]
		polygons[i] = findRegularPolygon(points, last, first, second)
	}

	best := 0

	for rot := range list {
		for i := range list {
			if len(polygons[(i+rot)%size]) < len(polygons[(i+best)%size]) {
				best = rot
				break
			}
			if len(polygons[(i+rot)%size]) > len(polygons[(i+best)%size]) {
				break
			}
		}
	}

	var f Face

	f.Center = points[index].Normalize()
	f.Normal = f.Center
	f.Neighbours = append(list[best:], list[:best]...)
	polygons = append(polygons[best:], polygons[:best]...)

	f.Polygon = make([]Vector3, size)
	for i, poly := range polygons {
		x, y, z := 0., 0., 0.
		for _, j := range poly {
			x, y, z = x+points[j].X, y+points[j].Y, z+points[j].Z
		}
		f.Polygon[i] = ScaleToPlaneOnNormal(Vector3{x, y, z}, f.Center)
	}

	return f
}

func CatalanDualFaces(points []Vector3) []Face {

	f := make([]Face, len(points))

	for i := range points {

		f[i] = CatalanDualFace(points, i)
	}

	return f
}

var (
	phi  = (1 + math.Sqrt(5)) / 2
	phi2 = phi * phi
	phi3 = phi * phi * phi
)

func RhombicosidodecahedronPoints() []Vector3 {

	p := []Vector3{}

	for _, x := range []float64{-1, 1} {
		for _, y := range []float64{-1, 1} {
			for _, z := range []float64{-phi3, phi3} {
				p = append(p, Vector3{x, y, z})
				p = append(p, Vector3{z, x, y})
				p = append(p, Vector3{y, z, x})
			}
		}
	}

	for _, x := range []float64{-phi2, phi2} {
		for _, y := range []float64{-phi, phi} {
			for _, z := range []float64{-2 * phi, 2 * phi} {
				p = append(p, Vector3{x, y, z})
				p = append(p, Vector3{z, x, y})
				p = append(p, Vector3{y, z, x})
			}
		}
	}

	for _, x := range []float64{-2 - phi, 2 + phi} {
		for _, z := range []float64{-phi2, phi2} {
			p = append(p, Vector3{x, 0, z})
			p = append(p, Vector3{z, x, 0})
			p = append(p, Vector3{0, z, x})
		}
	}

	return p
}

func DeltoidalhexecontahedronFaces() []Face {

	return CatalanDualFaces(RhombicosidodecahedronPoints())
}
