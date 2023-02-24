package artist

import "image/color"

// Orientation specifies an eight-way pattern orientation.
type Orientation int

const (
	OrientationVertical Orientation = iota
	OrientationDiagonalRight
	OrientationHorizontal
	OrientationDiagonalLeft
)

// Split is a pattern that is divided in half between two sub-patterns.
type Split struct {
	First  Pattern
	Second Pattern
	Orientation
}

// AtWhen satisfies the Pattern interface.
func (pattern Split) AtWhen (x, y, width, height int) (c color.RGBA) {
	var first bool
	switch pattern.Orientation {
	case OrientationVertical:
		first = x < width / 2
	case OrientationDiagonalRight:
		first = float64(x) / float64(width) +
			float64(y) / float64(height) < 1
	case OrientationHorizontal:
		first = y < height / 2
	case OrientationDiagonalLeft:
		first = float64(width - x) / float64(width) +
			float64(y) / float64(height) < 1
	}
	
	if first {
		return pattern.First.AtWhen(x, y, width, height)
	} else {
		return pattern.Second.AtWhen(x, y, width, height)
	}
}
