package main

import (
	"flag"
	"post6.net/goled/polyhedron"
	"post6.net/goled/stl"
	"os"
)

var name string

func init() {
	flag.StringVar(&name, "model", "", "conway name")
}

func main() {

	flag.Parse()

	model, ok := polyhedron.ConwayModels[name]

	if ! ok {
		panic("meh")
	}

	stl.WriteStl(os.Stdout, name, model())

}
