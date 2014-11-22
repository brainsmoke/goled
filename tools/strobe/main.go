package main

import (
	"flag"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"strconv"
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

	strip := led.NewLedStrip(300, ledOrder, gamma, brightness)

	out := drivers.LedDriver()

	frame1 := make([][3]byte, 300)
	frame2 := make([][3]byte, 300)

	for i := range frame1 {
		frame1[i] = [3]byte{byte(r), byte(g), byte(b)}
		frame2[i] = [3]byte{byte(rin), byte(gin), byte(bin)}
	}

	t := time.Tick(40 * time.Millisecond)

	for {
		<-t
		out.Write(strip.LoadFrame(frame1))
		frame1, frame2 = frame2, frame1
	}

}
