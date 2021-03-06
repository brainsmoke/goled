package greatcircles2

import (
	"post6.net/goled/model"
//	"post6.net/goled/model/poly"
//	"post6.net/goled/polyhedron"
	"post6.net/goled/vector"
	"math"
)

var rings = [6][10]int {

    {  0,  1,  2,  3,  4,  5,  6,  7,  8,  9 },
    {  1, 10, 11, 12, 13,  6, 14, 15, 16, 17 },
    { 10, 18, 19, 20,  7, 14, 21, 22, 23,  2 },
    { 18, 24,  9, 25, 15, 21, 26,  4, 27, 11 },
    { 24,  0, 17, 28, 22, 26,  5, 13, 29, 19 },
    { 28, 23,  3, 27, 12, 29, 20,  8, 25, 16 },
}

var ringsInv = [30][2][2]int {

	{{0, 0}, {4, 1}}, {{0, 1}, {1, 0}}, {{0, 2}, {2, 9}},
	{{0, 3}, {5, 2}}, {{0, 4}, {3, 7}}, {{0, 5}, {4, 6}},
	{{0, 6}, {1, 5}}, {{0, 7}, {2, 4}}, {{0, 8}, {5, 7}},
	{{0, 9}, {3, 2}}, {{1, 1}, {2, 0}}, {{1, 2}, {3, 9}},
	{{1, 3}, {5, 4}}, {{1, 4}, {4, 7}}, {{1, 6}, {2, 5}},
	{{1, 7}, {3, 4}}, {{1, 8}, {5, 9}}, {{1, 9}, {4, 2}},
	{{2, 1}, {3, 0}}, {{2, 2}, {4, 9}}, {{2, 3}, {5, 6}},
	{{2, 6}, {3, 5}}, {{2, 7}, {4, 4}}, {{2, 8}, {5, 1}},
	{{3, 1}, {4, 0}}, {{3, 3}, {5, 8}}, {{3, 6}, {4, 5}},
	{{3, 8}, {5, 3}}, {{4, 3}, {5, 0}}, {{4, 8}, {5, 5}},
}

var facesInward = []int { 21, 0, 0, 21, 21, 3, 3, 21, 21, 4, 4, 21, 30, 4, 4, 30, 30, 15, 15, 30, 30, 14, 14, 30, 29, 14, 14, 29, 29, 13, 13, 29, 29, 7, 7, 29, 29, 6, 6, 29, 25, 6, 6, 25, 24, 6, 6, 24, 24, 7, 7, 24, 28, 7, 7, 28, 28, 13, 13, 28, 1, 21, 21, 1, 2, 21, 21, 2, 2, 30, 30, 2, 5, 30, 30, 5, 5, 29, 29, 5, 5, 25, 25, 5, 2, 25, 25, 2, 1, 25, 25, 1, 18, 25, 25, 18, 18, 24, 24, 18, 17, 24, 24, 17, 8, 24, 24, 8, 8, 28, 28, 8, 9, 28, 28, 9, 10, 28, 28, 10, 0, 20, 20, 0, 16, 20, 20, 16, 16, 23, 23, 16, 12, 23, 23, 12, 12, 27, 27, 12, 12, 22, 22, 12, 16, 22, 22, 16, 0, 22, 22, 0, 3, 22, 22, 3, 3, 26, 26, 3, 4, 26, 26, 4, 15, 26, 26, 15, 15, 31, 31, 15, 14, 31, 31, 14, 13, 31, 31, 13, 20, 1, 1, 20, 20, 18, 18, 20, 20, 17, 17, 20, 23, 17, 17, 23, 23, 8, 8, 23, 23, 9, 9, 23, 27, 9, 9, 27, 27, 10, 10, 27, 27, 11, 11, 27, 27, 19, 19, 27, 22, 19, 19, 22, 26, 19, 19, 26, 26, 11, 11, 26, 31, 11, 11, 31, 31, 10, 31, 10}

var trianglesPentagons = []int {
0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0,
}

var phi float64 = (1+math.Sqrt(5))/2.

var a, b, c float64 = 1/2./phi, 1/2., phi/2

