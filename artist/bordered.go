package artist

import "image"
import "image/color"

// Bordered is a pattern with a border and a fill.
type Bordered struct {
	Fill Pattern
	Stroke
}

// AtWhen satisfies the Pattern interface.
func (pattern Bordered) AtWhen (x, y, width, height int) (c color.RGBA) {
	outerBounds := image.Rectangle { Max: image.Point { width, height }}
	innerBounds := outerBounds.Inset(pattern.Weight)
	if (image.Point { x, y }).In (innerBounds) {
		return pattern.Fill.AtWhen (
			x - pattern.Weight,
			y - pattern.Weight,
			innerBounds.Dx(), innerBounds.Dy())
	} else {
		return pattern.Stroke.AtWhen(x, y, width, height)
	}
}
