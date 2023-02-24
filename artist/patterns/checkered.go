package artist

import "image/color"

// Checkered is a pattern that produces a grid of two alternating colors.
type Checkered struct {
	First  Pattern
	Second Pattern
	CellWidth, CellHeight int
}

// AtWhen satisfies the Pattern interface.
func (pattern Checkered) AtWhen (x, y, width, height int) (c color.RGBA) {
	twidth  := pattern.CellWidth  * 2
	theight := pattern.CellHeight * 2
	x %= twidth
	y %= theight
	if x < 0 { x += twidth  }
	if y < 0 { x += theight }

	n := 0
	if x >= pattern.CellWidth  { n ++ }
	if y >= pattern.CellHeight { n ++ }
	
	x %= pattern.CellWidth
	y %= pattern.CellHeight

	if n % 2 == 0 {
		return pattern.First.AtWhen (
			x, y, pattern.CellWidth, pattern.CellHeight)
	} else {
		return pattern.Second.AtWhen (
			x, y, pattern.CellWidth, pattern.CellHeight)
	}
}

// Tiled is a pattern that tiles another pattern accross a grid.
type Tiled struct {
	Pattern
	CellWidth, CellHeight int
}

// AtWhen satisfies the Pattern interface.
func (pattern Tiled) AtWhen (x, y, width, height int) (c color.RGBA) {
	x %= pattern.CellWidth
	y %= pattern.CellHeight
	if x < 0 { x += pattern.CellWidth  }
	if y < 0 { y += pattern.CellHeight }
	return pattern.Pattern.AtWhen (
		x, y, pattern.CellWidth, pattern.CellHeight)
}
