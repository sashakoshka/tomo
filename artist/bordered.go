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

// Stroke represents a stoke that has a weight and a pattern.
type Stroke struct {
	Weight int
	Pattern
}

type borderInternal struct {
	weight int
	stroke Pattern
	bounds image.Rectangle
	dx, dy int
}

// MultiBordered is a pattern that allows multiple borders of different lengths
// to be inset within one another. The final border is treated as a fill color,
// and its weight does not matter.
type MultiBordered struct {
	borders []borderInternal
	lastWidth, lastHeight int
	maxBorder int
}

// NewMultiBordered creates a new MultiBordered pattern from the given list of
// borders.
func NewMultiBordered (borders ...Stroke) (multi *MultiBordered) {
	internalBorders := make([]borderInternal, len(borders))
	for index, border := range borders {
		internalBorders[index].weight = border.Weight
		internalBorders[index].stroke = border.Pattern
	}
	return &MultiBordered { borders: internalBorders }
}

// AtWhen satisfies the Pattern interface.
func (multi *MultiBordered) AtWhen (x, y, width, height int) (c color.RGBA) {
	if multi.lastWidth != width || multi.lastHeight != height {
		multi.recalculate(width, height)
	}
	point := image.Point { x, y }
	for index := multi.maxBorder; index >= 0; index -- {
		border := multi.borders[index]
		if point.In(border.bounds) {
			return border.stroke.AtWhen (
				point.X - border.bounds.Min.X,
				point.Y - border.bounds.Min.Y,
				border.dx, border.dy)
		}
	}
	return
}

func (multi *MultiBordered) recalculate (width, height int) {
	bounds := image.Rect (0, 0, width, height)
	multi.maxBorder = 0
	for index, border := range multi.borders {
		multi.maxBorder = index
		multi.borders[index].bounds = bounds
		multi.borders[index].dx = bounds.Dx()
		multi.borders[index].dy = bounds.Dy()
		bounds = bounds.Inset(border.weight)
		if bounds.Empty() { break }
	}
}
