package vector


type Matrix3x3 struct {
	v [3]Vector3
}

func RotationMatrix(eye, center, north Vector3) Matrix3x3 {

	f := eye.Sub(center).Normalize()
	f_neg := center.Sub(eye).Normalize()
	up := north.Normalize()
	s := f.CrossProduct(up).Normalize()
	u := s.CrossProduct(f)
	return Matrix3x3{ [3]Vector3{s, u, f_neg} }

}

func (m Matrix3x3) Mul(v Vector3) Vector3 {
	return Vector3{ X: m.v[0].ScalarProduct(v), Y: m.v[1].ScalarProduct(v), Z: m.v[2].ScalarProduct(v) }
}