var positions = [30]vector.Vector3 {

    { 0, 1, 0}, {-a, c,-b}, {-b, a,-c}, {-b,-a,-c}, {-a,-c,-b}, { 0,-1, 0},
    { a,-c, b}, { b,-a, c}, { b, a, c}, { a, c, b}, {-c, b,-a}, {-1, 0, 0},
    {-c,-b, a}, {-a,-c, b}, { c,-b, a}, { 1, 0, 0}, { c, b,-a}, { a, c,-b},
    {-c, b, a}, {-b, a, c}, { 0, 0, 1}, { c,-b,-a}, { b,-a,-c}, { 0, 0,-1},
    {-a, c, b}, { c, b, a}, { a,-c,-b}, {-c,-b,-a}, { b, a,-c}, {-b,-a, c},
}

var edges = [60][2]int {

		{17,  0}, {17, 28}, {22, 28}, {22, 26}, { 5, 26},
		{26,  4}, {26, 21}, {23, 22}, {21, 22}, {16, 17},
		{16, 15}, {25, 16}, { 1, 17}, {28, 23}, {28, 16},

		{ 0,  9}, { 8,  9}, { 8,  7}, { 6,  7}, { 6,  5},
		{ 9, 25}, {15, 25}, {14, 15}, {21, 14}, {14,  6},
		{15, 21}, { 7, 14}, { 7, 20}, {25,  8}, {20,  8},

		{ 0,  1}, { 2,  1}, { 2,  3}, { 4,  3}, { 4,  5},
		{ 1, 10}, {11, 10}, {27, 11}, {12, 27}, {27,  4},
		{11, 12}, { 3, 27}, { 3, 23}, {10,  2}, {23,  2},

		{24,  0}, {24, 19}, {29, 19}, {29, 13}, { 5, 13},
		{13,  6}, {13, 12}, {20, 29}, {12, 29}, {18, 24},
		{18, 11}, {10, 18}, { 9, 24}, {19, 20}, {19, 18},
}

var ledPositions = [][3]float64 {

	{ 118.522, -18.773,  1 },
	{ 118.522, -18.773, -1 },
	{ 118.522,  18.773, -1 },
	{ 118.522,  18.773,  1 },
}

var ledball *model.Model3D

func getLeds(vStart, vEnd vector.Vector3, faceID int) []model.Led3D {

	leds := make([]model.Led3D, 4)

	normal := vEnd.CrossProduct(vStart).Normalize()
	left := vEnd.Add(vStart).Normalize()

	top := normal.CrossProduct(left)

	for i := range leds {
		x, y, dir := ledPositions[i][0], ledPositions[i][1], ledPositions[i][2]
		leds[i] = model.Led3D{
			Position: left.Mul(x).Add(top.Mul(y)).Add(normal.Mul(dir)),
			Normal: normal.Mul(dir),
			Face: faceID,
			Inside: false,
		}
	}

	return leds
}

func cacheLedball() {

	ledball = new(model.Model3D)

	ledball.Leds = nil
	ringGroup := []int{}
	ringSideGroup := []int{}
	ledGroup := []int{}

	for i := range(edges) {
		s, e := edges[i][0], edges[i][1]
		r, pos1, pos2 := -1, -1, -1

		for _, rS := range ringsInv[s] {
			for _, rE := range ringsInv[e] {
				if rS[0] == rE[0] {
					r = rS[0]
					pos1 = rS[1]
					pos2 = rE[1]
				}
			}
		}

		ledball.Leds = append(ledball.Leds, getLeds(positions[s], positions[e], i)...)

		ringGroup = append(ringGroup, []int{ r, r, r, r }...)
		front, back := 0, 1
		if pos1 == (pos2+1)%10 {
			front, back = 1, 0
		}
		ringSideGroup = append(ringSideGroup, []int{ r*2+front, r*2+back, r*2+front, r*2+back }...)
		ledGroup = append(ledGroup, []int{ 0, 1, 2, 3 }...)
	}

	ledball.Groups = make(map[string][]int)
	ledball.Groups["circles"] = ringGroup
	ledball.Groups["sides"] = ringSideGroup
	ledball.Groups["facets"] = append([]int(nil), facesInward...)
	ledball.Groups["facets_types"] = append([]int(nil), trianglesPentagons...)
	ledball.Groups["leds"] = ledGroup
}

func Ledball() *model.Model3D {

	return ledball.Copy()
}

func init() {

	cacheLedball()
}

