package testing

import "fmt"
import "time"
import "image"
import "image/color"
import "tomo"
import "art"
import "art/shatter"
import "tomo/textdraw"
import "art/shapes"
import "art/artutil"
import "art/patterns"

// Artist is an element that displays shapes and patterns drawn by the art
// package in order to test it.
type Artist struct {
	entity tomo.Entity
}

// NewArtist creates a new art test element.
func NewArtist () (element *Artist) {
	element = &Artist { }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.entity.SetMinimumSize(240, 240)
	return
}

func (element *Artist) Entity () tomo.Entity {
	return element.entity
}

func (element *Artist) Draw (destination art.Canvas) {
	bounds := element.entity.Bounds()
	patterns.Uhex(0x000000FF).Draw(destination, bounds)

	drawStart := time.Now()

	// 0, 0 - 3, 0
	for x := 0; x < 4; x ++ {
		element.colorLines(destination, x + 1, element.cellAt(destination, x, 0).Bounds())
	}

	// 4, 0
	c40 := element.cellAt(destination, 4, 0)
	shapes.StrokeColorRectangle(c40, artutil.Hex(0x888888FF), c40.Bounds(), 1)
	shapes.ColorLine (
		c40, artutil.Hex(0xFF0000FF), 1,
		c40.Bounds().Min, c40.Bounds().Max)

	// 0, 1
	c01 := element.cellAt(destination, 0, 1)
	shapes.StrokeColorRectangle(c01, artutil.Hex(0x888888FF), c01.Bounds(), 1)
	shapes.FillColorEllipse(destination, artutil.Hex(0x00FF00FF), c01.Bounds())

	// 1, 1 - 3, 1
	for x := 1; x < 4; x ++ {
		c := element.cellAt(destination, x, 1)
		shapes.StrokeColorRectangle (
			destination, artutil.Hex(0x888888FF),
			c.Bounds(), 1)
		shapes.StrokeColorEllipse (
			destination,
			[]color.RGBA {
				artutil.Hex(0xFF0000FF),
				artutil.Hex(0x00FF00FF),
				artutil.Hex(0xFF00FFFF),
			} [x - 1],
			c.Bounds(), x)
	}

	// 4, 1
	c41 := element.cellAt(destination, 4, 1)
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
		[]art.Pattern {
			patterns.Uhex(0xFF0000FF),
			patterns.Uhex(0x00FF00FF),
			patterns.Uhex(0xFF00FFFF),
			patterns.Uhex(0xFFFF00FF),
			patterns.Uhex(0x00FFFFFF),
		} [index % 5].Draw(destination, tile)
	}

	// 0, 2
	c02 := element.cellAt(destination, 0, 2)
	shapes.StrokeColorRectangle(c02, artutil.Hex(0x888888FF), c02.Bounds(), 1)
	shapes.FillEllipse(c02, c41, c02.Bounds())

	// 1, 2
	c12 := element.cellAt(destination, 1, 2)
	shapes.StrokeColorRectangle(c12, artutil.Hex(0x888888FF), c12.Bounds(), 1)
	shapes.StrokeEllipse(c12, c41, c12.Bounds(), 5)
	
	// 2, 2
	c22 := element.cellAt(destination, 2, 2)
	shapes.FillRectangle(c22, c41, c22.Bounds())

	// 3, 2
	c32 := element.cellAt(destination, 3, 2)
	shapes.StrokeRectangle(c32, c41, c32.Bounds(), 5)
	
	// 4, 2
	c42 := element.cellAt(destination, 4, 2)
	
	// 0, 3
	c03 := element.cellAt(destination, 0, 3)
	patterns.Border {
		Canvas: element.thingy(c42),
		Inset:  art.Inset { 8, 8, 8, 8 },
	}.Draw(c03, c03.Bounds())
	
	// 1, 3
	c13 := element.cellAt(destination, 1, 3)
	patterns.Border {
		Canvas: element.thingy(c42),
		Inset:  art.Inset { 8, 8, 8, 8 },
	}.Draw(c13, c13.Bounds().Inset(10))
	
	// 2, 3
	c23 := element.cellAt(destination, 2, 3)
	patterns.Border {
		Canvas: element.thingy(c42),
		Inset:  art.Inset { 8, 8, 8, 8 },
	}.Draw(c23, c23.Bounds())
	patterns.Border {
		Canvas: element.thingy(c42),
		Inset:  art.Inset { 8, 8, 8, 8 },
	}.Draw(art.Cut(c23, c23.Bounds().Inset(16)), c23.Bounds())
	
	// how long did that take to render?
	drawTime := time.Since(drawStart)
	textDrawer := textdraw.Drawer { }
	textDrawer.SetFace(element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal,
		tomo.C("tomo", "art")))
	textDrawer.SetText ([]rune (fmt.Sprintf (
		"%dms\n%dus",
		drawTime.Milliseconds(),
		drawTime.Microseconds())))
	textDrawer.Draw (
		destination, artutil.Hex(0xFFFFFFFF),
		image.Pt(bounds.Min.X + 8, bounds.Max.Y - 24))
}

