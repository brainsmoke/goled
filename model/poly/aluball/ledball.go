package aluball

import (
	"post6.net/goled/model"
	"post6.net/goled/model/poly"
	"post6.net/goled/polyhedron"
	"post6.net/goled/vector"
)

const (
	BottomLeft  = 0
	Bottom      = 1
	BottomRight = 2
	TopRight    = 3
	TopLeft     = 4
)

const innerRadius = 100

var traversal = polyhedron.RemapReorientRoute{

	polyhedron.RemapStep{  0, TopLeft,      -1 }, //  1
	polyhedron.RemapStep{  1, BottomRight,  -1 }, //  2
	polyhedron.RemapStep{  2, TopRight,     -1 }, //  3
	polyhedron.RemapStep{  3, TopRight,     -1 }, //  4
	polyhedron.RemapStep{  4, BottomLeft,   -1 }, //  5
	polyhedron.RemapStep{  5, TopLeft,      -1 }, //  6
	polyhedron.RemapStep{  6, TopLeft,      -1 }, //  7

	polyhedron.RemapStep{  0, BottomRight,  -1 }, //  8
	polyhedron.RemapStep{  8, TopRight,     -1 }, //  9
	polyhedron.RemapStep{  9, BottomRight,  -1 }, // 10
	polyhedron.RemapStep{ 10, Bottom,       -1 }, // 11
	polyhedron.RemapStep{ 11, BottomRight,  -1 }, // 12
	polyhedron.RemapStep{ 12, TopRight,     -1 }, // 13
	polyhedron.RemapStep{ 13, TopRight,     -1 }, // 14
	polyhedron.RemapStep{ 14, TopRight,     -1 }, // 15

	polyhedron.RemapStep{  0, TopRight,     -1 }, // 16
	polyhedron.RemapStep{ 16, TopRight,     -1 }, // 17
	polyhedron.RemapStep{ 16, BottomRight,  -1 }, // 18
	polyhedron.RemapStep{ 18, TopRight,     -1 }, // 19
	polyhedron.RemapStep{ 19, BottomRight,  -1 }, // 20
	polyhedron.RemapStep{ 17, BottomRight,  -1 }, // 21
	polyhedron.RemapStep{ 19, TopRight,     -1 }, // 22
	polyhedron.RemapStep{ 22, TopRight,     -1 }, // 23

}

var ledPositions = []poly.FacePosition { // in mm, polygon point 0 defined as up, Center == (0, 0)

	{   0.000,   0.000,  true },
}

func polyhedronePositions() [][]poly.FacePosition {

	facesList := make([][]poly.FacePosition, 60)

	for i := range facesList {
		facesList[i] = ledPositions;
	}

	return facesList
}

func ledGroups(leds []model.Led3D, solid polyhedron.Solid) map[string][]int {

	groups := make(map[string][]int)

	groups["faces"] = make([]int, len(leds))

	for i := range leds {
		groups["faces"][i] = i
	}

	return groups
}

var ledball *model.Model3D

func cacheLedball() {

	solid := polyhedron.RemapReorientSolid(polyhedron.PentagonalIcositetrahedron(), 0, traversal)
	solid.Scale( innerRadius / solid.Faces[0].Center.Magnitude() )

	north := solid.Faces[0].Normal
	center := vector.Vector3{ 0, 0, 0 }
	eye := solid.Points[solid.Faces[0].Polygon[0]].Sub(solid.Points[solid.Faces[0].Polygon[1]]).Normalize()

	solid.Rotate(eye, center, north)

	ledball = new(model.Model3D)
	ledball.Leds = poly.PopulateLeds(solid, polyhedronePositions())
	ledball.Groups = ledGroups(ledball.Leds, solid)
}

func Ledball() *model.Model3D {

	return ledball.Copy()
}

func init() {

	cacheLedball()
}

