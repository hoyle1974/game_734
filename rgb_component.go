package main

import (
	bitmaps "image/color"

	"github.com/fatih/color"
)

type RGBComponent struct {
	width  int
	height int
	data   [][]bitmaps.RGBA
}

func NewRGBComponent(width, height int) *RGBComponent {
	data := make([][]bitmaps.RGBA, height)
	for y := 0; y < height; y++ {
		data[y] = make([]bitmaps.RGBA, width)
	}
	return &RGBComponent{width: width, height: height, data: data}
}

func mapGrayscaleToASCII(value uint8) string {
	ascii := "█"
	index := int(value) * (len(ascii) - 1) / 255
	return string(ascii[index])
}

func (p *RGBComponent) Get(x int, y int) bitmaps.RGBA {
	if x >= 0 && y >= 0 && x < p.width && y < p.height {
		return p.data[y][x]
	}
	return bitmaps.RGBA{}
}
func (p *RGBComponent) Set(x int, y int, c bitmaps.RGBA) {
	if x >= 0 && y >= 0 && x < p.width && y < p.height {
		p.data[y][x] = c
	}
}

func (p *RGBComponent) Clear() {
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			p.data[y][x] = bitmaps.RGBA{0, 0, 0, 0}
		}
	}
}

func (p *RGBComponent) Render() *Buffer {
	buffer := NewBuffer(p.width, p.height)

	for y := 0; y < buffer.height; y++ {
		for x := 0; x < buffer.width; x++ {
			// c := mapGrayscaleToASCII(((p.data[y][x].R + p.data[y][x].G + p.data[y][x].B) / 3))
			buffer.set(x, y, color.RGB(int(p.data[y][x].R), int(p.data[y][x].B), int(p.data[y][x].G)).Sprintf("\u2588"))
		}
	}

	return buffer
}
