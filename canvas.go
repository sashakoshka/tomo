package tomo

import "image"
import "image/draw"
import "image/color"

// Canvas is like draw.Image but is also able to return a raw pixel buffer for
// more efficient drawing. This interface can be easily satisfied using a
// BasicCanvas struct.
type Canvas interface {
	draw.Image
	Buffer () (data []color.RGBA, stride int)
}

// BasicCanvas is a general purpose implementation of tomo.Canvas.
type BasicCanvas struct {
	pix    []color.RGBA
	stride int
	rect   image.Rectangle
}

// NewBasicCanvas creates a new basic canvas with the specified width and
// height, allocating a buffer for it.
func NewBasicCanvas (width, height int) (canvas BasicCanvas) {
	canvas.pix    = make([]color.RGBA, height * width)
	canvas.stride = width
	canvas.rect = image.Rect(0, 0, width, height)
	return
}

// you know what it do
func (canvas BasicCanvas) Bounds () (bounds image.Rectangle) {
	return canvas.rect
}

// you know what it do
func (canvas BasicCanvas) At (x, y int) (color.Color) {
	if !image.Pt(x, y).In(canvas.rect) { return nil }
	return canvas.pix[x + y * canvas.stride]
}

// you know what it do
func (canvas BasicCanvas) ColorModel () (model color.Model) {
	return color.RGBAModel
}

// you know what it do
func (canvas BasicCanvas) Set (x, y int, c color.Color) {
	if !image.Pt(x, y).In(canvas.rect) { return }
	r, g, b, a := c.RGBA()
	canvas.pix[x + y * canvas.stride] = color.RGBA {
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}

// you know what it do
func (canvas BasicCanvas) Buffer () (data []color.RGBA, stride int) {
	return canvas.pix, canvas.stride
}

// Cut returns a sub-canvas of a given canvas.
func Cut (canvas Canvas, bounds image.Rectangle) (reduced BasicCanvas) {
	// println(canvas.Bounds().String(), bounds.String())
	bounds = bounds.Intersect(canvas.Bounds())
	if bounds.Empty() { return }
	reduced.rect = bounds
	reduced.pix, reduced.stride = canvas.Buffer()
	return
}
