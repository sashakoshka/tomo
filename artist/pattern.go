package artist

import "image/color"

// Pattern is capable of generating a pattern pixel by pixel.
type Pattern interface {
	// AtWhen returns the color of the pixel located at (x, y) relative to
	// the origin point of the pattern (0, 0), when the pattern has the
	// specified width and height. Patterns may ignore the width and height
	// parameters, but it may be useful for some patterns such as gradients.
	AtWhen (x, y, width, height int) (color.RGBA)
}
