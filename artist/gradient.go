package artist

import "image/color"

// Gradient is a pattern that interpolates between two colors.
type Gradient struct {
	First  Pattern
	Second Pattern
	Orientation
}

// AtWhen satisfies the Pattern interface.
func (pattern Gradient) AtWhen (x, y, width, height int) (c color.RGBA) {
	var position float64
	switch pattern.Orientation {
	case OrientationVertical:
		position = float64(x) / float64(width)
	case OrientationDiagonalRight:
		position = (float64(x) / float64(width) +
			float64(y) / float64(height)) / 2
	case OrientationHorizontal:
		position = float64(y) / float64(height)
	case OrientationDiagonalLeft:
		position = (float64(width - x) / float64(width) +
			float64(y) / float64(height)) / 2
	}

	firstColor  := pattern.First.AtWhen(x, y, width, height)
	secondColor := pattern.Second.AtWhen(x, y, width, height)
	return LerpRGBA(firstColor, secondColor, position)
}

// Lerp linearally interpolates between two integer values.
func Lerp (first, second int, fac float64) (n int) {
	return int(float64(first) * (1 - fac) + float64(second) * fac)
}

// LerpRGBA linearally interpolates between two color.RGBA values.
func LerpRGBA (first, second color.RGBA, fac float64) (c color.RGBA) {
	return color.RGBA {
		R: uint8(Lerp(int(first.R), int(second.R), fac)),
		G: uint8(Lerp(int(first.G), int(second.G), fac)),
		B: uint8(Lerp(int(first.G), int(second.B), fac)),
	}
}
