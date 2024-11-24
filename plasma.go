package main

import (
	"math"
	"math/rand"
	"time"
)

type Blob struct {
	x    float64
	y    float64
	dx   float64
	dy   float64
	size float64
	r    float64
	g    float64
	b    float64
}

type Plasma struct {
	rgb   *RGBComponent
	blobs []Blob
}

func NewPlasma(rgb *RGBComponent, dirty Dirty) *Plasma {
	p := &Plasma{rgb: rgb, blobs: make([]Blob, 60)}

	for t := 0; t < len(p.blobs); t++ {
		p.blobs[t] = Blob{
			x:    rand.Float64() * float64(rgb.width),
			y:    rand.Float64() * float64(rgb.height),
			dx:   (rand.Float64() - 0.5) * 4,
			dy:   (rand.Float64() - 0.5) * 2,
			size: 16 + float64(rand.Intn(8)),
			r:    rand.Float64() / 4.0,
			g:    rand.Float64() / 4.0,
			b:    rand.Float64() / 4.0,
		}
	}

	ticker := time.NewTicker(time.Second / 30)
	go func() {
		for range ticker.C {
			p.Tick()
			dirty.Dirty()
		}
	}()

	return p
}

func wrap(value, max float64) float64 {
	if max <= 0 {
		return 0 // Avoid division by zero, return 0 for non-positive max
	}

	// Wrap the value into the range [0, max)
	wrapped := value - max*float64(int(value/max))
	if wrapped < 0 {
		wrapped += max
	}
	return wrapped
}

func (p *Plasma) Tick() {
	p.rgb.Clear()

	for t := 0; t < len(p.blobs); t++ {
		blob := p.blobs[t]
		blob.x += blob.dx
		blob.y += blob.dy
		blob.x = wrap(blob.x, float64(p.rgb.width))
		blob.y = wrap(blob.y, float64(p.rgb.height))
		p.blobs[t] = blob

		// Draw the blob
		size2 := blob.size * blob.size
		for i := -blob.size; i < blob.size; i++ {
			for j := -blob.size; j < blob.size; j++ {
				d := (i * i) + ((j * j) * 2)
				if d < size2 {
					distance := float64(1-(math.Sqrt(float64(d)))/blob.size) * 255

					xx := blob.x + i
					yy := blob.y + j

					xx = wrap(xx, float64(p.rgb.width))
					yy = wrap(yy, float64(p.rgb.height))

					c := p.rgb.Get(int(xx), int(yy))

					c.R = uint8(math.Min(255, float64(c.R)+(distance*blob.r)))
					c.G = uint8(math.Min(255, float64(c.G)+(distance*blob.g)))
					c.B = uint8(math.Min(255, float64(c.B)+(distance*blob.b)))

					p.rgb.Set(int(xx), int(yy), c)
				}
			}
		}
	}
}