func (element *Artist) colorLines (destination art.Canvas, weight int, bounds image.Rectangle) {
	bounds = bounds.Inset(4)
	c := artutil.Hex(0xFFFFFFFF)
	shapes.ColorLine(destination, c, weight, bounds.Min, bounds.Max)
	shapes.ColorLine (
		destination, c, weight,
		image.Pt(bounds.Max.X, bounds.Min.Y),
		image.Pt(bounds.Min.X, bounds.Max.Y))
	shapes.ColorLine (
		destination, c, weight,
		image.Pt(bounds.Max.X, bounds.Min.Y + 16),
		image.Pt(bounds.Min.X, bounds.Max.Y - 16))
	shapes.ColorLine (
		destination, c, weight,
		image.Pt(bounds.Min.X, bounds.Min.Y + 16),
		image.Pt(bounds.Max.X, bounds.Max.Y - 16))
	shapes.ColorLine (
		destination, c, weight,
		image.Pt(bounds.Min.X + 20, bounds.Min.Y),
		image.Pt(bounds.Max.X - 20, bounds.Max.Y))
	shapes.ColorLine (
		destination, c, weight,
		image.Pt(bounds.Max.X - 20, bounds.Min.Y),
		image.Pt(bounds.Min.X + 20, bounds.Max.Y))
	shapes.ColorLine (
		destination, c, weight,
		image.Pt(bounds.Min.X, bounds.Min.Y + bounds.Dy() / 2),
		image.Pt(bounds.Max.X, bounds.Min.Y + bounds.Dy() / 2))
	shapes.ColorLine (
		destination, c, weight,
		image.Pt(bounds.Min.X + bounds.Dx() / 2, bounds.Min.Y),
		image.Pt(bounds.Min.X + bounds.Dx() / 2, bounds.Max.Y))
}

func (element *Artist) cellAt (destination art.Canvas, x, y int) (art.Canvas) {
	bounds := element.entity.Bounds()
	cellBounds := image.Rectangle { }
	cellBounds.Min = bounds.Min
	cellBounds.Max.X = bounds.Min.X + bounds.Dx() / 5
	cellBounds.Max.Y = bounds.Min.Y + (bounds.Dy() - 48) / 4
	return art.Cut (destination, cellBounds.Add (image.Pt (
		x * cellBounds.Dx(),
		y * cellBounds.Dy())))
}

func (element *Artist) thingy (destination art.Canvas) (result art.Canvas) {
	bounds := destination.Bounds()
	bounds = image.Rect(0, 0, 32, 32).Add(bounds.Min)
	shapes.FillColorRectangle(destination, artutil.Hex(0x440000FF), bounds)
	shapes.StrokeColorRectangle(destination, artutil.Hex(0xFF0000FF), bounds, 1)
	shapes.StrokeColorRectangle(destination, artutil.Hex(0x004400FF), bounds.Inset(4), 1)
	shapes.FillColorRectangle(destination, artutil.Hex(0x004444FF), bounds.Inset(12))
	shapes.StrokeColorRectangle(destination, artutil.Hex(0x888888FF), bounds.Inset(8), 1)
	return art.Cut(destination, bounds)
}
