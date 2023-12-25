package main

import (
    "fmt"
	"flag"
	"post6.net/goled/model"
	"post6.net/goled/projection"
	"post6.net/goled/vector"
	"post6.net/goled/model/poly/greatcircles"
	"post6.net/goled/model/poly/greatcircles2"
	"post6.net/goled/model/poly/polyhedrone"
	"post6.net/goled/model/poly/poly12"
	"post6.net/goled/model/poly/minipoly"
	"post6.net/goled/model/poly/aluball"
)

func writeMap(vmap []int, width, height int) {

	for y:=0; y<height; y++ {
		fmt.Printf("\t\t")
		for x:=0; x<width; x++ {
			fmt.Printf("%4d", vmap[y*width+x])
			if x+1 != width || y+1 != height {
				fmt.Printf(",")
			}
		}
		fmt.Printf("\n")
	}
}

func writeFloatarray(arr []float64) {

	fmt.Printf("\t\t")
    x := ""
	for _,v := range arr {
		fmt.Printf(x)
		fmt.Printf("%f", v)
		x = ", "
	}
	fmt.Printf("\n")
}

var modelName string
var width, height int

func main() {
	flag.IntVar(&width, "width", 42, "width")
	flag.IntVar(&height, "height", 20, "height")
	flag.StringVar(&modelName, "model", "polyhedrone", "model name")
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
/*	} else if modelName == "icosidode" {
		m = icosidode.Ledball()
	} else if modelName == "miniball" {
		m = miniball.Ledball()
*/	}

	m = m.UnitScale()

    points := make([]vector.Vector3, len(m.Leds))

    for i, led := range m.Leds {
		points[i] = led.Position
    }

	vmap, sinY := projection.Voronoi(width, height, points)


	fmt.Printf("{\n")
	fmt.Printf("\t\"width\": %d,\n", width)
	fmt.Printf("\t\"height\": %d,\n", height)
	fmt.Printf("\t\"voronoi\": [\n")
	writeMap(vmap, width, height)
	fmt.Printf("\t],\n")
	fmt.Printf("\t\"weighting\": [\n")
	writeFloatarray(sinY)
	fmt.Printf("\t]\n")
	fmt.Printf("}\n")
}
