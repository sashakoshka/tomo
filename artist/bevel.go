package artist

import "image/color"

// Beveled is a pattern that has a highlight section and a shadow section.
type Beveled [2]Pattern

// AtWhen satisfies the Pattern interface.
func (pattern Beveled) AtWhen (x, y, width, height int) (c color.RGBA) {
	return QuadBeveled {
		pattern[0],
		pattern[1],
		pattern[1],
		pattern[0],
	}.AtWhen(x, y, width, height)
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
		if x < y { side = 3 } else { side = 0 }
		
	case top && right:
		if width - x > y { side = 0 } else { side = 1 }
		
	case bottom && left:
		if x < height - y { side = 3 } else { side = 2 }
		
	case bottom && right:
		if width - x > height - y { side = 2 } else { side = 1 }
		
	}

	return pattern[side].AtWhen(x, y, width, height)
}
