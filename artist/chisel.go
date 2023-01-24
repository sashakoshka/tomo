package artist

import "image/color"

// Beveled is a pattern that has a highlight section and a shadow section.
type Beveled struct {
	Highlight Pattern
	Shadow    Pattern
}

// AtWhen satisfies the Pattern interface.
func (pattern Beveled) AtWhen (x, y, width, height int) (c color.RGBA) {
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
		return pattern.Highlight.AtWhen(x, y, width, height)
	} else {
		return pattern.Shadow.AtWhen(x, y, width, height)
	}
}

// QuadBeveled is like Beveled, but with four sides. A pattern can be specified
// for each one.
type QuadBeveled [4]Pattern

// AtWhen satisfies the Pattern interface.
func (pattern QuadBeveled) AtWhen (x, y, width, height int) (c color.RGBA) {
	bottom := y > height / 2
	right  := x > width / 2
	top    := !bottom
	left   := !right
	side := 0
	
	switch {
	case top && left:
		if x > y { side = 0 } else { side = 3 }
		
	case top && right:
		if width - x < y { side = 1 } else { side = 0 }
		
	case bottom && left:
		if x > height - y { side = 2 } else { side = 3 }
		
	case bottom && right:
		if width - x < height - y { side = 1 } else { side = 2 }
		
	}

	return pattern[side].AtWhen(x, y, width, height)
}
