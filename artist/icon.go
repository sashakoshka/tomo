package artist

import "image"
import "image/color"

type Icon interface {
	// Draw draws the icon to the destination canvas at the specified point,
	// using the specified color (if the icon is monochrome).
	Draw (destination Canvas, color color.RGBA, at image.Point)

	// Bounds returns the bounds of the icon.
	Bounds () image.Rectangle
}
