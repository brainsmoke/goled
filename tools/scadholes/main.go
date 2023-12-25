package main

import (
	"post6.net/goled/polyhedron"
	"post6.net/goled/vector"
	"os"
	"fmt"
	"flag"
)

const (
	TopLeft     = 0
	BottomLeft  = 1
	BottomRight = 2
	TopRight    = 3
)

const innerRadius = 1

var traversal = polyhedron.RemapRoute{

TopRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight,
}

var name string

func init() {
	flag.StringVar(&name, "model", "oD", "conway name")
}

func main() {

	flag.Parse()

	model, ok := polyhedron.ConwayModels[name]

	if !ok {
		panic("meh.")
	}

	solid := model()//polyhedron.RemapSolid(polyhedron.DeltoidalHexecontahedron(), 0, traversal)

    north := solid.Points[solid.Faces[1].Polygon[2]].Normalize()
    center := vector.Vector3{ 0, 0, 0 }
    eye := north.CrossProduct(solid.Points[solid.Faces[1].Polygon[3]]).Normalize()

    solid.Rotate(eye, center, north)

	for _,f := range(solid.Faces) {

		top := solid.Points[f.Polygon[0]]

		newZ := f.Normal
		vY := top.Sub(f.Center).Normalize()
		newX := vY.CrossProduct(newZ).Normalize()
		newY := newZ.CrossProduct(newX).Normalize()

		fmt.Fprintf(os.Stdout, "multmatrix (m=[[%e,%e,%e,%e],[%e,%e,%e,%e],[%e,%e,%e,%e],[0, 0, 0, 1]]) {hole();};\n",
			newX.X, newY.X, newZ.X, 0.,
			newX.Z, newY.Z, newZ.Z, 0.,
			-newX.Y, -newY.Y, -newZ.Y, 0.,
			)
	}
}
