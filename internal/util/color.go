package util

import (
	"image/color"
	"unsafe"
)

type Color struct {
	R uint8
	G uint8
	B uint8
	W uint8
}

func (c *Color) DimColor(b float64) Color {
	return Color{
		R: uint8(float64(c.R) * b),
		G: uint8(float64(c.G) * b),
		B: uint8(float64(c.B) * b),
	}
}

func (c *Color) ToBits() uint32 {
	return uint32(c.W)<<24 | uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}

func (c *Color) ToSysColor() color.RGBA {
	return color.RGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: 255,
	}
}

type PixelBuffer struct {
	RawWidth  int
	RawHeight int
	Width     int
	Height    int
	minX      int
	minY      int
	stride    int
	buffer    []Color
}

func NewPixelBuffer(width, height, minX, minY, stride int) *PixelBuffer {
	return &PixelBuffer{
		Width:  width,
		Height: height,
		minX:   minX,
		minY:   minY,
		stride: stride,
		buffer: make([]Color, width*height*stride),
	}
}

func (pb *PixelBuffer) GetUnsafePointer() unsafe.Pointer {
	return unsafe.Pointer(&pb.buffer[0])
}

func (pb *PixelBuffer) GetPixel(p *Point) Color {
	mappedX := (p.X - pb.minX) * pb.stride
	mappedY := (p.Y - pb.minY) * pb.stride
	return pb.buffer[mappedX+mappedY*pb.Width]
}

func (pb *PixelBuffer) GetPixelPointer(p *Point) *Color {
	mappedX := (p.X - pb.minX) * pb.stride
	mappedY := (p.Y - pb.minY) * pb.stride
	return &pb.buffer[mappedX+mappedY*pb.Width]
}

func (pb *PixelBuffer) BlackOut() {
	for i := range pb.buffer {
		pb.buffer[i] = Color{}
	}
}
