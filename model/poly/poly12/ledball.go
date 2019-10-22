package poly12

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

	{  39.606,   8.504, false }, // U1
	{  30.715,   7.336,  true }, // U2
	{  17.972,  21.880, false }, // U3
	{  -5.381,  30.860, false }, // U4
	{  -8.078,  20.266,  true }, // U5
	{   1.531,  15.508, false }, // U6
	{  22.184,  -2.811, false }, // U7
	{  22.633, -19.605, false }, // U8
	{   4.494, -29.663, false }, // U9
	{   3.954, -33.255,  true }, // U10
	{   1.709, -49.329, false }, // U11
	{ -12.749, -34.692, false }, // U12
	{ -29.542,  -9.098, false }, // U13
	{ -17.329,  16.047, false }, // U14
	{ -34.570,  13.711, false }, // U15

}

func polyhedronePositions() [][]poly.FacePosition {

	facesList := make([][]poly.FacePosition, 60)

	for i := range facesList {
		facesList[i] = ledPositions;
	}

	return facesList
}

func ledNeighbours(solid polyhedron.Solid) [][]int {

	neighbourList := make([][]int, 900)
	for i := range neighbourList {
		neighbourList[i] = make([]int, 4)
	}

	for i, f := range solid.Faces {

		tl, tr := f.Neighbours[TopLeft], f.Neighbours[TopRight]
		bl, br := f.Neighbours[BottomLeft], f.Neighbours[BottomRight]

		neighbourList[i*15 + 0 ][TopLeft] = i*15 + 2
		neighbourList[i*15 + 0 ][TopRight] = tr*15 + 14
		neighbourList[i*15 + 0 ][BottomLeft] = i*15 + 6
		neighbourList[i*15 + 0 ][BottomRight] = br*15 + 12

		neighbourList[i*15 + 2 ][TopLeft] = tr*15 + 3
		neighbourList[i*15 + 2 ][TopRight] = tr*15 + 13
		neighbourList[i*15 + 2 ][BottomLeft] = i*15 + 5
		neighbourList[i*15 + 2 ][BottomRight] = i*15 + 6

		neighbourList[i*15 + 3 ][TopLeft] = tl*15 + 3
		neighbourList[i*15 + 3 ][TopRight] = tr*15 + 3
		neighbourList[i*15 + 3 ][BottomLeft] = i*15 + 13
		neighbourList[i*15 + 3 ][BottomRight] = i*15 + 5

		neighbourList[i*15 + 5 ][TopLeft] = i*15 + 3
		neighbourList[i*15 + 5 ][TopRight] = i*15 + 2
		neighbourList[i*15 + 5 ][BottomLeft] = i*15 + 12
		neighbourList[i*15 + 5 ][BottomRight] = i*15 + 6

		neighbourList[i*15 + 6 ][TopLeft] = i*15 + 5
		neighbourList[i*15 + 6 ][TopRight] = i*15 + 2
		neighbourList[i*15 + 6 ][BottomLeft] = i*15 + 7
		neighbourList[i*15 + 6 ][BottomRight] = br*15 + 12

		neighbourList[i*15 + 7 ][TopLeft] = i*15 + 6
		neighbourList[i*15 + 7 ][TopRight] = br*15 + 12
		neighbourList[i*15 + 7 ][BottomLeft] = i*15 + 8
		neighbourList[i*15 + 7 ][BottomRight] = br*15 + 11

		neighbourList[i*15 + 8 ][TopLeft] = i*15 + 11
		neighbourList[i*15 + 8 ][TopRight] = i*15 + 7
		neighbourList[i*15 + 8 ][BottomLeft] = i*15 + 10
		neighbourList[i*15 + 8 ][BottomRight] = br*15 + 11

		neighbourList[i*15 + 10 ][TopLeft] = i*15 + 11
		neighbourList[i*15 + 10 ][TopRight] = i*15 + 8
		neighbourList[i*15 + 10 ][BottomLeft] = bl*15 + 10
		neighbourList[i*15 + 10 ][BottomRight] = br*15 + 10

		neighbourList[i*15 + 11 ][TopLeft] = i*15 + 12
		neighbourList[i*15 + 11 ][TopRight] = i*15 + 6
		neighbourList[i*15 + 11 ][BottomLeft] = bl*15 + 10
		neighbourList[i*15 + 11 ][BottomRight] = i*15 + 10

		neighbourList[i*15 + 12 ][TopLeft] = bl*15 + 0
		neighbourList[i*15 + 12 ][TopRight] = i*15 + 14
		neighbourList[i*15 + 12 ][BottomLeft] = bl*15 + 7
		neighbourList[i*15 + 12 ][BottomRight] = i*15 + 11

		neighbourList[i*15 + 13 ][TopLeft] = tl*15 + 2
		neighbourList[i*15 + 13 ][TopRight] = tr*15 + 3
		neighbourList[i*15 + 13 ][BottomLeft] = i*15 + 12
		neighbourList[i*15 + 13 ][BottomRight] = i*15 + 5

		neighbourList[i*15 + 14 ][TopLeft] = tl*15 + 0
		neighbourList[i*15 + 14 ][TopRight] = tl*15 + 3
		neighbourList[i*15 + 14 ][BottomLeft] = bl*15 + 0
		neighbourList[i*15 + 14 ][BottomRight] = i*15 + 12

		neighbourList[i*15 + 1 ][TopLeft] = i*15 + 4
		neighbourList[i*15 + 1 ][TopRight] = tr*15 + 4
		neighbourList[i*15 + 1 ][BottomLeft] = i*15 + 10
		neighbourList[i*15 + 1 ][BottomRight] = br*15 + 4

		neighbourList[i*15 + 4 ][TopLeft] = tl*15 + 4
		neighbourList[i*15 + 4 ][TopRight] = tr*15 + 4
		neighbourList[i*15 + 4 ][BottomLeft] = bl*15 + 1
		neighbourList[i*15 + 4 ][BottomRight] = i*15 + 1

		neighbourList[i*15 + 9 ][TopLeft] = i*15 + 4
		neighbourList[i*15 + 9 ][TopRight] = i*15 + 1
		neighbourList[i*15 + 9 ][BottomLeft] = bl*15 + 9
		neighbourList[i*15 + 9 ][BottomRight] = br*15 + 9

	}
	return neighbourList
}

