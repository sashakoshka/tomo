package testing

import "fmt"
import "time"
import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/defaultfont"
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
	element.core.SetMinimumSize(480, 600)
	return
}

func (element *Artist) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	bounds := element.Bounds()
	element.cellBounds.Max.X = bounds.Dx() / 5
	element.cellBounds.Max.Y = (bounds.Dy() - 48) / 8

	drawStart := time.Now()

	// 0, 0
	artist.FillRectangle (
		element,
		artist.Beveled {
			artist.NewUniform(hex(0xFF0000FF)),
			artist.NewUniform(hex(0x0000FFFF)),
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

	// 4, 0
	artist.FillRectangle (
		element,
		artist.Padded {
			Stroke: uhex(0xFFFFFFFF),
			Fill: uhex(0x666666FF),
			Sides: []int { 4, 13, 2, 0 },
		},
		element.cellAt(4, 0))

	// 0, 1 - 3, 1
	for x := 0; x < 4; x ++ {
		artist.FillRectangle (
			element,
			artist.Striped {
				First:  artist.Stroke { Pattern: uhex(0xFF8800FF), Weight: 7 },
				Second: artist.Stroke { Pattern: uhex(0x0088FFFF), Weight: 2 },
				Orientation: artist.Orientation(x),
				
			},
			element.cellAt(x, 1))
	}

	// 0, 2 - 3, 2
	for x := 0; x < 4; x ++ {
		element.lines(x + 1, element.cellAt(x, 2))
	}

	// 0, 3
	artist.StrokeRectangle (
		element,uhex(0x888888FF), 1,
		element.cellAt(0, 3))
	artist.FillEllipse(element, uhex(0x00FF00FF), element.cellAt(0, 3))

	// 1, 3 - 3, 3
	for x := 1; x < 4; x ++ {
		artist.StrokeRectangle (
			element,uhex(0x888888FF), 1,
			element.cellAt(x, 3))
		artist.StrokeEllipse (
			element,
			[]artist.Pattern {
				uhex(0xFF0000FF),
				uhex(0x00FF00FF),
				uhex(0xFF00FFFF),
			} [x - 1],
			x, element.cellAt(x, 3))
	}

	// 0, 4 - 3, 4
	for x := 0; x < 4; x ++ {
		artist.FillEllipse (
			element, 
			artist.Split {
				First:  uhex(0xFF0000FF),
				Second: uhex(0x0000FFFF),
				Orientation: artist.Orientation(x),
			},
			element.cellAt(x, 4))
	}
	
	// how long did that take to render?
	drawTime := time.Since(drawStart)
	textDrawer := artist.TextDrawer { }
	textDrawer.SetFace(defaultfont.FaceRegular)
	textDrawer.SetText ([]rune (fmt.Sprintf (
		"%dms\n%dus",
		drawTime.Milliseconds(),
		drawTime.Microseconds())))
	textDrawer.Draw(element, uhex(0xFFFFFFFF), image.Pt(8, bounds.Max.Y - 24))
	
	// 0, 5
	artist.FillRectangle (
		element,
		artist.QuadBeveled {
			uhex(0x880000FF),
			uhex(0x00FF00FF),
			uhex(0x0000FFFF),
			uhex(0xFF00FFFF),
		},
		element.cellAt(0, 5))
		
	// 1, 5
	artist.FillRectangle (
		element,
		artist.Checkered {
			First: artist.QuadBeveled {
				uhex(0x880000FF),
				uhex(0x00FF00FF),
				uhex(0x0000FFFF),
				uhex(0xFF00FFFF),
			},
			Second: artist.Striped {
				First:  artist.Stroke { Pattern: uhex(0xFF8800FF), Weight: 1 },
				Second: artist.Stroke { Pattern: uhex(0x0088FFFF), Weight: 1 },
				Orientation: artist.OrientationVertical,
			},
			CellWidth: 32,
			CellHeight: 16,
		},
		element.cellAt(1, 5))
	
	// 2, 5
	artist.FillRectangle (
		element,
		artist.Dotted {
			Foreground: uhex(0x00FF00FF),
			Background: artist.Checkered {
				First:  uhex(0x444444FF),
				Second: uhex(0x888888FF),
				CellWidth: 16,
				CellHeight: 16,
			},
			Size: 8,
			Spacing: 16,
		},
		element.cellAt(2, 5))
		
	// 3, 5
	artist.FillRectangle (
		element,
		artist.Tiled {
			Pattern: artist.QuadBeveled {
				uhex(0x880000FF),
				uhex(0x00FF00FF),
				uhex(0x0000FFFF),
				uhex(0xFF00FFFF),
			},
			CellWidth: 17,
			CellHeight: 23,
		},
		element.cellAt(3, 5))
		
	// 0, 6 - 3, 6
	for x := 0; x < 4; x ++ {
		artist.FillRectangle (
			element, 
			artist.Gradient {
				First:  uhex(0xFF0000FF),
				Second: uhex(0x0000FFFF),
				Orientation: artist.Orientation(x),
			},
			element.cellAt(x, 6))
	}

	// 0, 7
	artist.FillEllipse (
		element,
		artist.EllipticallyBordered {
			Fill: artist.Gradient {
				First:  uhex(0x00FF00FF),
				Second: uhex(0x0000FFFF),
				Orientation: artist.OrientationVertical,
			},
			Stroke: artist.Stroke { Pattern: uhex(0x00FF00), Weight: 5 },
		},
		element.cellAt(0, 7))

	// 1, 7
	artist.FillRectangle (
		element,
		artist.Noisy {
			Low:  uhex(0x000000FF),
			High: uhex(0xFFFFFFFF),
			Seed: 0,
		},
		element.cellAt(1, 7),
	)

	// 2, 7
	artist.FillRectangle (
		element,
		artist.Noisy {
			Low:  uhex(0x000000FF),
			High: artist.Gradient {
				First:  uhex(0x000000FF),
				Second: uhex(0xFFFFFFFF),
				Orientation: artist.OrientationVertical,
			},
			Seed: 0,
		},
		element.cellAt(2, 7),
	)

	// 3, 7
	artist.FillRectangle (
		element,
		artist.Noisy {
			Low:  uhex(0x000000FF),
			High: uhex(0xFFFFFFFF),
			Seed: 0,
			Harsh: true,
		},
		element.cellAt(3, 7),
	)
}

func (element *Artist) lines (weight int, bounds image.Rectangle) {
	bounds = bounds.Inset(4)
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
