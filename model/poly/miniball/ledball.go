package miniball

import (
	"post6.net/goled/model"
	"post6.net/goled/model/poly"
	"post6.net/goled/polyhedron"
	"post6.net/goled/vector"
)

const (
	TopLeft     = 0
	BottomLeft  = 1
	BottomRight = 2
	TopRight    = 3
)

const innerRadius = 1

var traversal = polyhedron.RemapRoute{

BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight,
BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight,
//TopRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight,
}



var ledPositions = []poly.FacePosition { // in mm, polygon point 0 defined as up, Center == (0, 0)

	{ 0, 0, false },

}

func polyhedronePositions() [][]poly.FacePosition {

	facesList := make([][]poly.FacePosition, 60)

	for i := range facesList {
		facesList[i] = ledPositions;
	}

	return facesList
}

func ledNeighbours(faces []polyhedron.Face) [][]int {

	neighbourList := make([][]int, 60)
	for i := range neighbourList {
		neighbourList[i] = make([]int, 4)
	}

	for i, f := range faces {

		neighbourList[i][TopLeft] = f.Neighbours[TopLeft]
		neighbourList[i][TopRight] = f.Neighbours[TopRight]
		neighbourList[i][BottomLeft] = f.Neighbours[BottomLeft]
		neighbourList[i][BottomRight] = f.Neighbours[BottomRight]

	}
	return neighbourList
}

func ledGroups(leds []model.Led3D, faces []polyhedron.Face) map[string][]int {

	groups := make(map[string][]int)

	groups["leds"] = make([]int, len(leds))

	for i := range leds {
		groups["leds"][i] = i
	}

	return groups
}

var ledball *model.Model3D

func cacheLedball() {

	faces := polyhedron.RemapFaces(polyhedron.DeltoidalHexecontahedronFaces(), 0, traversal)
	factor := innerRadius / faces[0].Center.Magnitude()
	faces = polyhedron.Scale(faces, factor)

	north := faces[5].Polygon[2].Normalize()
	center := vector.Vector3{ 0, 0, 0 }
	eye := north.CrossProduct(faces[5].Polygon[0]).Normalize()

	faces = polyhedron.Rotate(faces, eye, center, north)

	ledball = new(model.Model3D)
	ledball.Leds = poly.PopulateLeds(faces, polyhedronePositions())
	ledball.Neighbours = ledNeighbours(faces)
	ledball.Groups = ledGroups(ledball.Leds, faces)
}

func Ledball() *model.Model3D {

	return ledball.Copy()
}

func init() {

	cacheLedball()
}

