package main

import (
	"flag"
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
	play := make([]*color.ColorPlay, 300)

	for i := range play {
		play[i] = color.NewColorPlay(256, 3)
	}

	t := time.Tick(10 * time.Millisecond)

	for {
		<-t
		for i := range frame {
			frame[i] = play[i].NextColor()
		}
		out.Write(strip.LoadFrame(frame))
	}

}
