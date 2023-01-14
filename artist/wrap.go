package artist

import "image"
import "image/color"

// WrappedPattern is a pattern that is able to behave like an image.Image.
type WrappedPattern struct {
	Pattern
	Width, Height int
}

// At satisfies the image.Image interface.
func (pattern WrappedPattern) At (x, y int) (c color.Color) {
	return pattern.Pattern.AtWhen(x, y, pattern.Width, pattern.Height)
}

// Bounds satisfies the image.Image interface.
func (pattern WrappedPattern) Bounds () (rectangle image.Rectangle) {
	rectangle.Min = image.Point { -1e9, -1e9 }
	rectangle.Max = image.Point {  1e9,  1e9 }
	return
}

// ColorModel satisfies the image.Image interface.
func (pattern WrappedPattern) ColorModel () (model color.Model) {
	return color.RGBAModel
}
