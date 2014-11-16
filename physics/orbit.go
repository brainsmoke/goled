package physics

import . "post6.net/goled/vector"

const G = 6.674e-11

type Object struct {
	m       float64
	p, v, a Vector3
	fixed   bool
}

func (o *Object) P() Vector3 {
	return o.p
}

func (o *Object) V() Vector3 {
	return o.v
}

func (o *Object) M() float64 {
	return o.m
}

func NewFixedObject(m float64, p Vector3) *Object {

	return &Object{m: m, p: p, fixed: true}
}

func NewObject(m float64, p, v Vector3) *Object {

	return &Object{m: m, p: p, v: v, fixed: false}
}

func Update(dt float64, objects ...*Object) {

	for _, o1 := range objects {
		if !o1.fixed {

			o1.a = Vector3{0, 0, 0}
			for _, o2 := range objects {
				if o1 != o2 {
					relP := o2.P().Sub(o1.P())
					d := relP.Magnitude()
					//F := relP.Normalize().Mul(G*(o1.m*o2.m)/(d*d))
					//a := F.Mul(1/o1.m)
					//o1.a = o1.a.Add( a )

					o1.a = o1.a.Add(relP.Mul(G * o2.m / (d * d * d)))
				}
			}
		}
	}
	for i := range objects {
		objects[i].v = objects[i].v.Add(objects[i].a.Mul(dt))
		objects[i].p = objects[i].p.Add(objects[i].v.Mul(dt))
	}
}
