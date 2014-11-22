package main

import (
	"bufio"
	"flag"
	"os"
	"path"
	"post6.net/goled/ani"
	"post6.net/goled/ani/blend"
	"post6.net/goled/ani/cache"
	"post6.net/goled/ani/fire"
	"post6.net/goled/ani/five"
	"post6.net/goled/ani/gameoflife"
	"post6.net/goled/ani/gradient"
	"post6.net/goled/ani/image"
	"post6.net/goled/ani/onion"
	"post6.net/goled/ani/orbit"
	"post6.net/goled/ani/radar"
	"post6.net/goled/ani/shadowplay"
	"post6.net/goled/ani/shadowwalk"
	"post6.net/goled/ani/snake"
	"post6.net/goled/ani/topo"
	"post6.net/goled/ani/wobble"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"post6.net/goled/model"
	"time"
)

const (
	next = iota
	previous
	quit
)

func cmdHandler(file *os.File, events chan<- int) {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		switch scanner.Text() {

		case "next":
			events <- next
		case "previous":
			events <- previous
		case "quit":
			events <- quit
		}
	}

	close(events)
}

var gamma, brightness float64
var fps, switchTime, blendTime int
var ledOrder led.LedOrder
var ambient bool
var animations = []ani.Animation(nil)

func addAni(a ani.Animation) {
	animations = append(animations, a)
}

func init() {
	flag.Float64Var(&gamma, "gamma", 2.2, "used gamma correction setting")
	flag.Float64Var(&brightness, "brightness", 1., "used brighness setting")
	flag.IntVar(&fps, "fps", 80, "frames per second")
	flag.IntVar(&switchTime, "switchtime", 60, "seconds per animation")
	flag.BoolVar(&ambient, "ambient", false, "don't load bright animations")
	ledOrder = led.RGB
	flag.Var(&ledOrder, "ledorder", "led order")
}

func main() {

	blendTime = 3

	flag.Parse()

	strip := led.NewLedStrip(300, ledOrder, gamma, brightness)

	in := os.Stdin
	out := drivers.LedDriver()

	events := make(chan int)

	go cmdHandler(in, events)

	tick := time.NewTicker(time.Second / time.Duration(fps))

	nextAni := time.NewTicker(time.Second * time.Duration(switchTime))

	baseDir := path.Dir(os.Args[0])
	earth, _ := os.Open(baseDir + "/earth.png")

	addAni(wobble.NewWobble(model.LedballSmooth(), wobble.Inside))
	addAni(shadowwalk.NewShadowWalk(model.LedballSmooth()))
	addAni(snake.NewSnake())
	addAni(fire.NewInnerFire(model.LedballSmooth()))
	addAni(fire.NewFire(model.LedballSmooth()))
	if !ambient {
		addAni(wobble.NewWobble(model.LedballSmooth(), wobble.Outside))
	}
	addAni(cache.NewCachedAni(image.NewImageAni(model.LedballSmooth(), earth, 0, 0, 0), 300, 256))
	addAni(shadowplay.NewShadowPlay(512, 3))
	if !ambient {
		addAni(topo.NewTopo())
	}
	addAni(orbit.NewOrbitAni(model.Ledball()))
	if !ambient {
		addAni(gradient.NewGradient(model.LedballSmooth(), gradient.Hard))
	}
	addAni(gameoflife.NewGameOfLife())
	if !ambient {
		addAni(gradient.NewGradient(model.LedballSmooth(), gradient.Smooth))
	}
	addAni(shadowplay.NewShadowPlay(1024, 12))
	if !ambient {
		addAni(radar.NewRadar(model.LedballSmooth()))
	}
	if !ambient {
		addAni(onion.NewOnion(model.LedballSmooth()))
	}
	addAni(five.NewFive())
	addAni(five.NewUniformInside())

	current, last := 0, -1
	blendIter := 0
	blendFrames := fps * blendTime

	var curFrame, lastFrame, tmpFrame [][3]byte

	tmpFrame = make([][3]byte, 300)

	for {
		select {

		case <-tick.C:
			curFrame = animations[current].Next()
			if blendIter > 0 {

				if last >= 0 {
					lastFrame = animations[last].Next()
				}
				blend.Blend(tmpFrame, curFrame, lastFrame, float64(blendIter)/float64(blendFrames))
				curFrame = tmpFrame
				blendIter--
			}
			out.Write(strip.LoadFrame(curFrame))

		case <-nextAni.C:
			last = current
			current = (current + 1) % len(animations)
			blendIter = blendFrames

		case e, more := <-events:
			switch e {
			case next:
				current = (current + 1) % len(animations)
				lastFrame, last = curFrame, -1
				blendIter = blendFrames
				nextAni.Stop()
				nextAni = time.NewTicker(time.Second * time.Duration(switchTime))
			case previous:
				current = (current + len(animations) - 1) % len(animations)
				lastFrame, last = curFrame, -1
				blendIter = blendFrames
				nextAni.Stop()
				nextAni = time.NewTicker(time.Second * time.Duration(switchTime))
			case quit:
				os.Exit(0)
			}
			if !more {
				events = nil
			}
		}
	}
}
