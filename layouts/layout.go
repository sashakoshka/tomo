package layouts

import "image"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements"

// LayoutEntry associates an element with layout and positioning information so
// it can be arranged by a Layout.
type LayoutEntry struct {
	elements.Element
	Bounds image.Rectangle
	Expand bool
}

// Layout is capable of arranging elements within a container. It is also able
// to determine the minimum amount of room it needs to do so.
type Layout interface {
	// Arrange takes in a slice of entries and a bounding width and height,
	// and changes the position of the entiries in the slice so that they
	// are properly laid out. The given width and height should not be less
	// than what is returned by MinimumSize.
	Arrange (
		entries []LayoutEntry,
		margin  image.Point,
		padding artist.Inset,
		bounds image.Rectangle,
	)

	// MinimumSize returns the minimum width and height that the layout
	// needs to properly arrange the given slice of layout entries.
	MinimumSize (
		entries []LayoutEntry,
		margin  image.Point,
		padding artist.Inset,
	) (
		width, height int,
	)
}
