package artist

import "image"
import "image/color"

// Uniform is an infinite-sized pattern of uniform color. It implements the
// Pattern, color.Color, color.Model, and image.Image interfaces.
type Uniform color.RGBA

// NewUniform returns a new Uniform image of the given color.
func NewUniform (c color.Color) (uniform Uniform) {
	r, g, b, a := c.RGBA()
	uniform.R = uint8(r >> 8)
	uniform.G = uint8(g >> 8)
	uniform.B = uint8(b >> 8)
	uniform.A = uint8(a >> 8)
	return
}

// ColorModel satisfies the image.Image interface.
func (uniform Uniform) ColorModel () (model color.Model) {
	return uniform
}

// Convert satisfies the color.Model interface.
func (uniform Uniform) Convert (in color.Color) (c color.Color) {
	return color.RGBA(uniform)
}

// Bounds satisfies the image.Image interface.
func (uniform Uniform) Bounds () (rectangle image.Rectangle) {
	rectangle.Min = image.Point { -1e9, -1e9 }
	rectangle.Max = image.Point {  1e9,  1e9 }
	return
}

// At satisfies the image.Image interface.
func (uniform Uniform) At (x, y int) (c color.Color) {
	return color.RGBA(uniform)
}

// AtWhen satisfies the Pattern interface.
func (uniform Uniform) AtWhen (x, y, width, height int) (c color.RGBA) {
	return color.RGBA(uniform)
}

// RGBA satisfies the color.Color interface.
func (uniform Uniform) RGBA () (r, g, b, a uint32) {
	return color.RGBA(uniform).RGBA()
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (uniform Uniform) Opaque () (opaque bool) {
	return uniform.A == 0xFF
}
