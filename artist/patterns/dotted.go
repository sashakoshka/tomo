package artist

import "math"
import "image/color"

// Dotted is a pattern that produces a grid of circles.
type Dotted struct {
	Background Pattern
	Foreground Pattern
	Size int
	Spacing int
}

// AtWhen satisfies the Pattern interface.
func (pattern Dotted) AtWhen (x, y, width, height int) (c color.RGBA) {
	xm := x % pattern.Spacing
	ym := y % pattern.Spacing
	if xm < 0 { xm += pattern.Spacing }
	if ym < 0 { xm += pattern.Spacing }
	radius  := float64(pattern.Size) / 2
	spacing := float64(pattern.Spacing) / 2 - 0.5
	xf := float64(xm) - spacing
	yf := float64(ym) - spacing

	if math.Sqrt(xf * xf + yf * yf) > radius {
		return pattern.Background.AtWhen(x, y, width, height)
	} else {
		return pattern.Foreground.AtWhen(x, y, width, height)
	}
}
