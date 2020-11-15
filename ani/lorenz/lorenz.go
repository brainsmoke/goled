package lorenz

import (
	"post6.net/goled/model"
	"post6.net/goled/util/clip"
	"post6.net/goled/vector"
)

type LorenzAni struct {
	attractors [6]vector.Vector3
	rotation [6]vector.Matrix3x3
	lights  [6][3]float64
	leds []model.Led3D
	buf  [][3]byte
}

func NewLorenzAni(leds []model.Led3D) (a *LorenzAni) {

	a = new(LorenzAni)

	a.leds = leds
	a.buf = make([][3]byte, len(leds))

	var ill, r, g, b float64

	// Object
	ill = 1.2 // distance for normalized illumination
	r, g, b = 0, 0, 255
	a.lights[0] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	ill = 1.2 // distance for normalized illumination
	r, g, b = 255, 255, 64
	a.lights[1] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	ill = 1.2 // distance for normalized illumination
	r, g, b = 255, 32, 127
	a.lights[2] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	// Object
	ill = 1.2 // distance for normalized illumination
	r, g, b = 0, 255, 25
	a.lights[3] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	ill = 1.2 // distance for normalized illumination
	r, g, b = 255, 55, 34
	a.lights[4] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	ill = 1.2 // distance for normalized illumination
	r, g, b = 55, 32, 127
	a.lights[5] = [3]float64{r * ill * ill, g * ill * ill, b * ill * ill}

	a.attractors = [6]vector.Vector3{ {X:1, Y:1, Z:.5}, {X:1, Y:.5, Z: 1}, {X:.5, Y:1, Z:1},
	                                  {X:1, Y:2, Z:.6}, {X:1, Y:.6, Z: 2}, {X:.6, Y:2, Z:1} }

	center := vector.Vector3{X:0, Y:0, Z:0}
	north := vector.Vector3{X:0, Y:1, Z:0}
	eye := vector.Vector3{X:0, Y:0, Z:1}
	identity := vector.RotationMatrix( eye, center, north )

	somerotation := vector.RotationMatrix( north, center, eye )

	someotherrotation := vector.RotationMatrix( north.Mul(-1), center, eye.Mul(-1) )

	a.rotation = [6]vector.Matrix3x3{ identity, identity, identity,
	                                   somerotation, somerotation, someotherrotation }

	return a
}

func shader(pos, normal, lightPos vector.Vector3, light [3]float64) (r, g, b float64) {

	relP := lightPos.Sub(pos)
	sprod := relP.ScalarProduct(normal) // distance * 'angle'

	if sprod < 0 {
		return 0, 0, 0
	}
	d := relP.Magnitude()
	intensity := sprod / (d * d * d) //

	return light[0] * intensity, light[1] * intensity, light[2] * intensity
}

func (a *LorenzAni) Next() [][3]byte {

	sigma, rho, beta := 10., 28., 8/3.
	dt := 1/5000.
	for i := 0; i < 10; i++ {

		for j, o := range a.attractors {
			x, y, z := o.X, o.Y, o.Z
			dx, dy, dz := sigma*(y-x)*dt, (x*(rho-z)-y)*dt, (x*y-beta*z)*dt
			a.attractors[j] = vector.Vector3{X:x+dx, Y:y+dy, Z:z+dz}
		}
	}

	for i, led := range a.leds {

		r, g, b := 0., 0., 0.
		for j, o := range a.attractors {
			rj, gj, bj := shader(led.Position, led.Normal, a.rotation[j].Mul(o.Mul(.11)), a.lights[j])
//fmt.Println(o,rj, gj, bj, led.Normal,led.Position, a.lights[j])
			r, g, b = r+rj, gj+g, bj+b
		}
		a.buf[i] = [3]byte{clip.FloatToByte(r), clip.FloatToByte(g), clip.FloatToByte(b)}
	}

	return a.buf[:]
}
