package main

import (
	"flag"
	"post6.net/goled/ani/fire"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"post6.net/goled/model"
)

var gamma, brightness float64
var ledOrder led.LedOrder

func init() {
	flag.Float64Var(&gamma, "gamma", 1, "used gamma correction setting")
	flag.Float64Var(&brightness, "brightness", 1., "used brighness setting")
	ledOrder = led.RGB
	flag.Var(&ledOrder, "ledorder", "led order")
}

func main() {

	flag.Parse()

	strip := led.NewLedStrip(300, ledOrder, gamma, brightness)
	out := drivers.LedDriver()
	animation := fire.NewFire(model.LedballSmooth())

	for {
		out.Write(strip.LoadFrame(animation.Next()))
	}
}
