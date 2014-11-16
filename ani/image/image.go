package image

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"post6.net/goled/model"
	"post6.net/goled/projection"
	"post6.net/goled/util/clip"
	"post6.net/goled/vector"
)

type ImageAni struct {
	buf         [][3]byte
	accum       [][3]float64
	max         []float64
	sinY        []float64
	vmap        []int
	img         *image.RGBA
	rot         int
	insideColor [3]byte
}

func NewImageAni(leds []model.Led3D, r io.Reader, inr, ing, inb byte) (a *ImageAni) {

	a = new(ImageAni)

	srcImg, _, err := image.Decode(r)

	if err != nil {
		panic("meh")
	}

	a.img = image.NewRGBA(srcImg.Bounds())
	draw.Draw(a.img, a.img.Bounds(), srcImg, image.ZP, draw.Src)

	width := a.img.Rect.Max.X - a.img.Rect.Min.X
	height := a.img.Rect.Max.Y - a.img.Rect.Min.Y

	points := make([]vector.Vector3, len(leds))

	for i, led := range leds {

		if led.Inside {
			points[i] = vector.Vector3{X: 0, Y: 0, Z: 0}
		} else {
			points[i] = led.Position
		}
	}

	a.vmap, a.sinY = projection.Voronoi(width, height, points)

	a.buf = make([][3]byte, len(leds))
	a.accum = make([][3]float64, len(leds))
	a.max = make([]float64, len(leds))

	a.rot = 0

	a.insideColor = [3]byte{inr, ing, inb}

	return a
}

func (a *ImageAni) Next() [][3]byte {

	width := a.img.Rect.Max.X - a.img.Rect.Min.X
	height := a.img.Rect.Max.Y - a.img.Rect.Min.Y
	stride := a.img.Stride
	pix := a.img.Pix
	sinY := a.sinY

	for i := range a.buf {

		a.accum[i] = [3]float64{0, 0, 0}
		a.max[i] = 0
	}

	i := 0
	for ty := 0; ty < height; ty++ {
		pixel := ty * stride
		density := sinY[ty]
		i += width - a.rot
		for tx := 0; tx < width; tx++ {

			if tx == a.rot {
				i -= width
			}

			led := a.vmap[i]
			a.max[led] += density
			a.accum[led][0] += density * float64(pix[pixel])
			a.accum[led][1] += density * float64(pix[pixel+1])
			a.accum[led][2] += density * float64(pix[pixel+2])

			i += 1
			pixel += 4
		}
		i += a.rot
	}

	a.rot = (a.rot + width - 1) % width

	for i := range a.buf {
		if a.max[i] == 0 {
			a.buf[i] = a.insideColor
		} else {
			a.buf[i][0] = clip.FloatToByte(a.accum[i][0] / a.max[i])
			a.buf[i][1] = clip.FloatToByte(a.accum[i][1] / a.max[i])
			a.buf[i][2] = clip.FloatToByte(a.accum[i][2] / a.max[i])
		}
	}
	return a.buf[:]
}
