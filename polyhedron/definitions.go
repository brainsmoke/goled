package polyhedron

import (
	"math"
	"post6.net/goled/vector"
)

/* Point cloud definitions */

var (
	phi  = (1 + math.Sqrt(5)) / 2
	phi2 = phi * phi
	phi3 = phi * phi * phi
)

func evenPermutations(x, y, z float64) []vector.Vector3 {

	return []vector.Vector3{

		{x, y, z},
		{z, x, y},
		{y, z, x},
	}
}

func allPermutations(x, y, z float64) []vector.Vector3 {

	return []vector.Vector3{

		{x, y, z},
		{x, z, y},
		{y, x, z},
		{y, z, x},
		{z, x, y},
		{z, y, x},
	}
}

func signPermutations(x, y, z float64) []vector.Vector3 {

	p := []vector.Vector3{}

	for _, sx := range []float64{ -x, x } {
		for _, sy := range []float64{ -y, y } {
			for _, sz := range []float64{ -z, z } {
				p = append(p, vector.Vector3{sx, sy, sz})
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

func signEvenPermutations(x, y, z float64) []vector.Vector3 {

	p := []vector.Vector3{}

	for _, v := range signPermutations(x, y, z) {
		p = append(p, evenPermutations(v.X, v.Y, v.Z)...)
	}

	return p
}

func allSignPermutations(x, y, z float64) []vector.Vector3 {

	p := []vector.Vector3{}

	for _, v := range signPermutations(x, y, z) {
		p = append(p, allPermutations(v.X, v.Y, v.Z)...)
	}

	return p
}

/* point lists */

func TetrahedronPoints() []vector.Vector3 {

	return []vector.Vector3 {
		{ -1, 0, -1/math.Sqrt(2) },
		{  1, 0, -1/math.Sqrt(2) },
		{  0,-1,  1/math.Sqrt(2) },
		{  0, 1,  1/math.Sqrt(2) },
	}
}

func CubePoints() []vector.Vector3 {

	return signPermutations(1, 1, 1)
}

func IcosahedronPoints() []vector.Vector3 {

	return signEvenPermutations(0, 1, phi)
}

func TruncatedTetrahedronPoints() []vector.Vector3 {

	p := []vector.Vector3{}

	p = append(p, evenPermutations( 3,  1,  1)...)
	p = append(p, evenPermutations( 3, -1, -1)...)
	p = append(p, evenPermutations(-3,  1, -1)...)
	p = append(p, evenPermutations(-3, -1,  1)...)

	return p
}

func TruncatedCubePoints() []vector.Vector3 {

	return signEvenPermutations( math.Sqrt(2)-1, 1, 1 )
}

func TruncatedCuboctahedronPoints() []vector.Vector3 {

	return allSignPermutations(1, 1+math.Sqrt(2), 1+2*math.Sqrt(2))
}

func TruncatedOctahedronPoints() []vector.Vector3 {

	return allSignPermutations(0, 1, 2)
}

func TruncatedDodecahedronPoints() []vector.Vector3 {

	p := []vector.Vector3{}

	p = append(p, signEvenPermutations(0, 1/phi, 2+phi)...)
	p = append(p, signEvenPermutations(1/phi, phi, 2*phi)...)
	p = append(p, signEvenPermutations(phi, 2, phi2)...)

	return p
}

func TruncatedIcosidodecahedronPoints() []vector.Vector3 {

	p := []vector.Vector3{}

	p = append(p, signEvenPermutations(1/phi, 1/phi, 3+phi)...)
	p = append(p, signEvenPermutations(2/phi, phi, 1+2*phi)...)
	p = append(p, signEvenPermutations(1/phi, phi2, -1+3*phi)...)
	p = append(p, signEvenPermutations(-1+2*phi, 2, 2+phi)...)
	p = append(p, signEvenPermutations(phi, 3, 2*phi)...)

	return p
}

func TruncatedIcosahedronPoints() []vector.Vector3 {

	p := []vector.Vector3{}

	p = append(p, signEvenPermutations(2, 1+2*phi, phi)...)
	p = append(p, signEvenPermutations(1, 2+phi, 2*phi)...)
	p = append(p, signEvenPermutations(1, 3*phi, 0)...)

	return p
}

func CuboctahedronPoints() []vector.Vector3 {

	return signEvenPermutations( 1, 1, 0 )
}

func IcosidodecahedronPoints() []vector.Vector3 {

	p := []vector.Vector3{}

	p = append(p, signEvenPermutations(phi, 0, 0)...)
	p = append(p, signEvenPermutations(.5, phi/2, (1+phi)/2)...)

	return p
}


func RhombicuboctahedronPoints() []vector.Vector3 {
	return signEvenPermutations(1, 1, 1+math.Sqrt(2))
}

func RhombicosidodecahedronPoints() []vector.Vector3 {

	p := []vector.Vector3{}

	p = append(p, signEvenPermutations(1, 1, phi3)...)
	p = append(p, signEvenPermutations(phi2, phi, 2*phi)...)
	p = append(p, signEvenPermutations(2+phi, 0, phi2)...)

	return p
}

func SnubCubePoints() []vector.Vector3 {

	p := []vector.Vector3{}

	xi := ( math.Pow( 17.+3.*math.Sqrt(33.), 1./3.) -
            math.Pow(-17.+3.*math.Sqrt(33.), 1./3.) - 1.) / 3.

	evenPlusses := [...]vector.Vector3 {
		{-1,-1,-1 },
		{ 1, 1,-1 },
		{ 1,-1, 1 },
		{-1, 1, 1 },
	}

	coords := [...]vector.Vector3 {
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
func SnubDodecahedronPoints() []vector.Vector3 {

	p := []vector.Vector3{}

	xi := ( math.Pow(phi/2.+math.Sqrt(phi-(5./27.))/2., 1./3.) +
	        math.Pow(phi/2.-math.Sqrt(phi-(5./27.))/2., 1./3.) )

	a := xi - 1/xi
	b := xi*phi + phi2 + phi/xi

	evenPlusses := [...]vector.Vector3 {
		{-1,-1,-1 },
		{ 1, 1,-1 },
		{ 1,-1, 1 },
		{-1, 1, 1 },
	}

	coords := [...]vector.Vector3 {
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

/* Platonic solids */

func Tetrahedron() Solid { /* T */
    return CatalanDual(TetrahedronPoints())
}

func Octahedron() Solid { /* O */
    return CatalanDual(CubePoints())
}

func Cube() Solid { /* C */
    return Archimedean(CubePoints())
}

func Dodecahedron() Solid { /* D */
    return CatalanDual(IcosahedronPoints())
}

func Icosahedron() Solid { /* I */
    return Archimedean(IcosahedronPoints())
}

/* Archimedean solids */

func TruncatedTetrahedron() Solid { /* tT */
    return Archimedean(TruncatedTetrahedronPoints())
}

func TruncatedCube() Solid { /* tC */
    return Archimedean(TruncatedCubePoints())
}

func TruncatedCuboctahedron() Solid { /* bC */
    return Archimedean(TruncatedCuboctahedronPoints())
}

func TruncatedOctahedron() Solid { /* tO */
    return Archimedean(TruncatedOctahedronPoints())
}

func TruncatedDodecahedron() Solid { /* tD */
    return Archimedean(TruncatedDodecahedronPoints())
}

func TruncatedIcosidodecahedron() Solid { /* bD */
    return Archimedean(TruncatedIcosidodecahedronPoints())
}

func TruncatedIcosahedron() Solid { /* tI */
    return Archimedean(TruncatedIcosahedronPoints())
}

func Cuboctahedron() Solid { /* aC */
    return Archimedean(CuboctahedronPoints())
}

func Icosidodecahedron() Solid { /* aD */
    return Archimedean(IcosidodecahedronPoints())
}

func Rhombicuboctahedron() Solid { /* eC */
    return Archimedean(RhombicuboctahedronPoints())
}

func Rhombicosidodecahedron() Solid { /* eD */
    return Archimedean(RhombicosidodecahedronPoints())
}

func SnubCube() Solid { /* sC */
    return Archimedean(SnubCubePoints())
}

func SnubDodecahedron() Solid { /* sD */
    return Archimedean(SnubDodecahedronPoints())
}

/* Catalan solids */

func TriakisTetrahedron() Solid { /* kT */
    return CatalanDual(TruncatedTetrahedronPoints())
}

func TriakisOctahedron() Solid { /* kO */
    return CatalanDual(TruncatedCubePoints())
}

func DisdyakisDodecahedron() Solid { /* mC */
    return CatalanDual(TruncatedCuboctahedronPoints())
}

func TetrakisHexahedron() Solid { /* kC */
    return CatalanDual(TruncatedOctahedronPoints())
}

func TriakisIcosahedron() Solid { /* kI */
    return CatalanDual(TruncatedDodecahedronPoints())
}

func DisdyakisTriacontahedron() Solid { /* mD */
    return CatalanDual(TruncatedIcosidodecahedronPoints())
}

func PentakisDodecahedron() Solid { /* kD */
    return CatalanDual(TruncatedIcosahedronPoints())
}

func RhombicDodecahedron() Solid { /* jC */
    return CatalanDual(CuboctahedronPoints())
}

func RhombicTriacontahedron() Solid { /* jD */
    return CatalanDual(IcosidodecahedronPoints())
}

func DeltoidalIcositetrahedron() Solid { /* oC */
    return CatalanDual(RhombicuboctahedronPoints())
}

func DeltoidalHexecontahedron() Solid { /* oD */
    return CatalanDual(RhombicosidodecahedronPoints())
}

func PentagonalIcositetrahedron() Solid { /* gC */
    return CatalanDual(SnubCubePoints())
}

func PentagonalHexecontahedron() Solid { /* gD */
    return CatalanDual(SnubDodecahedronPoints())
}

