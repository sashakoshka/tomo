package testing

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Artist is an element that displays shapes and patterns drawn by the artist
// package in order to test it.
type Artist struct {
	*core.Core
	core core.CoreControl
	cellBounds image.Rectangle
}

// NewArtist creates a new artist test element.
func NewArtist () (element *Artist) {
	element = &Artist { }
	element.Core, element.core = core.NewCore(element)
	element.core.SetMinimumSize(400, 300)
	return
}

func (element *Artist) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	bounds := element.Bounds()
	element.cellBounds.Max.X = bounds.Dx() / 4
	element.cellBounds.Max.Y = bounds.Dy() / 4

	// 0, 0
	artist.FillRectangle (
		element,
		artist.Chiseled {
			Highlight: artist.NewUniform(hex(0xFF0000FF)),
			Shadow:    artist.NewUniform(hex(0x0000FFFF)),
		},
		element.cellAt(0, 0))

	// 1, 0
	artist.StrokeRectangle (
		element,
		artist.NewUniform(hex(0x00FF00FF)), 3,
		element.cellAt(1, 0))

	// 2, 0
	artist.FillRectangle (
		element,
		artist.NewMultiBordered (
			artist.Stroke { Pattern: uhex(0xFF0000FF), Weight: 1 },
			artist.Stroke { Pattern: uhex(0x888800FF), Weight: 2 },
			artist.Stroke { Pattern: uhex(0x00FF00FF), Weight: 3 },
			artist.Stroke { Pattern: uhex(0x008888FF), Weight: 4 },
			artist.Stroke { Pattern: uhex(0x0000FFFF), Weight: 5 },
			),
		element.cellAt(2, 0))

	// 3, 0
	artist.FillRectangle (
		element,
		artist.Bordered {
			Stroke: artist.Stroke { Pattern: uhex(0x0000FFFF), Weight: 5 },
			Fill: uhex(0xFF0000FF),
		},
		element.cellAt(3, 0))

	// 0, 1 - 0, 3
	for x := 0; x < 4; x ++ {
		artist.FillRectangle (
			element,
			artist.Striped {
				First:  artist.Stroke { Pattern: uhex(0xFF8800FF), Weight: 7 },
				Second: artist.Stroke { Pattern: uhex(0x0088FFFF), Weight: 2 },
				Direction: artist.StripeDirection(x),
				
			},
			element.cellAt(x, 1))
	}

	for x := 0; x < 4; x ++ {
		element.lines(x, element.cellAt(x, 2))
	}
}

func (element *Artist) lines (weight int, bounds image.Rectangle) {
	bounds = bounds.Inset(8)
	c := uhex(0xFFFFFFFF)
	artist.Line(element, c, weight, bounds.Min, bounds.Max)
	artist.Line (
		element, c, weight,
		image.Pt(bounds.Max.X, bounds.Min.Y),
		image.Pt(bounds.Min.X, bounds.Max.Y))
	artist.Line (
		element, c, weight,
		image.Pt(bounds.Max.X, bounds.Min.Y + 16),
		image.Pt(bounds.Min.X, bounds.Max.Y - 16))
	artist.Line (
		element, c, weight,
		image.Pt(bounds.Min.X, bounds.Min.Y + 16),
		image.Pt(bounds.Max.X, bounds.Max.Y - 16))
	artist.Line (
		element, c, weight,
		image.Pt(bounds.Min.X + 20, bounds.Min.Y),
		image.Pt(bounds.Max.X - 20, bounds.Max.Y))
	artist.Line (
		element, c, weight,
		image.Pt(bounds.Max.X - 20, bounds.Min.Y),
		image.Pt(bounds.Min.X + 20, bounds.Max.Y))
	artist.Line (
		element, c, weight,
		image.Pt(bounds.Min.X, bounds.Min.Y + bounds.Dy() / 2),
		image.Pt(bounds.Max.X, bounds.Min.Y + bounds.Dy() / 2))
	artist.Line (
		element, c, weight,
		image.Pt(bounds.Min.X + bounds.Dx() / 2, bounds.Min.Y),
		image.Pt(bounds.Min.X + bounds.Dx() / 2, bounds.Max.Y))
}

func (element *Artist) cellAt (x, y int) (image.Rectangle) {
	return element.cellBounds.Add (image.Pt (
		x * element.cellBounds.Dx(),
		y * element.cellBounds.Dy()))
}

func hex (n uint32) (c color.RGBA) {
	c.A = uint8(n)
	c.B = uint8(n >>  8)
	c.G = uint8(n >> 16)
	c.R = uint8(n >> 24)
	return
}

func uhex (n uint32) (artist.Pattern) {
	return artist.NewUniform (color.RGBA {
		A: uint8(n),
		B: uint8(n >>  8),
		G: uint8(n >> 16),
		R: uint8(n >> 24),
	})
}