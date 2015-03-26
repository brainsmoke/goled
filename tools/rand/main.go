package main

import (
	"flag"
	"math/rand"
	"post6.net/goled/color"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"time"
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

	frame := make([][3]byte, 300)

	t := time.Tick(80 * time.Millisecond)

	for {
		<-t
		for i := range frame {
			frame[i] = color.HSIToRGB(rand.Float64(), 1, .5+rand.Float64()/2)
		}
		out.Write(strip.LoadFrame(frame))
	}

}
