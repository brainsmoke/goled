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

func evenPermutations(x, y, z float64) []Vector3 {

	return []Vector3{

		{x, y, z},
		{z, x, y},
		{y, z, x},
	}
}

func allPermutations(x, y, z float64) []Vector3 {

	return []Vector3{

		{x, y, z},
		{x, z, y},
		{y, x, z},
		{y, z, x},
		{z, x, y},
		{z, y, x},
	}
}

func signPermutations(x, y, z float64) []Vector3 {

	p := []Vector3{}

	for _, sx := range []float64{ -x, x } {
		for _, sy := range []float64{ -y, y } {
			for _, sz := range []float64{ -z, z } {
				p = append(p, Vector3{sx, sy, sz})
				if z == -z {
					break
				}
			}
			if y == -y {
				break
			}
		}
		if x == -x {
			break
		}
	}

	return p
}

func signEvenPermutations(x, y, z float64) []Vector3 {

	p := []Vector3{}

	for _, v := range signPermutations(x, y, z) {
		p = append(p, evenPermutations(v.X, v.Y, v.Z)...)
	}

	return p
}

func allSignPermutations(x, y, z float64) []Vector3 {

	p := []Vector3{}

	for _, v := range signPermutations(x, y, z) {
		p = append(p, allPermutations(v.X, v.Y, v.Z)...)
	}

	return p
}

/* point lists */

func TetrahedronPoints() []Vector3 {

	return []Vector3 {
		{ -1, 0, -1/math.Sqrt(2) },
		{  1, 0, -1/math.Sqrt(2) },
		{ -1, 0,  1/math.Sqrt(2) },
		{  1, 0,  1/math.Sqrt(2) },
	}
}

func CubePoints() []Vector3 {

	return signPermutations(1, 1, 1)
}

func IcosahedronPoints() []Vector3 {

	return signEvenPermutations(0, 1, phi)
}

func TruncatedTetrahedronPoints() []Vector3 {

	p := []Vector3{}

	p = append(p, evenPermutations( 3,  1,  1)...)
	p = append(p, evenPermutations( 3, -1, -1)...)
	p = append(p, evenPermutations(-3,  1, -1)...)
	p = append(p, evenPermutations(-3, -1,  1)...)

	return p
}

func TruncatedCubePoints() []Vector3 {

	return signEvenPermutations( math.Sqrt(2)-1, 1, 1 )
}

func TruncatedCuboctahedronPoints() []Vector3 {

	return allSignPermutations(1, 1+math.Sqrt(2), 1+2*math.Sqrt(2))
}

func TruncatedOctahedronPoints() []Vector3 {

	return allSignPermutations(0, 1, 2)
}

func TruncatedDodecahedronPoints() []Vector3 {

	p := []Vector3{}

	p = append(p, signEvenPermutations(0, 1/phi, 2+phi)...)
	p = append(p, signEvenPermutations(1/phi, phi, 2*phi)...)
	p = append(p, signEvenPermutations(phi, 2, phi2)...)

	return p
}

func TruncatedIcosidodecahedronPoints() []Vector3 {

	p := []Vector3{}

	p = append(p, signEvenPermutations(1/phi, 1/phi, 3+phi)...)
	p = append(p, signEvenPermutations(2/phi, phi, 1+2*phi)...)
	p = append(p, signEvenPermutations(1/phi, phi2, -1+3*phi)...)
	p = append(p, signEvenPermutations(-1+2*phi, 2, 2+phi)...)
	p = append(p, signEvenPermutations(phi, 3, 2*phi)...)

	return p
}

func TruncatedIcosahedronPoints() []Vector3 {

	p := []Vector3{}

	p = append(p, signEvenPermutations(2, 1+2*phi, phi)...)
	p = append(p, signEvenPermutations(1, 2+phi, 2*phi)...)
	p = append(p, signEvenPermutations(1, 3*phi, 0)...)

	return p
}

func CuboctahedronPoints() []Vector3 {

	return signEvenPermutations( 1, 0, 0 )
}

func IcosidodecahedronPoints() []Vector3 {

	p := []Vector3{}

	p = append(p, signEvenPermutations(phi, 0, 0)...)
	p = append(p, signEvenPermutations(.5, phi/2, (1+phi)/2)...)

	return p
}


func RhombicuboctahedronPoints() []Vector3 {
	return signEvenPermutations(1, 1, 1+math.Sqrt(2))
}

