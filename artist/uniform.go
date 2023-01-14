package artist

import "image"
import "image/color"

// Uniform is an infinite-sized pattern of uniform color. It implements the
// color.Color, color.Model, and image.Image interfaces.
type Uniform struct {
	C color.RGBA
}

// NewUniform returns a new Uniform image of the given color.
func NewUniform (c color.Color) (uniform *Uniform) {
	uniform = &Uniform { }
	r, g, b, a := c.RGBA()
	uniform.C.R = uint8(r >> 8)
	uniform.C.G = uint8(g >> 8)
	uniform.C.B = uint8(b >> 8)
	uniform.C.A = uint8(a >> 8)
	return
}

func (uniform *Uniform) RGBA () (r, g, b, a uint32) {
	r = uint32(uniform.C.R) << 8 | uint32(uniform.C.R)
	g = uint32(uniform.C.G) << 8 | uint32(uniform.C.G)
	b = uint32(uniform.C.B) << 8 | uint32(uniform.C.B)
	a = uint32(uniform.C.A) << 8 | uint32(uniform.C.A)
	return
}

func (uniform *Uniform) ColorModel () (model color.Model) {
	return uniform
}

func (uniform *Uniform) Convert (in color.Color) (c color.Color) {
	return uniform.C
}

func (uniform *Uniform) Bounds () (rectangle image.Rectangle) {
	rectangle.Min = image.Point { -1e9, -1e9 }
	rectangle.Max = image.Point {  1e9,  1e9 }
	return
}

func (uniform *Uniform) At (x, y int) (c color.Color) {
	return uniform.C
}

func (uniform *Uniform) AtWhen (x, y, width, height int) (c color.RGBA) {
	return uniform.C
}

func (uniform *Uniform) RGBA64At (x, y int) (c color.RGBA64) {
	r := uint16(uniform.C.R) << 8 | uint16(uniform.C.R)
	g := uint16(uniform.C.G) << 8 | uint16(uniform.C.G)
	b := uint16(uniform.C.B) << 8 | uint16(uniform.C.B)
	a := uint16(uniform.C.A) << 8 | uint16(uniform.C.A)
	return color.RGBA64 { R: r, G: g, B: b, A: a }
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (uniform *Uniform) Opaque () (opaque bool) {
	return uniform.C.A == 0xFF
}
