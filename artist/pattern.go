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
