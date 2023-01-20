package artist

import "image/color"

// Chiseled is a pattern that has a highlight section and a shadow section.
type Chiseled struct {
	Highlight Pattern
	Shadow    Pattern
}

// AtWhen satisfies the Pattern interface.
func (chiseled Chiseled) AtWhen (x, y, width, height int) (c color.RGBA) {
	var highlighted  bool
	var bottomCorner bool
	
	if width > height {
		bottomCorner = y > height / 2
	} else {
		bottomCorner = x < width / 2
	}
	
	if bottomCorner {
		highlighted = float64(x) < float64(height) - float64(y)
	} else {
		highlighted = float64(width) - float64(x) > float64(y)
	}

	if highlighted {
		return chiseled.Highlight.AtWhen(x, y, width, height)
	} else {
		return chiseled.Shadow.AtWhen(x, y, width, height)
	}
}
