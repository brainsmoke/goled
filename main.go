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
	"post6.net/goled/ani/fifteen"
	"post6.net/goled/ani/gameoflife"
	"post6.net/goled/ani/gradient"
	"post6.net/goled/ani/image"
	//"post6.net/goled/ani/onion"
	"post6.net/goled/ani/orbit"
	"post6.net/goled/ani/radar"
	"post6.net/goled/ani/shadowplay"
	"post6.net/goled/ani/shadowwalk"
	"post6.net/goled/ani/snake"
	"post6.net/goled/ani/topo"
//	"post6.net/goled/ani/uniform"
	"post6.net/goled/ani/wobble"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"post6.net/goled/model"
	"post6.net/goled/model/poly/minipoly"
	"post6.net/goled/model/poly/polyhedrone"
	"post6.net/goled/model/poly/poly12"
	"time"
)

const (
	nop = iota
	next
	previous
	quit
	lock
	unlock
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
		case "lock":
			events <- lock
		case "unlock":
			events <- unlock
		}
	}

	close(events)
}

var gamma, brightness float64
var fps, switchTime, blendTime int
var ledOrder led.LedOrder
//var ambient bool
var mini, p12 bool
var animations = []ani.Animation(nil)

func addAni(a ani.Animation) {
	animations = append(animations, a)
}

func init() {
	flag.Float64Var(&gamma, "gamma", 2.2, "used gamma correction setting")
	flag.Float64Var(&brightness, "brightness", 1., "used brighness setting")
	flag.IntVar(&fps, "fps", 80, "frames per second")
	flag.IntVar(&switchTime, "switchtime", 60, "seconds per animation")
//	flag.BoolVar(&ambient, "ambient", false, "don't load bright animations")
	flag.BoolVar(&mini, "mini", false, "use small polyhedron model")
	flag.BoolVar(&p12, "poly12", false, "use new polyhedron model")
	ledOrder = led.RGB
	flag.Var(&ledOrder, "ledorder", "led order")
}

func nextTicker() *time.Ticker {
	if switchTime == 0 {
		t := time.NewTicker(time.Second * time.Duration(1000000))
		t.Stop()
		return t
	} else {
		return time.NewTicker(time.Second * time.Duration(switchTime))
	}
}

func main() {

	blendTime = 3

	flag.Parse()

	var ball *model.Model3D
	var inside bool

	if p12 {
		ball = poly12.Ledball()
		inside = true
	} else if mini {
		ball = minipoly.Ledball()
		inside = false
	} else {
		ball = polyhedrone.Ledball()
		inside = true
	}
	unitBall := ball.UnitScale()
	//smoothUnitBall := ball.UnitScale().Smooth()
	smooth := ball.Smooth()

	strip := led.NewLedStrip(len(ball.Leds), ledOrder, gamma, brightness)

	in := os.Stdin
	out := drivers.LedDriver()

	events := make(chan int)

	go cmdHandler(in, events)

	tick := time.NewTicker(time.Second / time.Duration(fps))

	nextAni := nextTicker()

	baseDir := path.Dir(os.Args[0])
	earth, _ := os.Open(baseDir + "/earth.png")
	newImg, _ := os.Open(baseDir + "/newimg.png")

	if inside {
		addAni(wobble.NewWobble(smooth.Leds, wobble.Inside))
	}
	addAni(fire.NewFire(smooth.Leds))
	addAni(wobble.NewWobble(unitBall.Leds, wobble.Outside))
	addAni(snake.NewSnake(ball))
	addAni(cache.NewCachedAni(image.NewImageAni(smooth.Leds, earth, 0, 0, 0), len(smooth.Leds), 256))
	addAni(cache.NewCachedAni(image.NewImageAni(smooth.Leds, newImg, 0, 0, 0), len(smooth.Leds), 256))
	if inside {
		addAni(shadowwalk.NewShadowWalk(smooth.Leds))
	}
	addAni(shadowplay.NewShadowPlay(ball.Leds, 512, 8))
	addAni(shadowplay.NewShadowPlay(ball.Leds, 512, 32))
	addAni(topo.NewTopo(ball))
	addAni(orbit.NewOrbitAni(unitBall.Leds))
	addAni(gameoflife.NewGameOfLife(ball))
	addAni(gradient.NewGradient(smooth.Leds, gradient.Hard, gradient.Outside))
	addAni(gradient.NewGradient(smooth.Leds, gradient.Hard, gradient.Inside))
	addAni(gradient.NewGradient(smooth.Leds, gradient.Binary, gradient.Outside))
	addAni(gradient.NewGradient(smooth.Leds, gradient.Smooth, gradient.Outside))
	addAni(gradient.NewGradient(smooth.Leds, gradient.Striped, gradient.Outside))
	addAni(gradient.NewGradient(smooth.Leds, gradient.Striped, gradient.Inside))

	addAni(gradient.NewSpiral(smooth.Leds, gradient.Hard, gradient.Outside))
	addAni(gradient.NewSpiral(smooth.Leds, gradient.Binary, gradient.Outside))
	addAni(gradient.NewSpiral(smooth.Leds, gradient.Smooth, gradient.Outside))
	addAni(gradient.NewSpiral(smooth.Leds, gradient.Striped, gradient.Outside))
	addAni(gradient.NewSpiral(smooth.Leds, gradient.Striped, gradient.Inside))
	if p12 {
		addAni(fifteen.NewFifteen(ball.Leds))
	} else {
		addAni(five.NewFive(ball.Leds))
	}
	addAni(radar.NewRadar(smooth.Leds))
	if p12 {
		addAni(fifteen.NewFifteenWave(smooth.Leds))
	} else {
		addAni(five.NewFiveWave(smooth.Leds))
	}

	current, last := 0, -1
	blendIter := 0
	blendFrames := fps * blendTime

	var curFrame, lastFrame, tmpFrame [][3]byte

	tmpFrame = make([][3]byte, len(ball.Leds))

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
			if more {
				switch e {
				case next:
					current = (current + 1) % len(animations)
					lastFrame, last = curFrame, -1
					blendIter = blendFrames
					if nextAni.C != nil {
						nextAni.Stop()
						nextAni = nextTicker()

					}
				case previous:
					current = (current + len(animations) - 1) % len(animations)
					lastFrame, last = curFrame, -1
					blendIter = blendFrames
					if nextAni.C != nil {
						nextAni.Stop()
						nextAni = nextTicker()
					}
				case lock:
					if nextAni.C != nil {
						nextAni.Stop()
						nextAni.C = nil
					}
				case unlock:
					if nextAni.C == nil {
						nextAni = nextTicker()
					}
				case quit:
					os.Exit(0)
				}
			} else {
				events = nil
			}
		}
	}
}
