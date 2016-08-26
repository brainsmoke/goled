package main

import (
	"flag"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"post6.net/goled/model"
	"post6.net/goled/model/poly/minipoly"
	"post6.net/goled/model/poly/polyhedrone"
	"post6.net/goled/model/poly/poly12"
	"strconv"
)

var gamma, brightness float64
var ledOrder led.LedOrder
var mini, p12 bool

func init() {
	flag.Float64Var(&gamma, "gamma", 1, "used gamma correction setting")
	flag.Float64Var(&brightness, "brightness", 1., "used brighness setting")
	ledOrder = led.RGB
	flag.Var(&ledOrder, "ledorder", "led order")
	flag.BoolVar(&mini, "mini", false, "use small polyhedron model")
	flag.BoolVar(&p12, "poly12", false, "use new polyhedron model")
}

func main() {

	flag.Parse()

	var ball *model.Model3D

	if p12 {
		ball = poly12.Ledball()
	} else if mini {
		ball = minipoly.Ledball()
	} else {
		ball = polyhedrone.Ledball()
	}

	args := flag.Args()
	r, g, b := 0, 0, 0
	var errr, errg, errb error
	if len(args) >= 3 {
		r, errr = strconv.Atoi(args[0])
		g, errg = strconv.Atoi(args[1])
		b, errb = strconv.Atoi(args[2])
		if errr != nil || errg != nil || errb != nil {
			r, g, b = 0, 0, 0
		}
	}

	rin, gin, bin := r, g, b
	if len(args) >= 6 {
		rin, errr = strconv.Atoi(args[3])
		gin, errg = strconv.Atoi(args[4])
		bin, errb = strconv.Atoi(args[5])
		if errr != nil || errg != nil || errb != nil {
			rin, gin, bin = 0, 0, 0
		}
	}

	strip := led.NewLedStrip(len(ball.Leds), ledOrder, gamma, brightness)

	out := drivers.LedDriver()

	frame := make( [][3]byte, len(ball.Leds) )

	for i := range frame {
		frame[i] = [3]byte{byte(r), byte(g), byte(b)}
		if ball.Leds[i].Inside {
			frame[i] = [3]byte{byte(rin), byte(gin), byte(bin)}
		}
	}

	out.Write(strip.LoadFrame(frame[:]))
}
