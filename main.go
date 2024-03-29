package main

import (
	"bufio"
	"flag"
	"os"
	"path"
//"fmt"
	"post6.net/goled/ani"
	"post6.net/goled/ani/blend"
	"post6.net/goled/ani/cache"
	"post6.net/goled/ani/fire"
	"post6.net/goled/ani/fire2"
	"post6.net/goled/ani/five"
	"post6.net/goled/ani/fifteen"
	"post6.net/goled/ani/gameoflife"
	"post6.net/goled/ani/gradient"
	"post6.net/goled/ani/image"
	//"post6.net/goled/ani/onion"
	"post6.net/goled/ani/orbit"
	"post6.net/goled/ani/lorenz"
	"post6.net/goled/ani/radar"
	"post6.net/goled/ani/shadowplay"
	"post6.net/goled/ani/shadowwalk"
	"post6.net/goled/ani/snake"
	"post6.net/goled/ani/topo"
//	"post6.net/goled/ani/uniform"
	"post6.net/goled/ani/wobble"
	"post6.net/goled/ani/ring"
	"post6.net/goled/drivers"
	"post6.net/goled/led"
	"post6.net/goled/model"
	"post6.net/goled/model/poly/minipoly"
	"post6.net/goled/model/poly/polyhedrone"
	"post6.net/goled/model/poly/poly12"
	"post6.net/goled/model/poly/greatcircles"
	"post6.net/goled/model/poly/greatcircles2"
	"post6.net/goled/model/poly/aluball"
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
var mini, p12, gc, gc2, alu bool
var animations = []ani.Animation(nil)

func addAni(a ani.Animation) {
	animations = append(animations, a)
}

func init() {
	flag.Float64Var(&gamma, "gamma", 2.5, "used gamma correction setting")
	flag.Float64Var(&brightness, "brightness", 1., "used brighness setting")
	flag.IntVar(&fps, "fps", 100, "frames per second")
	flag.IntVar(&switchTime, "switchtime", 60, "seconds per animation")
//	flag.BoolVar(&ambient, "ambient", false, "don't load bright animations")
	flag.BoolVar(&mini, "mini", false, "use small polyhedron model")
	flag.BoolVar(&p12, "poly12", false, "use new polyhedron model")
	flag.BoolVar(&alu, "alu", false, "use aluball model")
	flag.BoolVar(&gc, "greatcircles0", false, "use old greatcircles model")
	flag.BoolVar(&gc2, "greatcircles", false, "use new greatcircles model")
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
	var inside, neighbours bool

	if p12 {
		ball = poly12.Ledball()
		inside = true
		neighbours = true
	} else if mini {
		ball = minipoly.Ledball()
		inside = false
		neighbours = true
	} else if gc2 {
		ball = greatcircles2.Ledball()
		inside = false
		neighbours = false
	} else if gc {
		ball = greatcircles.Ledball()
		inside = false
		neighbours = false
	} else if alu {
		ball = aluball.Ledball()
		inside = true
		neighbours = false
	} else {
		ball = polyhedrone.Ledball()
		inside = true
		neighbours = true
	}
	unitBall := ball.UnitScale()
	//smoothUnitBall := ball.UnitScale().Smooth()
	smooth := ball.Smooth()


	in := os.Stdin
	driver := drivers.GetLedDriver()
	strip := led.NewLedStrip(len(ball.Leds), ledOrder, driver.Bpp(), driver.MaxValue(), gamma, brightness)
	if driver.Bpp() == 16 {
		strip.SetCutoff(0x18, 0x18)
	}

	var frameBuffer []byte
	if p12 {
		frameBuffer = make([]byte, 930*strip.LedSize())
		strip.MapRange(0, 450, 0)
		strip.MapRange(450, 450, 480*strip.LedSize())
	} else {
		frameBuffer = make([]byte, len(ball.Leds)*strip.LedSize())
		strip.MapRange(0, len(ball.Leds), 0)
	}

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
	addAni(fire2.NewFire(smooth.Leds))
	addAni(wobble.NewWobble(unitBall.Leds, wobble.Outside))
	if neighbours {
		addAni(snake.NewSnake(ball))
	}
//	if !(gc || gc2) {
		addAni(cache.NewCachedAni(image.NewImageAni(smooth.Leds, earth, 0, 0, 0), len(smooth.Leds), 256))
		addAni(cache.NewCachedAni(image.NewImageAni(smooth.Leds, newImg, 0, 0, 0), len(smooth.Leds), 256))
//	}
	if inside {
		addAni(shadowwalk.NewShadowWalk(smooth.Leds))
		addAni(shadowplay.NewShadowPlay(ball.Leds, 512, 8))
		addAni(shadowplay.NewShadowPlay(ball.Leds, 512, 32))
	}
	addAni(topo.NewTopo(ball))
	addAni(orbit.NewOrbitAni(unitBall.Copy().Leds))
	if neighbours {
		addAni(gameoflife.NewGameOfLife(ball))
	}
	if !(gc || gc2) {
	addAni(gradient.NewGradient(smooth.Leds, gradient.Hard, gradient.Outside))
	}
	if inside {
		addAni(gradient.NewGradient(smooth.Leds, gradient.Hard, gradient.Inside))
	}
	if !(gc || gc2) {
	addAni(gradient.NewGradient(smooth.Leds, gradient.Binary, gradient.Outside))
	}
	addAni(gradient.NewGradient(smooth.Leds, gradient.Smooth, gradient.Outside))
	if !(gc || gc2) {
	addAni(gradient.NewGradient(smooth.Leds, gradient.Striped, gradient.Outside))
	}
	if inside {
		addAni(gradient.NewGradient(smooth.Leds, gradient.Striped, gradient.Inside))
	}

	addAni(lorenz.NewLorenzAni(unitBall.Leds))
	if !(gc || gc2) {
	addAni(gradient.NewSpiral(smooth.Leds, 1, 1, gradient.Hard, gradient.Outside))
	}
	addAni(gradient.NewSpiral(smooth.Leds, 1, 1, gradient.Smooth, gradient.Outside))
	if p12 {
		addAni(fifteen.NewFifteen(ball.Leds))
	} else {
		addAni(five.NewFive(ball.Leds))
	}
	addAni(radar.NewRadar(smooth.Leds))
	if p12 {
		addAni(fifteen.NewFifteenWave(smooth.Leds))
	} else if (gc || gc2) {
		addAni(ring.NewRingWave(ball))
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
			strip.LoadFrame(curFrame, frameBuffer)
			_, err := driver.Write(frameBuffer)
			if err != nil {
				panic("write error")
			}

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
					curFrame = make([][3]byte, len(ball.Leds))
					strip.LoadFrame(curFrame, frameBuffer)
					driver.Write(frameBuffer)
					os.Exit(0)
				}
			} else {
				events = nil
			}
		}
	}
}
