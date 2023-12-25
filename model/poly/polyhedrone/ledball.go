package polyhedrone

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

const (
	LedTop    = 0
	LedRight  = 1
	LedBottom = 2
	LedLeft   = 3
	LedInside = 4
)

const innerRadius = 160

var traversal = polyhedron.RemapRoute{

	BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, TopLeft,

	BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, BottomLeft, BottomLeft, TopLeft,
	TopLeft,
	BottomLeft, BottomLeft, TopLeft,
}

var ledPositions = []poly.FacePosition { // in mm, polygon point 0 defined as up, Center == (0, 0)

	{   5.381,  30.860, false }, { -17.972,  21.880, false }, {  -1.531,  15.508, false }, /* top */
	{  17.329,  16.047, false }, {  34.570,  13.711, false }, {  29.542,  -9.098, false }, /* right */
	{  -4.494, -29.663, false }, {  12.749, -34.692, false }, {  -1.709, -49.329, false }, /* bottom */
	{ -39.606,   8.504, false }, { -22.184,  -2.811, false }, { -22.633, -19.605, false }, /* right */
	{   8.078,  20.266,  true }, { -30.715,   7.336,  true }, {  -3.954, -33.255,  true }, /* inside */

}

func polyhedronePositions() [][]poly.FacePosition {

	list := []poly.FacePosition(nil)

	for i:=0; i<len(ledPositions); i+=3 {

		a, b, c := ledPositions[i], ledPositions[i+1], ledPositions[i+2]

		list = append(list, poly.FacePosition{ (a.X+b.X+c.X)/3., (a.Y+b.Y+c.Y)/3., a.Inside } )
	}

	facesList := make([][]poly.FacePosition, 60)

	for i := range facesList {
		facesList[i] = list
	}

	return facesList
}

func ledNeighbours(solid polyhedron.Solid) [][]int {

	neighbourList := make([][]int, 300)

	for i, f := range solid.Faces {

		for j := i*5 ; j < (i+1)*5 ; j++ {
			neighbourList[j] = make([]int, 4)
		}

		tl, tr := f.Neighbours[TopLeft], f.Neighbours[TopRight]
		bl, br := f.Neighbours[BottomLeft], f.Neighbours[BottomRight]

		neighbourList[i*5+LedTop][TopLeft] = tl*5 + LedTop
		neighbourList[i*5+LedTop][TopRight] = tr*5 + LedTop
		neighbourList[i*5+LedTop][BottomLeft] = i*5 + LedLeft
		neighbourList[i*5+LedTop][BottomRight] = i*5 + LedRight

		neighbourList[i*5+LedLeft][TopLeft] = tl*5 + LedRight
		neighbourList[i*5+LedLeft][TopRight] = i*5 + LedTop
		neighbourList[i*5+LedLeft][BottomLeft] = bl*5 + LedRight
		neighbourList[i*5+LedLeft][BottomRight] = i*5 + LedBottom

		neighbourList[i*5+LedRight][TopLeft] = i*5 + LedTop
		neighbourList[i*5+LedRight][TopRight] = tr*5 + LedLeft
		neighbourList[i*5+LedRight][BottomLeft] = i*5 + LedBottom
		neighbourList[i*5+LedRight][BottomRight] = br*5 + LedLeft

		neighbourList[i*5+LedBottom][TopLeft] = i*5 + LedLeft
		neighbourList[i*5+LedBottom][TopRight] = i*5 + LedRight
		neighbourList[i*5+LedBottom][BottomLeft] = bl*5 + LedBottom
		neighbourList[i*5+LedBottom][BottomRight] = br*5 + LedBottom

		neighbourList[i*5+LedInside][TopLeft] = tl*5 + LedInside
		neighbourList[i*5+LedInside][TopRight] = tr*5 + LedInside
		neighbourList[i*5+LedInside][BottomLeft] = bl*5 + LedInside
		neighbourList[i*5+LedInside][BottomRight] = br*5 + LedInside
	}

	return neighbourList
}

func ledGroups(leds []model.Led3D, solid polyhedron.Solid) map[string][]int {

	groups := make(map[string][]int)

	groups["icosaedron"] = make([]int, len(leds))
	groups["dodecahedron"] = make([]int, len(leds))
	groups["faces"] = make([]int, len(leds))
	groups["rhomb"] = make([]int, len(leds))
	groups["leds"] = make([]int, len(leds))

	icosaedron := make([]int, 60)
	dodecahedron := make([]int, 60)

	for i := 0; i < 60; i++ {
		icosaedron[i] = -1
		dodecahedron[i] = -1
	}

	icoCount, dodeCount := 0, 0
	for i := 0; i < 60; i++ {
		if icosaedron[i] == -1 {

			for j := i; icosaedron[j] == -1; j = solid.Faces[j].Neighbours[TopLeft] {
				icosaedron[j] = icoCount
			}
			icoCount++
		}
		if dodecahedron[i] == -1 {

			for j := i; dodecahedron[j] == -1; j = solid.Faces[j].Neighbours[BottomLeft] {
				dodecahedron[j] = dodeCount
			}
			dodeCount++
		}
	}

	rhombCount := icoCount + dodeCount
	for i, l := range leds {

		f := l.Face
		groups["icosaedron"][i] = icosaedron[f]
		groups["dodecahedron"][i] = dodecahedron[f]
		groups["faces"][i] = f
		if i%5 == 0 {
			groups["rhomb"][i] = icosaedron[f]
		} else if i%5 == 2 {
			groups["rhomb"][i] = icoCount + dodecahedron[f]
		} else if i%5 == 1 && groups["rhomb"][i] == 0 {

			topRight := solid.Faces[f].Neighbours[TopRight]
			bottomRight := solid.Faces[f].Neighbours[BottomRight]
			opposite := solid.Faces[topRight].Neighbours[BottomLeft]

			groups["rhomb"][i] = rhombCount
			groups["rhomb"][5*topRight+3] = rhombCount
			groups["rhomb"][5*opposite+1] = rhombCount
			groups["rhomb"][5*bottomRight+3] = rhombCount

			rhombCount++
		}
		groups["leds"][i] = i % 5
	}

	return groups
}

var ledball *model.Model3D

func cacheLedball() {

	solid := polyhedron.RemapSolid(polyhedron.DeltoidalHexecontahedron(), 0, traversal)
	solid.Scale( innerRadius / solid.Faces[0].Center.Magnitude() )

	north := solid.Points[solid.Faces[32].Polygon[2]].Normalize()
	center := vector.Vector3{ 0, 0, 0 }
	eye := north.CrossProduct(solid.Points[solid.Faces[32].Polygon[0]]).Normalize()

	solid.Rotate(eye, center, north)

	ledball = new(model.Model3D)
	ledball.Leds = poly.PopulateLeds(solid, polyhedronePositions())
	ledball.Neighbours = ledNeighbours(solid)
	ledball.Groups = ledGroups(ledball.Leds, solid)
}

func Ledball() *model.Model3D {

	return ledball.Copy()
}

func init() {

	cacheLedball()
}

