package model

import (
	"math"
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

var traversal = [...]int{

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

func remapFaces(faces []polyhedron.Face) []polyhedron.Face {

	newFaces := make([]polyhedron.Face, len(faces))
	mapping := make([]int, len(faces))

	current := 0

	for i := range newFaces {

		newFaces[i] = faces[current]
		mapping[current] = i
		current = faces[current].Neighbours[traversal[i]]
	}

	for i := range newFaces {
		for j, k := range newFaces[i].Neighbours {

			newFaces[i].Neighbours[j] = mapping[k]
		}
	}

	return newFaces
}

func faceLeds(f polyhedron.Face, index int) []Led3D {

	top, left, bottom, right := f.Polygon[0], f.Polygon[1], f.Polygon[2], f.Polygon[3]
	leds := make([]Led3D, len(faceTraversal))

	for i, p := range faceTraversal {

		vY := bottom.Sub(top)
		vX := right.Sub(left).Normalize().Mul(top.Distance(bottom))
		pos := top.Add(vX.Mul(p.x)).Add(vY.Mul(p.y))
		normal := f.Normal
		if p.inside {
			normal = normal.Mul(-1)
		}
		leds[i] = Led3D{pos, normal, index, p.inside}
	}

	return leds
}

var ledballRawCached = false
var ledballRaw []Led3D

var ledballCached = false
var ledball []Led3D

var ledballSmoothCached = false
var ledballSmooth []Led3D

var ledballFacesCached = false
var ledballFaces []polyhedron.Face

func LedballFaces() []polyhedron.Face {

	if !ledballFacesCached {
		ledballFaces = remapFaces(polyhedron.DeltoidalhexecontahedronFaces())
		ledballFacesCached = true
	}
	return append([]polyhedron.Face(nil), ledballFaces...)
}

func LedballLedNeighbours() [][4]int {

	neighbourList := make([][4]int, 300)

	faces := LedballFaces()

	for i, f := range faces {
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

func LedballRaw() []Led3D {

	if !ledballRawCached {

		ledballRaw = make([]Led3D, 300)
		poly := LedballFaces()

		for i, f := range poly {

			copy(ledballRaw[i*5:i*5+5], faceLeds(f, i))
		}

		ledballRawCached = true
	}

	return append([]Led3D(nil), ledballRaw...)
}

func Ledball() []Led3D {

	if !ledballCached {

		ledball = LedballRaw()

		max := 0.
		for _, p := range ledball {
			max = math.Max(max, p.Position.Magnitude())
		}

		for i := range ledball {
			ledball[i].Position = ledball[i].Position.Mul(1. / max)
		}

		ledballCached = true
	}

	return append([]Led3D(nil), ledball...)
}

func LedballSmooth() []Led3D {

	if !ledballSmoothCached {

		ledballSmooth = LedballRaw()

		for i := range ledball {
			p := ledballSmooth[i].Position.Normalize()
			ledballSmooth[i].Position = p

			if ledballSmooth[i].Inside {
				ledballSmooth[i].Normal = p.Mul(-1)
			} else {
				ledballSmooth[i].Normal = p
			}
		}

		ledballSmoothCached = true
	}

	return append([]Led3D(nil), ledballSmooth...)
}