func ledGroups(leds []model.Led3D, solid polyhedron.Solid) map[string][]int {

	groups := make(map[string][]int)

	groups["shapes"] = make([]int, len(leds))
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
		groups["dodecahedron"][i] = f
		groups["faces"][i] = dodecahedron[f]

		switch i%15 {
			case 3:
				groups["shapes"][i] = icosaedron[f]
			case 14, 13, 5, 2:
				groups["shapes"][i] = icosaedron[f] + icoCount
			case 10:
				groups["shapes"][i] = icoCount*2 + dodecahedron[f]
			case 11, 8:
				groups["shapes"][i] = icoCount*2 + dodeCount + dodecahedron[f]
			case 0, 6, 7:
				groups["shapes"][i] = icoCount*2 + dodeCount *2 + dodecahedron[f]
		}


		switch i%15 {

			case 2, 3, 5:
				groups["rhomb"][i] = icosaedron[f]

			case 8, 10, 11:
				groups["rhomb"][i] = icoCount + dodecahedron[f]

			case 0:
				if groups["rhomb"][i] == 0 {

					topRight := solid.Faces[f].Neighbours[TopRight]
					bottomRight := solid.Faces[f].Neighbours[BottomRight]
					opposite := solid.Faces[topRight].Neighbours[BottomLeft]

					groups["rhomb"][i + 0] = rhombCount
					groups["rhomb"][i + 6] = rhombCount
					groups["rhomb"][i + 7] = rhombCount
					groups["rhomb"][15*topRight + 12] = rhombCount
					groups["rhomb"][15*topRight + 13] = rhombCount
					groups["rhomb"][15*topRight + 14] = rhombCount
					groups["rhomb"][15*opposite + 0] = rhombCount
					groups["rhomb"][15*opposite + 6] = rhombCount
					groups["rhomb"][15*opposite + 7] = rhombCount
					groups["rhomb"][15*bottomRight + 12] = rhombCount
					groups["rhomb"][15*bottomRight + 13] = rhombCount
					groups["rhomb"][15*bottomRight + 14] = rhombCount

					rhombCount++
				}
		}
		groups["leds"][i] = i % 15
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

