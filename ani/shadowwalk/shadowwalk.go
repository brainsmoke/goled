package shadowwalk

import (
	"math"
	"post6.net/goled/color"
	"post6.net/goled/model"
	"post6.net/goled/vector"
)

type ShadowWalk struct {
	colorPlay                                   *color.ColorPlay
	yRot, xRot, zRot, yRotMax, xRotMax, zRotMax int
	buf                                         [][3]byte
	points                                      []vector.Vector3
}

func NewShadowWalk(leds []model.Led3D) *ShadowWalk {

	points := make([]vector.Vector3, len(leds))

	for i, l := range leds {
		if l.Inside {
			points[i] = l.Position
		} else {
			points[i] = vector.Vector3{X: 0, Y: 0, Z: 0}
		}
	}

	return &ShadowWalk{
		buf:       make([][3]byte, len(leds)),
		points:    points,
		xRotMax:   500,
		yRotMax:   2000,
		zRotMax:   50000,
		colorPlay: color.NewColorPlay(1024, 3),
	}
}

func (s *ShadowWalk) Next() [][3]byte {

	color := s.colorPlay.NextColor()

	x, y, z := 0., 1., 0.

	xPhase := 2 * math.Pi * float64(s.xRot) / float64(s.xRotMax)
	yPhase := 2 * math.Pi * float64(s.yRot) / float64(s.yRotMax)
	zPhase := 2 * math.Pi * float64(s.zRot) / float64(s.zRotMax)

	cosX, sinX := math.Cos(xPhase), math.Sin(xPhase)
	cosY, sinY := math.Cos(yPhase), math.Sin(yPhase)
	cosZ, sinZ := math.Cos(zPhase), math.Sin(zPhase)

	x, z = cosY*x+sinY*z, cosY*z+sinY*x
	y, z = cosX*y+sinX*z, cosX*z+sinX*y
	x, y = cosZ*x+sinZ*y, cosZ*y+sinZ*x

	v := vector.Vector3{X: x, Y: y, Z: z}

	for i, p := range s.points {

		d := v.SquaredDistance(p)
		if d < .5 {
			s.buf[i][0] = byte(float64(color[1]) * (1 - d*2))
			s.buf[i][1] = byte(float64(color[2]) * (1 - d*2))
			s.buf[i][2] = byte(float64(color[0]) * (1 - d*2))
		} else {
			s.buf[i] = [3]byte{0, 0, 0}
		}

	}

	s.yRot++
	s.yRot %= s.yRotMax

	s.xRot++
	s.xRot %= s.xRotMax

	s.zRot++
	s.zRot %= s.zRotMax

	return s.buf[:]
}
