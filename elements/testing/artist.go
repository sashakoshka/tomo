package testing

import "fmt"
import "time"
import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/shatter"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/defaultfont"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/artist/shapes"
import "git.tebibyte.media/sashakoshka/tomo/artist/patterns"

// Artist is an element that displays shapes and patterns drawn by the artist
// package in order to test it.
type Artist struct {
	*core.Core
	core core.CoreControl
}

// NewArtist creates a new artist test element.
func NewArtist () (element *Artist) {
	element = &Artist { }
	element.Core, element.core = core.NewCore(element.draw)
	element.core.SetMinimumSize(240, 240)
	return
}

func (element *Artist) draw () {
	bounds := element.Bounds()
	patterns.Uhex(0x000000FF).Draw(element.core, bounds)

	drawStart := time.Now()

	// 0, 0 - 3, 0
	for x := 0; x < 4; x ++ {
		element.colorLines(x + 1, element.cellAt(x, 0).Bounds())
	}

	// 4, 0
	c40 := element.cellAt(4, 0)
	shapes.StrokeColorRectangle(c40, artist.Hex(0x888888FF), c40.Bounds(), 1)
	shapes.ColorLine (
		c40, artist.Hex(0xFF0000FF), 1,
		c40.Bounds().Min, c40.Bounds().Max)

	// 0, 1
	c01 := element.cellAt(0, 1)
	shapes.StrokeColorRectangle(c01, artist.Hex(0x888888FF), c01.Bounds(), 1)
	shapes.FillColorEllipse(element.core, artist.Hex(0x00FF00FF), c01.Bounds())

	// 1, 1 - 3, 1
	for x := 1; x < 4; x ++ {
		c := element.cellAt(x, 1)
		shapes.StrokeColorRectangle (
			element.core, artist.Hex(0x888888FF),
			c.Bounds(), 1)
		shapes.StrokeColorEllipse (
			element.core,
			[]color.RGBA {
				artist.Hex(0xFF0000FF),
				artist.Hex(0x00FF00FF),
				artist.Hex(0xFF00FFFF),
			} [x - 1],
			c.Bounds(), x)
	}

	// 4, 1
	c41 := element.cellAt(4, 1)
	shatterPos := c41.Bounds().Min
	rocks := []image.Rectangle {
		image.Rect(3, 12, 13, 23).Add(shatterPos),
		// image.Rect(30, 10, 40, 23).Add(shatterPos),
		image.Rect(55, 40, 70, 49).Add(shatterPos),
		image.Rect(30, -10, 40, 43).Add(shatterPos),
		image.Rect(80, 30, 90, 45).Add(shatterPos),
	}
	tiles := shatter.Shatter(c41.Bounds(), rocks...)
	for index, tile := range tiles {
		artist.DrawBounds (
			element.core,
			[]artist.Pattern {
				patterns.Uhex(0xFF0000FF),
				patterns.Uhex(0x00FF00FF),
				patterns.Uhex(0xFF00FFFF),
				patterns.Uhex(0xFFFF00FF),
				patterns.Uhex(0x00FFFFFF),
			} [index % 5], tile)
	}

	// 0, 2
	c02 := element.cellAt(0, 2)
	shapes.StrokeColorRectangle(c02, artist.Hex(0x888888FF), c02.Bounds(), 1)
	shapes.FillEllipse(c02, c41)

	// 1, 2
	c12 := element.cellAt(1, 2)
	shapes.StrokeColorRectangle(c12, artist.Hex(0x888888FF), c12.Bounds(), 1)
	shapes.StrokeEllipse(c12, c41, 5)
	
	// 2, 2
	c22 := element.cellAt(2, 2)
	shapes.FillRectangle(c22, c41)

	// 3, 2
	c32 := element.cellAt(3, 2)
	shapes.StrokeRectangle(c32, c41, 5)
	
	// how long did that take to render?
	drawTime := time.Since(drawStart)
	textDrawer := textdraw.Drawer { }
	textDrawer.SetFace(defaultfont.FaceRegular)
	textDrawer.SetText ([]rune (fmt.Sprintf (
		"%dms\n%dus",
		drawTime.Milliseconds(),
		drawTime.Microseconds())))
	textDrawer.Draw (
		element.core, artist.Hex(0xFFFFFFFF),
		image.Pt(bounds.Min.X + 8, bounds.Max.Y - 24))
}

func (element *Artist) colorLines (weight int, bounds image.Rectangle) {
	bounds = bounds.Inset(4)
	c := artist.Hex(0xFFFFFFFF)
	shapes.ColorLine(element.core, c, weight, bounds.Min, bounds.Max)
	shapes.ColorLine (
		element.core, c, weight,
		image.Pt(bounds.Max.X, bounds.Min.Y),
		image.Pt(bounds.Min.X, bounds.Max.Y))
	shapes.ColorLine (
		element.core, c, weight,
		image.Pt(bounds.Max.X, bounds.Min.Y + 16),
		image.Pt(bounds.Min.X, bounds.Max.Y - 16))
	shapes.ColorLine (
		element.core, c, weight,
		image.Pt(bounds.Min.X, bounds.Min.Y + 16),
		image.Pt(bounds.Max.X, bounds.Max.Y - 16))
	shapes.ColorLine (
		element.core, c, weight,
		image.Pt(bounds.Min.X + 20, bounds.Min.Y),
		image.Pt(bounds.Max.X - 20, bounds.Max.Y))
	shapes.ColorLine (
		element.core, c, weight,
		image.Pt(bounds.Max.X - 20, bounds.Min.Y),
		image.Pt(bounds.Min.X + 20, bounds.Max.Y))
	shapes.ColorLine (
		element.core, c, weight,
		image.Pt(bounds.Min.X, bounds.Min.Y + bounds.Dy() / 2),
		image.Pt(bounds.Max.X, bounds.Min.Y + bounds.Dy() / 2))
	shapes.ColorLine (
		element.core, c, weight,
		image.Pt(bounds.Min.X + bounds.Dx() / 2, bounds.Min.Y),
		image.Pt(bounds.Min.X + bounds.Dx() / 2, bounds.Max.Y))
}

func (element *Artist) cellAt (x, y int) (canvas.Canvas) {
	bounds := element.Bounds()
	cellBounds := image.Rectangle { }
	cellBounds.Min = bounds.Min
	cellBounds.Max.X = bounds.Min.X + bounds.Dx() / 5
	cellBounds.Max.Y = bounds.Min.Y + (bounds.Dy() - 48) / 4
	return canvas.Cut (element.core, cellBounds.Add (image.Pt (
		x * cellBounds.Dx(),
		y * cellBounds.Dy())))
}
