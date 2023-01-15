package artist

import "image"
import "image/color"

// Border represents a border that can be fed to MultiBorder.
type Border struct {
	Weight int
	Stroke Pattern
	bounds image.Rectangle
	dx, dy int
}

// MultiBorder is a pattern that allows multiple borders of different lengths to
// be inset within one another. The final border is treated as a fill color, and
// its weight does not matter.
type MultiBorder struct {
	borders []Border
	lastWidth, lastHeight int
	maxBorder int
}

// NewMultiBorder creates a new MultiBorder pattern from the given list of
// borders.
func NewMultiBorder (borders ...Border) (multi *MultiBorder) {
	return &MultiBorder { borders: borders }
}

// AtWhen satisfies the Pattern interface.
func (multi *MultiBorder) AtWhen (x, y, width, height int) (c color.RGBA) {
	if multi.lastWidth != width || multi.lastHeight != height {
		multi.recalculate(width, height)
	}
	point := image.Point { x, y }
	for index := multi.maxBorder; index >= 0; index -- {
		border := multi.borders[index]
		if point.In(border.bounds) {
			return border.Stroke.AtWhen (
				point.X - border.bounds.Min.X,
				point.Y - border.bounds.Min.Y,
				border.dx, border.dy)
		}
	}
	return
}

func (multi *MultiBorder) recalculate (width, height int) {
	bounds := image.Rect (0, 0, width, height)
	multi.maxBorder = 0
	for index, border := range multi.borders {
		multi.maxBorder = index
		multi.borders[index].bounds = bounds
		multi.borders[index].dx = bounds.Dx()
		multi.borders[index].dy = bounds.Dy()
		bounds = bounds.Inset(border.Weight)
		if bounds.Empty() { break }
	}
}
