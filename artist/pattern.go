package artist

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/shatter"

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
) (
	updatedRegion image.Rectangle,
) {
	for _, clip := range clips {
		source.Draw(destination, clip)
		updatedRegion = updatedRegion.Union(clip)
	}
	return
}

// DrawBounds lets you specify an overall bounding rectangle for drawing a
// pattern. The destination is cut to this rectangle.
func DrawBounds (
	destination canvas.Canvas,
	source      Pattern,
	bounds      image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	return Draw(canvas.Cut(destination, bounds), source, bounds)
}

// DrawShatter is like an inverse of Draw, drawing nothing in the areas
// specified in "rocks".
func DrawShatter (
	destination canvas.Canvas,
	source      Pattern,
	rocks       ...image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	tiles := shatter.Shatter(destination.Bounds(), rocks...)
	return Draw(destination, source, tiles...)
}

// AllocateSample returns a new canvas containing the result of a pattern. The
// resulting canvas can be sourced from shape drawing functions. I beg of you
// please do not call this every time you need to draw a shape with a pattern on
// it because that is horrible and cruel to the computer.
func AllocateSample (source Pattern, width, height int) (allocated canvas.Canvas) {
	allocated = canvas.NewBasicCanvas(width, height)
	source.Draw(allocated, allocated.Bounds())
	return
} 
