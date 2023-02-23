package artist

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// Pattern is capable of drawing to a canvas within the bounds of a given
// clipping rectangle.
type Pattern interface {
	// Draw draws to destination, using the bounds of destination as a width
	// and height for things like gradients, bevels, etc. The pattern may
	// not draw outside the union of destination.Bounds() and clip. The
	// clipping rectangle effectively takes a subset of the pattern. To
	// change the bounds of the pattern itself, use canvas.Cut() on the
	// destination before passing it to Draw().
	Draw (destination canvas.Canvas, clip image.Rectangle)
}

// Draw lets you use several clipping rectangles to draw a pattern.
func Draw (
	destination canvas.Canvas,
	source      Pattern,
	clips       ...image.Rectangle,
) {
	for _, clip := range clips {
		source.Draw(destination, clip)
	}
}

// DrawBounds is like Draw, but lets you specify an overall bounding rectangle
// for the pattern. The destination is cut to this rectangle.
func DrawBounds (
	destination canvas.Canvas,
	bounds      image.Rectangle,
	source      Pattern,
	clips       ...image.Rectangle,
) {
	cut := canvas.Cut(destination, bounds)
	for _, clip := range clips {
		source.Draw(cut, clip)
	}
}
