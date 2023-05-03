// Package artutil provides utility functions for working with graphical types
// defined in artist, canvas, and image.
package artutil

import "image"
import "image/color"
import "tomo/artist"
import "tomo/shatter"

// Fill fills the destination canvas with the given pattern.
func Fill (destination artist.Canvas, source artist.Pattern) (updated image.Rectangle) {
	source.Draw(destination, destination.Bounds())
	return destination.Bounds()
}

// DrawClip lets you draw several subsets of a pattern at once.
func DrawClip (
	destination artist.Canvas,
	source      artist.Pattern,
	bounds      image.Rectangle,
	subsets     ...image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	for _, subset := range subsets {
		source.Draw(artist.Cut(destination, subset), bounds)
		updatedRegion = updatedRegion.Union(subset)
	}
	return
}

// DrawShatter is like an inverse of DrawClip, drawing nothing in the areas
// specified by "rocks".
func DrawShatter (
	destination artist.Canvas,
	source      artist.Pattern,
	bounds      image.Rectangle,
	rocks       ...image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	tiles := shatter.Shatter(bounds, rocks...)
	return DrawClip(destination, source, bounds, tiles...)
}

// AllocateSample returns a new canvas containing the result of a pattern. The
// resulting canvas can be sourced from shape drawing functions. I beg of you
// please do not call this every time you need to draw a shape with a pattern on
// it because that is horrible and cruel to the computer.
func AllocateSample (source artist.Pattern, width, height int) artist.Canvas {
	allocated := artist.NewBasicCanvas(width, height)
	Fill(allocated, source)
	return allocated
}

// Hex creates a color.RGBA value from an RGBA integer value.
func Hex (color uint32) (c color.RGBA) {
	c.A = uint8(color)
	c.B = uint8(color >>  8)
	c.G = uint8(color >> 16)
	c.R = uint8(color >> 24)
	return
}
