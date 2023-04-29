package artist

import "image"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// Pattern is capable of drawing to a canvas within the bounds of a given
// clipping rectangle.
type Pattern interface {
	// Draw draws the pattern onto the destination canvas, using the
	// specified bounds. The given bounds can be smaller or larger than the
	// bounds of the destination canvas. The destination canvas can be cut
	// using canvas.Cut() to draw only a specific subset of a pattern.
	Draw (destination canvas.Canvas, bounds image.Rectangle)
}
