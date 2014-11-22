package main

import (
	"flag"
	"post6.net/goled/ani/uniform"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"time"
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
	animation := uniform.NewUniformInside()

	t := time.Tick(time.Second / time.Duration(fps))

	for {
		<-t
		out.Write(strip.LoadFrame(animation.Next()))
	}
}