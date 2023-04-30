package tomo

import "image"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Entity is a handle given to elements by the backend. Extended entity
// interfaces are defined in the ability module.
type Entity interface {
	// Invalidate marks the element's current visual as invalid. At the end
	// of every event, the backend will ask all invalid entities to redraw
	// themselves.
	Invalidate ()

	// Bounds returns the bounds of the element to be used for drawing and
	// layout.
	Bounds () image.Rectangle

	// Window returns the window that the element is in.
	Window () Window

	// SetMinimumSize reports to the system what the element's minimum size
	// can be. The minimum size of child elements should be taken into
	// account when calculating this.
	SetMinimumSize (width, height int)

	// DrawBackground asks the parent element to draw its background pattern
	// to a canvas. This should be used for transparent elements like text
	// labels. If there is no parent element (that is, the element is
	// directly inside of the window), the backend will draw a default
	// background pattern.
	DrawBackground (artist.Canvas)
}
