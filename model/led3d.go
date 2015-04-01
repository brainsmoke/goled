package model

import (
	. "post6.net/goled/vector"
)

type Led3D struct {
	Position, Normal Vector3
	Face             int
	Inside           bool
}

type Led2D struct {
	X, Y   float64
	Inside bool
}
