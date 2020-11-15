package projection

import (
	"math"
	"post6.net/goled/util"
	"post6.net/goled/vector"
)

func Voronoi(width, height int, points []vector.Vector3) ([]int, []float64) {

	sinY := util.TrigTable(height, math.Pi*.5/float64(height), math.Pi*(float64(height)+.5)/float64(height), math.Sin)
	cosY := util.TrigTable(height, math.Pi*.5/float64(height), math.Pi*(float64(height)+.5)/float64(height), math.Cos)
	sinX := util.TrigTable(width, 0, 2*math.Pi, math.Sin)
	cosX := util.TrigTable(width, 0, 2*math.Pi, math.Cos)

	vmap := make([]int, width*height)

	for ty := 0; ty < height; ty++ {
		for tx := 0; tx < width; tx++ {

			ay := cosY[ty]
			ax := sinY[ty] * sinX[tx]
			az := sinY[ty] * cosX[tx]

			best := -1
			bestd2 := 5.

			for i, v := range points {
				dx, dy, dz := ax-v.X, ay-v.Y, az-v.Z
				d2 := dx*dx + dy*dy + dz*dz
				if d2 < bestd2 {
					best = i
					bestd2 = d2
				}
			}
			vmap[ty*width+tx] = best
		}
	}

	return vmap, sinY
}
