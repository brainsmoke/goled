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

	ball := polyhedrone.Ledball().Smooth()

	out := drivers.GetLedDriver()
	strip := led.NewLedStrip(len(ball.Leds), ledOrder, out.Bpp(), out.MaxValue(), gamma, brightness)
	strip.MapRange(0, len(ball.Leds), 0)
	frameBuffer := make([]byte, len(ball.Leds)*strip.LedSize())

	animation := fire.NewFire(ball.Leds)

	t := time.Tick(time.Second / time.Duration(fps))

	for {
		<-t
		strip.LoadFrame(animation.Next(), frameBuffer)
		out.Write(frameBuffer)
	}
}