func RhombicosidodecahedronPoints() []Vector3 {

	p := []Vector3{}

	p = append(p, signEvenPermutations(1, 1, phi3)...)
	p = append(p, signEvenPermutations(phi2, phi, 2*phi)...)
	p = append(p, signEvenPermutations(2+phi, 0, phi2)...)

	return p
}

func SnubCubePoints() []Vector3 {

	p := []Vector3{}

	xi := ( math.Pow( 17.+3.*math.Sqrt(33.), 1./3.) -
            math.Pow(-17.+3.*math.Sqrt(33.), 1./3.) - 1.) / 3.

	evenPlusses := [...]Vector3 {
		{-1,-1,-1 },
		{ 1, 1,-1 },
		{ 1,-1, 1 },
		{-1, 1, 1 },
	}

	coords := [...]Vector3 {
		{ 1, xi, 1/xi },
		{-xi, -1, -1/xi },
	}

	for _, s := range evenPlusses {
		for _, c := range coords {
			p = append(p, evenPermutations(s.X*c.X, s.Y*c.Y, s.Z*c.Z)...)
		}
	}

	return p
}
func SnubDodecahedronPoints() []Vector3 {

	p := []Vector3{}

	xi := ( math.Pow(phi/2.+math.Sqrt(phi-(5./27.))/2., 1./3.) +
	        math.Pow(phi/2.-math.Sqrt(phi-(5./27.))/2., 1./3.) )

	a := xi - 1/xi
	b := xi*phi + phi2 + phi/xi

	evenPlusses := [...]Vector3 {
		{-1,-1,-1 },
		{ 1, 1,-1 },
		{ 1,-1, 1 },
		{-1, 1, 1 },
	}

	coords := [...]Vector3 {
		{ 2*a, 2, 2*b },
		{ (a + b/phi + phi),    (-a*phi + b + 1/phi), (a/phi + b*phi - 1) },
        { (-a/phi + b*phi + 1), (-a + b/phi - phi),   (a*phi + b - 1/phi) },
        { (-a/phi + b*phi - 1), ( a - b/phi - phi),   (a*phi + b + 1/phi) },
        { (a + b/phi - phi),    (a*phi - b + 1/phi),  (a/phi + b*phi + 1) },
	}

	for _, s := range evenPlusses {
		for _, c := range coords {
			p = append(p, evenPermutations(s.X*c.X, s.Y*c.Y, s.Z*c.Z)...)
		}
	}

	return p
}

/* Catalan solids */

func TriakisTetrahedronFaces() []Face { /*  kT */
    return CatalanDualFaces(TruncatedTetrahedronPoints())
}

func TriakisOctahedronFaces() []Face { /*  kO */
    return CatalanDualFaces(TruncatedCubePoints())
}

func DisdyakisDodecahedronFaces() []Face { /*  mC */
    return CatalanDualFaces(TruncatedCuboctahedronPoints())
}

func TetrakisHexahedronFaces() []Face { /*  kC */
    return CatalanDualFaces(TruncatedOctahedronPoints())
}

func TriakisIcosahedronFaces() []Face { /*  kI */
    return CatalanDualFaces(TruncatedDodecahedronPoints())
}

func DisdyakisTriacontahedronFaces() []Face { /*  mD */
    return CatalanDualFaces(TruncatedIcosidodecahedronPoints())
}

func PentakisDodecahedronFaces() []Face { /*  kD */
    return CatalanDualFaces(TruncatedIcosahedronPoints())
}

func RhombicDodecahedronFaces() []Face { /*  jC */
    return CatalanDualFaces(CuboctahedronPoints())
}

func RhombicTriacontahedronFaces() []Face { /*  jD */
    return CatalanDualFaces(IcosidodecahedronPoints())
}

func DeltoidalIcositetrahedronFaces() []Face { /*  oC */
    return CatalanDualFaces(RhombicuboctahedronPoints())
}

func DeltoidalHexecontahedronFaces() []Face { /*  oD */
    return CatalanDualFaces(RhombicosidodecahedronPoints())
}

func PentagonalIcositetrahedronFaces() []Face { /*  gC */
    return CatalanDualFaces(SnubCubePoints())
}

func PentagonalHexecontahedronFaces() []Face { /*  gD */
    return CatalanDualFaces(SnubDodecahedronPoints())
}

