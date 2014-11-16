package model

import (
	. "post6.net/goled/vector"
)

type Led3D struct {
	Position, Normal Vector3
	Face             int
	Inside           bool
}
