package polyhedrone

import (
	"post6.net/goled/model"
	"post6.net/goled/polyhedron"
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

var faceTraversal = [...]struct {
	x, y   float64
	inside bool
}{

	{.04, .17, false},  // top
	{.27, .43, false},  // right
	{0, .73, false},    // bottom
	{-.25, .33, false}, // left
	{.04, .40, true},   // inside
}

func faceLeds(f polyhedron.Face, index int) []model.Led3D {

	top, left, bottom, right := f.Polygon[0], f.Polygon[1], f.Polygon[2], f.Polygon[3]
	leds := make([]model.Led3D, len(faceTraversal))

	for i, p := range faceTraversal {

		vY := bottom.Sub(top)
		vX := right.Sub(left).Normalize().Mul(top.Distance(bottom))
		pos := top.Add(vX.Mul(p.x)).Add(vY.Mul(p.y))
		normal := f.Normal
		if p.inside {
			normal = normal.Mul(-1)
		}
		leds[i] = model.Led3D{pos, normal, index, p.inside}
	}

	return leds
}

func ledNeighbours(faces []polyhedron.Face) [][]int {

	neighbourList := make([][]int, 300)

	for i, f := range faces {

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

func ledGroups(leds []model.Led3D, faces []polyhedron.Face) map[string][]int {

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

			for j := i; icosaedron[j] == -1; j = faces[j].Neighbours[TopLeft] {
				icosaedron[j] = icoCount
			}
			icoCount++
		}
		if dodecahedron[i] == -1 {

			for j := i; dodecahedron[j] == -1; j = faces[j].Neighbours[BottomLeft] {
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
		if i%5 == 0 {
			groups["rhomb"][i] = icosaedron[f]
		} else if i%5 == 2 {
			groups["rhomb"][i] = icoCount + dodecahedron[f]
		} else if i%5 == 1 && groups["rhomb"][i] == 0 {

			topRight := faces[f].Neighbours[TopRight]
			bottomRight := faces[f].Neighbours[BottomRight]
			opposite := faces[topRight].Neighbours[BottomLeft]

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

	ledball = new(model.Model3D)

	faces := polyhedron.RemapFaces(polyhedron.DeltoidalHexecontahedronFaces(), 0, traversal)

	ledball.Leds = make([]model.Led3D, 300)

	for i, f := range faces {

		copy(ledball.Leds[i*5:i*5+5], faceLeds(f, i))
	}

	ledball.Neighbours = ledNeighbours(faces)
	ledball.Groups = ledGroups(ledball.Leds, faces)

	ledball = ledball.UnitScale()
}

func Ledball() *model.Model3D {

	return ledball.Copy()
}

func init() {

	cacheLedball()
}

