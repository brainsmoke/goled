package cache

import (
	"post6.net/goled/ani"
)

type CachedAni struct {
	frames                [][3]byte
	frameSize, frameCount int
	rot                   int
}

func NewCachedAni(orig ani.Animation, frameSize, frameCount int) (a *CachedAni) {

	a = new(CachedAni)
	a.frames = make([][3]byte, frameSize*frameCount)

	a.frameSize, a.frameCount = frameSize, frameCount

	a.rot = 0

	for i := 0; i < frameCount; i++ {
		copy(a.frames[i*frameSize:(i+1)*frameSize], orig.Next())
	}

	return a
}

func (a *CachedAni) Next() [][3]byte {

	cur := a.rot * a.frameSize
	a.rot = (a.rot + 1) % a.frameCount
	return a.frames[cur : cur+a.frameSize]
}
