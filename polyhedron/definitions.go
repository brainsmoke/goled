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

func TetrahedronFaces() []Face { /* T */
    return CatalanDualFaces(TetrahedronPoints())
}

func OctahedronFaces() []Face { /* O */
    return CatalanDualFaces(CubePoints())
}

func CubeFaces() []Face { /* C */
    return ArchimedeanFaces(CubePoints())
}

func DodecahedronFaces() []Face { /* D */
    return CatalanDualFaces(IcosahedronPoints())
}

func IcosahedronFaces() []Face { /* I */
    return ArchimedeanFaces(IcosahedronPoints())
}

/* Archimedean solids */

func TruncatedTetrahedronFaces() []Face { /* tT */
    return ArchimedeanFaces(TruncatedTetrahedronPoints())
}

func TruncatedCubeFaces() []Face { /* tC */
    return ArchimedeanFaces(TruncatedCubePoints())
}

func TruncatedCuboctahedronFaces() []Face { /* bC */
    return ArchimedeanFaces(TruncatedCuboctahedronPoints())
}

func TruncatedOctahedronFaces() []Face { /* tO */
    return ArchimedeanFaces(TruncatedOctahedronPoints())
}

func TruncatedDodecahedronFaces() []Face { /* tD */
    return ArchimedeanFaces(TruncatedDodecahedronPoints())
}

func TruncatedIcosidodecahedronFaces() []Face { /* bD */
    return ArchimedeanFaces(TruncatedIcosidodecahedronPoints())
}

func TruncatedIcosahedronFaces() []Face { /* tI */
    return ArchimedeanFaces(TruncatedIcosahedronPoints())
}

func CuboctahedronFaces() []Face { /* aC */
    return ArchimedeanFaces(CuboctahedronPoints())
}

func IcosidodecahedronFaces() []Face { /* aD */
    return ArchimedeanFaces(IcosidodecahedronPoints())
}

func RhombicuboctahedronFaces() []Face { /* eC */
    return ArchimedeanFaces(RhombicuboctahedronPoints())
}

func RhombicosidodecahedronFaces() []Face { /* eD */
    return ArchimedeanFaces(RhombicosidodecahedronPoints())
}

func SnubCubeFaces() []Face { /* sC */
    return ArchimedeanFaces(SnubCubePoints())
}

func SnubDodecahedronFaces() []Face { /* sD */
    return ArchimedeanFaces(SnubDodecahedronPoints())
}

/* Catalan solids */

func TriakisTetrahedronFaces() []Face { /* kT */
    return CatalanDualFaces(TruncatedTetrahedronPoints())
}

func TriakisOctahedronFaces() []Face { /* kO */
    return CatalanDualFaces(TruncatedCubePoints())
}

func DisdyakisDodecahedronFaces() []Face { /* mC */
    return CatalanDualFaces(TruncatedCuboctahedronPoints())
}

func TetrakisHexahedronFaces() []Face { /* kC */
    return CatalanDualFaces(TruncatedOctahedronPoints())
}

func TriakisIcosahedronFaces() []Face { /* kI */
    return CatalanDualFaces(TruncatedDodecahedronPoints())
}

func DisdyakisTriacontahedronFaces() []Face { /* mD */
    return CatalanDualFaces(TruncatedIcosidodecahedronPoints())
}

func PentakisDodecahedronFaces() []Face { /* kD */
    return CatalanDualFaces(TruncatedIcosahedronPoints())
}

func RhombicDodecahedronFaces() []Face { /* jC */
    return CatalanDualFaces(CuboctahedronPoints())
}

func RhombicTriacontahedronFaces() []Face { /* jD */
    return CatalanDualFaces(IcosidodecahedronPoints())
}

func DeltoidalIcositetrahedronFaces() []Face { /* oC */
    return CatalanDualFaces(RhombicuboctahedronPoints())
}

func DeltoidalHexecontahedronFaces() []Face { /* oD */
    return CatalanDualFaces(RhombicosidodecahedronPoints())
}

func PentagonalIcositetrahedronFaces() []Face { /* gC */
    return CatalanDualFaces(SnubCubePoints())
}

func PentagonalHexecontahedronFaces() []Face { /* gD */
    return CatalanDualFaces(SnubDodecahedronPoints())
}

