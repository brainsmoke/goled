package main

import (
    "fmt"
	"os"
	"post6.net/goled/model"
	"post6.net/goled/model/poly/polyhedrone"
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

func main() {

	writeModel(os.Stdout, polyhedrone.Ledball().Leds)
}
