package artist

import "image/color"

// Striped is a pattern that produces stripes of two alternating colors.
type Striped struct {
	First  Stroke
	Second Stroke
	Orientation
}

// AtWhen satisfies the Pattern interface.
func (pattern Striped) AtWhen (x, y, width, height int) (c color.RGBA) {
	position := 0
	switch pattern.Orientation {
	case OrientationVertical:
		position = x
	case OrientationDiagonalRight:
		position = x + y
	case OrientationHorizontal:
		position = y
	case OrientationDiagonalLeft:
		position = x - y
	}

	phase := pattern.First.Weight + pattern.Second.Weight
	position %= phase
	if position < 0 {
		position += phase
	}
	
	if position < pattern.First.Weight {
		return pattern.First.AtWhen(x, y, width, height)
	} else {
		return pattern.Second.AtWhen(x, y, width, height)
	}
}
