package main

import (
	"flag"
	"post6.net/goled/polyhedron"
	"post6.net/goled/scad"
	"post6.net/goled/vector"
	"os"
)

var name string

func init() {
	flag.StringVar(&name, "model", "", "conway name")
}

func main() {

	flag.Parse()

	modelFunc, ok := polyhedron.ConwayModels[name]

	if ! ok {
		panic("meh")
	}

	model := modelFunc()

    north := model.Faces[0].Normal
    center := vector.Vector3{ 0, 0, 0 }
    eye := north.CrossProduct(model.Points[model.Faces[0].Polygon[0]])

    model.Rotate(eye, center, north)

	scad.WriteOpenScad(os.Stdout, model)
}
