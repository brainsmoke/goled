package vector

import (
	"fmt"
	"math"
)

type Vector3 struct {
	X, Y, Z float64
}

func (v1 Vector3) SquaredDistance(v2 Vector3) float64 {

	var dx, dy, dz = v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z

	return dx*dx + dy*dy + dz*dz
}

func (v1 Vector3) Distance(v2 Vector3) float64 {

	var dx, dy, dz = v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z

	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (v Vector3) Magnitude() float64 {

	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) Normalize() Vector3 {

	var d = v.Magnitude()
	return Vector3{v.X / d, v.Y / d, v.Z / d}
}

func (v1 Vector3) Add(v2 Vector3) Vector3 {

	return Vector3{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

func (v1 Vector3) Sub(v2 Vector3) Vector3 {

	return Vector3{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

func (v Vector3) Mul(s float64) Vector3 {

	return Vector3{v.X * s, v.Y * s, v.Z * s}
}

func (v1 Vector3) CrossProduct(v2 Vector3) Vector3 {

	return Vector3{v1.Y*v2.Z - v2.Y*v1.Z,
		v1.Z*v2.X - v2.Z*v1.X,
		v1.X*v2.Y - v2.X*v1.Y,
	}
}

func (v1 Vector3) ScalarProduct(v2 Vector3) float64 {

	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vector3) Interpolate(v2 Vector3, fraction float64) Vector3 {

	return Vector3{v2.X*fraction + v1.X*(1.0-fraction),
		v2.Y*fraction + v1.Y*(1.0-fraction),
		v2.Z*fraction + v1.Z*(1.0-fraction),
	}
}

func (v Vector3) String() string {

	return fmt.Sprintf("(%.5f, %.5f, %.5f)", v.X, v.Y, v.Z)
}
