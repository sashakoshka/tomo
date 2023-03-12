package artist

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/shatter"

// Pattern is capable of drawing to a canvas within the bounds of a given
// clipping rectangle.
type Pattern interface {
	// Draw draws the pattern onto the destination canvas, using the
	// specified bounds. The given bounds can be smaller or larger than the
	// bounds of the destination canvas. The destination canvas can be cut
	// using canvas.Cut() to draw only a specific subset of a pattern.
	Draw (destination canvas.Canvas, bounds image.Rectangle)
}

// Fill fills the destination canvas with the given pattern.
func Fill (destination canvas.Canvas, source Pattern) (updated image.Rectangle) {
	source.Draw(destination, destination.Bounds())
	return destination.Bounds()
}

// DrawClip lets you draw several subsets of a pattern at once.
func DrawClip (
	destination canvas.Canvas,
	source      Pattern,
	bounds      image.Rectangle,
	subsets     ...image.Rectangle,
) (
	updatedRegion image.Rectangle,
) {
	for _, subset := range subsets {
		source.Draw(canvas.Cut(destination, subset), bounds)
		updatedRegion = updatedRegion.Union(subset)
	}
	return
}

// DrawShatter is like an inverse of DrawClip, drawing nothing in the areas
// specified by "rocks".
func DrawShatter (
	destination canvas.Canvas,
	source      Pattern,
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
func AllocateSample (source Pattern, width, height int) canvas.Canvas {
	allocated := canvas.NewBasicCanvas(width, height)
	Fill(allocated, source)
	return allocated
} 
