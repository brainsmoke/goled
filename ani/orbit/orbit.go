package orbit

import (
	"math"
	"post6.net/goled/model"
	"post6.net/goled/physics"
	"post6.net/goled/util/clip"
	. "post6.net/goled/vector"
)

type OrbitAni struct {
	ball    *physics.Object
	objects [3]*physics.Object
	lights  [3][3]float64

	leds []model.Led3D
	buf  [][3]byte
}

func NewOrbitAni(leds []model.Led3D) (a *OrbitAni) {

	a = new(OrbitAni)

	a.leds = leds
	a.buf = make([][3]byte, len(leds))

	rBall := 6371.0 * 1000
	mBall := 5.97219e24
	pBall := Vector3{0, 0, 0}

	for i := range a.leds {
		a.leds[i].Position = a.leds[i].Position.Mul(rBall)
	}

	a.ball = physics.NewFixedObject(mBall, pBall)

	var angle, m, d, speed, ill, r, g, b float64
	var p, v Vector3

	// Object
	angle = .3
	m = 10000.
	p = Vector3{rBall + 300*1000., 0., 0.}
	d = p.Magnitude()
	speed = math.Sqrt(physics.G*(mBall+m)/d) * 1.1
	v = Vector3{0., math.Sin(angle) * speed, math.Cos(angle) * speed}

	a.objects[0] = physics.NewObject(m, p, v)

	ill = 4000 * 1000. // distance for normalized illumination
	r, g, b = 0, 0, 255
	a.lights[0] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	// Object
	angle = 1.3
	m = 10.
	p = Vector3{0., 0., rBall + 3000*1000.}
	d = p.Magnitude()
	speed = math.Sqrt(physics.G*(mBall+m)/d) * 1.15
	v = Vector3{math.Sin(angle) * speed, math.Cos(angle) * speed, 0.}

	a.objects[1] = physics.NewObject(m, p, v)

	ill = 4000 * 1000. // distance for normalized illumination
	r, g, b = 255, 255, 64
	a.lights[1] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	// Object
	angle = -.7
	m = 10000.
	p = Vector3{0., rBall + 10000*1000., 0.}
	d = p.Magnitude()
	speed = math.Sqrt(physics.G*(mBall+m)/d) * 1
	v = Vector3{math.Sin(angle) * speed, 0, math.Cos(angle) * speed}

	a.objects[2] = physics.NewObject(m, p, v)

	ill = 10000 * 1000. // distance for normalized illumination
	r, g, b = 255, 32, 127
	a.lights[2] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	return a
}

func shader(pos, normal, lightPos Vector3, light [3]float64) (r, g, b float64) {

	relP := lightPos.Sub(pos)
	sprod := relP.ScalarProduct(normal) // distance * 'angle'

	if sprod < 0 {
		return 0, 0, 0
	}
	d := relP.Magnitude()
	intensity := sprod / (d * d * d) //

	return light[0] * intensity, light[1] * intensity, light[2] * intensity
}

func (a *OrbitAni) Next() [][3]byte {

	for i := 0; i < 10; i++ {

		for _, o := range a.objects {
			physics.Update(1, a.ball, o)
		}
	}

	for i, led := range a.leds {

		r, g, b := 0., 0., 0.
		for j, o := range a.objects {
			rj, gj, bj := shader(led.Position, led.Normal, o.P(), a.lights[j])
			r, g, b = r+rj, gj+g, bj+b
		}
		a.buf[i] = [3]byte{clip.FloatToByte(r), clip.FloatToByte(g), clip.FloatToByte(b)}
	}

	return a.buf[:]
}
