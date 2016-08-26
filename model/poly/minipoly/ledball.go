package minipoly

import (
	"post6.net/goled/model"
	"post6.net/goled/model/poly"
	"post6.net/goled/polyhedron"
)

const (
	TopLeft     = 0
	BottomLeft  = 1
	BottomRight = 2
	TopRight    = 3
)

const (
	LedLeft   = 0
	LedBottom = 1
	LedMiddle = 2
	LedTop    = 3
	LedRight  = 4
)

const innerRadius = 54

var traversal = polyhedron.RemapRoute{

	BottomRight, TopRight,
	BottomRight, TopRight,
	BottomRight, TopRight,
	TopRight,
	BottomRight, BottomRight, BottomRight, TopRight,
	BottomRight,

	BottomRight, TopRight,
	BottomRight, TopRight,
	BottomRight, TopRight,
	TopRight,
	BottomRight, BottomRight, BottomRight, TopRight,
	BottomRight,
}

var ledPositions = []poly.FacePosition { // in mm, polygon point 0 defined as up, Center == (0, 0)

	{ -17.898,  -1.141, false },
	{   2.961, -18.257, false },
	{   0.840,   1.992, false },
	{  -3.050,  12.503, false },
	{  18.517,   4.371, false },

}

func polyhedronePositions() [][]poly.FacePosition {

	facesList := make([][]poly.FacePosition, 24)

	for i := range facesList {
		facesList[i] = ledPositions
	}

	return facesList
}

func ledNeighbours(faces []polyhedron.Face) [][]int {

	neighbourList := make([][]int, 120)

	for i, f := range faces {

		for j := i*5 ; j < (i+1)*5 ; j++ {
			neighbourList[j] = make([]int, 4)
		}

		tl, tr := f.Neighbours[TopLeft], f.Neighbours[TopRight]
		bl, br := f.Neighbours[BottomLeft], f.Neighbours[BottomRight]

		neighbourList[i*5+LedTop][TopLeft] = tl*5 + LedTop
		neighbourList[i*5+LedTop][TopRight] = tr*5 + LedTop
		neighbourList[i*5+LedTop][BottomLeft] = i*5 + LedLeft
		neighbourList[i*5+LedTop][BottomRight] = i*5 + LedMiddle

		neighbourList[i*5+LedLeft][TopLeft] = tl*5 + LedRight
		neighbourList[i*5+LedLeft][TopRight] = i*5 + LedTop
		neighbourList[i*5+LedLeft][BottomLeft] = bl*5 + LedRight
		neighbourList[i*5+LedLeft][BottomRight] = i*5 + LedMiddle

		neighbourList[i*5+LedRight][TopLeft] = i*5 + LedTop
		neighbourList[i*5+LedRight][TopRight] = tr*5 + LedLeft
		neighbourList[i*5+LedRight][BottomLeft] = i*5 + LedMiddle
		neighbourList[i*5+LedRight][BottomRight] = br*5 + LedLeft

		neighbourList[i*5+LedBottom][TopLeft] = i*5 + LedLeft
		neighbourList[i*5+LedBottom][TopRight] = i*5 + LedMiddle
		neighbourList[i*5+LedBottom][BottomLeft] = bl*5 + LedBottom
		neighbourList[i*5+LedBottom][BottomRight] = br*5 + LedBottom

		neighbourList[i*5+LedMiddle][TopLeft] = i*5 + LedLeft
		neighbourList[i*5+LedMiddle][TopRight] = i*5 + LedTop
		neighbourList[i*5+LedMiddle][BottomLeft] = i*5 + LedBottom
		neighbourList[i*5+LedMiddle][BottomRight] = i*5 + LedRight
	}

	return neighbourList
}

func ledGroups(leds []model.Led3D, faces []polyhedron.Face) map[string][]int {

	groups := make(map[string][]int)

	groups["cube"] = make([]int, len(leds))
	groups["octahedron"] = make([]int, len(leds))
	groups["faces"] = make([]int, len(leds))
	groups["rhomb"] = make([]int, len(leds))
	groups["leds"] = make([]int, len(leds))

	cube := [24]int { 0, 0, 1, 1, 2, 2, 3, 4, 4, 4, 4, 3, 3, 3, 2, 2, 1, 1, 0, 5, 5, 5, 5, 0 }
	octahedron := [24]int { 0, 1, 1, 2, 2, 3, 3, 3, 2, 1, 0, 0, 4, 5, 5, 6, 6, 7, 7, 7, 6, 5, 4, 4 }

	rhombCount := 6+8
	for i, l := range leds {

		f := l.Face
		groups["cube"][i] = cube[f]
		groups["octahedron"][i] = octahedron[f]
		groups["faces"][i] = f
		if i%5 == LedBottom {
			groups["rhomb"][i] = cube[f]
		} else if i%5 == LedTop {
			groups["rhomb"][i] = 8 + octahedron[f]
		} else if i%5 == LedRight && groups["rhomb"][i] == 0 {

			topRight := faces[f].Neighbours[TopRight]
			bottomRight := faces[f].Neighbours[BottomRight]
			opposite := faces[topRight].Neighbours[BottomLeft]

			groups["rhomb"][i] = rhombCount
			groups["rhomb"][5*topRight+LedLeft] = rhombCount
			groups["rhomb"][5*opposite+LedRight] = rhombCount
			groups["rhomb"][5*bottomRight+LedLeft] = rhombCount

			rhombCount++
		}
		groups["leds"][i] = i % 5
	}

	return groups
}

var ledball *model.Model3D

func cacheLedball() {

	faces := polyhedron.RemapFaces(polyhedron.DeltoidalIcositetrahedronFaces(), 0, traversal)
	factor := innerRadius / faces[0].Center.Magnitude()
	faces = polyhedron.Scale(faces, factor)

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

