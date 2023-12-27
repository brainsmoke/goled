package main

import (
    "fmt"
	"flag"
	"post6.net/goled/model"
	"post6.net/goled/model/poly/greatcircles"
	"post6.net/goled/model/poly/greatcircles2"
	"post6.net/goled/model/poly/polyhedrone"
	"post6.net/goled/model/poly/poly12"
	"post6.net/goled/model/poly/minipoly"
	"post6.net/goled/model/poly/aluball"
)


var modelName string

func main() {
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

	fmt.Print("{\n\t\"neighbours\":\n\t[")
	s := "\n"

	for _,list := range m.Neighbours {
		fmt.Print(s)
		s = ",\n"
		fmt.Print("\t\t[ ")
		for i, v := range list {
			if i != 0 {
				fmt.Print(",")
			}
			fmt.Printf("%3d", v)
		}
		fmt.Print(" ]")
	}
	fmt.Print("\n\t]\n}\n")
}
