package artist

import "math"
import "image/color"

// EllipticallyBordered is a pattern with a border and a fill that is elliptical
// in shape.
type EllipticallyBordered struct {
	Fill Pattern
	Stroke
}

// AtWhen satisfies the Pattern interface.
func (pattern EllipticallyBordered) AtWhen (x, y, width, height int) (c color.RGBA) {
	xf := (float64(x) + 0.5) / float64(width ) * 2 - 1
	yf := (float64(y) + 0.5) / float64(height) * 2 - 1
	distance := math.Sqrt(xf * xf + yf * yf)

	var radius float64
	if width < height {
		// vertical
		radius = 1 - float64(pattern.Weight * 2) / float64(width)
	} else {
		// horizontal
		radius = 1 - float64(pattern.Weight * 2) / float64(height)
	}

	if distance < radius {
		return pattern.Fill.AtWhen(x, y, width, height)
	} else {
		return pattern.Stroke.AtWhen(x, y, width, height)
	}
}
