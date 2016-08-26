package main

import (
	"flag"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"strconv"
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

	args := flag.Args()
	n := 0
	var err error
	if len(args) >= 1 {
		n, err = strconv.Atoi(args[0])
		if err != nil {
			n = -1
		}
	}

	strip := led.NewLedStrip(900, ledOrder, gamma, brightness)

	out := drivers.LedDriver()

	var frame [900][3]byte

	frame[n] = [3]byte{255, 0, 0}

	out.Write(strip.LoadFrame(frame[:]))
}
