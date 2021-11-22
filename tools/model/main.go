package main

import (
    "fmt"
    "flag"
	"os"
	"post6.net/goled/model"
	"post6.net/goled/model/poly/minipoly"
	"post6.net/goled/model/poly/greatcircles"
	"post6.net/goled/model/poly/greatcircles2"
	"post6.net/goled/model/poly/polyhedrone"
	"post6.net/goled/model/poly/poly12"
	"post6.net/goled/model/poly/aluball"
//	"post6.net/goled/model/poly/icosidode"
//	"post6.net/goled/model/poly/miniball"
    "io"
)


func writeLed(out io.Writer, led model.Led3D) {
	boolStr := map[bool]string{ true: "true", false: "false" }
	fmt.Fprintf(out, "\t\t{\n"+
	                 "\t\t\t\"position\": [%f, %f, %f],\n"+
	                 "\t\t\t\"normal\": [%f, %f, %f],\n"+
	                 "\t\t\t\"inside\": %s\n"+
	                 "\t\t}",
	                 led.Position.X, led.Position.Y, led.Position.Z,
	                 led.Normal.X, led.Normal.Y, led.Normal.Z,
	                 boolStr[led.Inside])
}

func writeModel(out io.Writer, leds []model.Led3D) {

	fmt.Fprintf(out, "{\n\t\"leds\": [\n")
	for i,led := range leds {
		writeLed(out, led)
		if i < len(leds)-1 {
			fmt.Fprintf(out, ",\n")
		} else {
			fmt.Fprintf(out, "\n")
		}
	}
	fmt.Fprintf(out, "\t]\n}\n")
}

var modelName string
var unitScale bool

func main() {
	flag.StringVar(&modelName, "model", "polyhedrone", "model name")
	flag.BoolVar(&unitScale, "unitscale", false, "scale model to r=1")
	flag.Parse()

	var m *model.Model3D

	if modelName == "polyhedrone" {
		m = polyhedrone.Ledball()
	} else if modelName == "poly12" {
		m = poly12.Ledball()
	} else if modelName == "greatcircles" {
		m = greatcircles.Ledball()
	} else if modelName == "greatcircles2" {
		m = greatcircles2.Ledball()
	} else if modelName == "minipoly" {
		m = minipoly.Ledball()
	} else if modelName == "aluball" {
		m = aluball.Ledball()
/*
	} else if modelName == "icosidode" {
		m = icosidode.Ledball()
	} else if modelName == "miniball" {
		m = miniball.Ledball()
*/	}

	if (unitScale) {
		m = m.UnitScale()
	}

	writeModel(os.Stdout, m.Leds)
}
