package main

import (
	"flag"
	"time"
	"post6.net/goled/ani/fire"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"post6.net/goled/model/poly/polyhedrone"
)

var gamma, brightness float64
var ledOrder led.LedOrder
var fps int

func init() {
	flag.Float64Var(&gamma, "gamma", 1, "used gamma correction setting")
	flag.Float64Var(&brightness, "brightness", 1., "used brighness setting")
	flag.IntVar(&fps, "fps", 80, "frames per second")
	ledOrder = led.RGB
	flag.Var(&ledOrder, "ledorder", "led order")
}

func main() {

	flag.Parse()

	strip := led.NewLedStrip(300, ledOrder, gamma, brightness)
	out := drivers.LedDriver()

	animation := fire.NewFire(polyhedrone.Ledball().Smooth().Leds)

	t := time.Tick(time.Second / time.Duration(fps))

	for {
		<-t
		out.Write(strip.LoadFrame(animation.Next()))
	}
}
