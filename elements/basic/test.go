package basic

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Test is a simple element that can be used as a placeholder.
type Test struct {
	*core.Core
	core core.CoreControl
	drawing      bool
	color        tomo.Image
	lastMousePos image.Point
}

// NewTest creates a new test element.
func NewTest () (element *Test) {
	element = &Test { }
	element.Core, element.core = core.NewCore(element)
	element.core.SetMinimumSize(32, 32)
	element.color = artist.NewUniform(color.Black)
	return
}

func (element *Test) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		artist.Rectangle (
			element.core,
			theme.AccentImage(),
			artist.NewUniform(color.Black),
			1, element.Bounds())
		artist.Line (
			element.core, artist.NewUniform(color.White), 1,
			image.Pt(1, 1),
			image.Pt(resizeEvent.Width - 2, resizeEvent.Height - 2))
		artist.Line (
			element.core, artist.NewUniform(color.White), 1,
			image.Pt(1, resizeEvent.Height - 2),
			image.Pt(resizeEvent.Width - 2, 1))
	
	case tomo.EventMouseDown:
		element.drawing = true
		mouseDownEvent := event.(tomo.EventMouseDown)
		element.lastMousePos = image.Pt (
			mouseDownEvent.X,
			mouseDownEvent.Y)

	case tomo.EventMouseUp:
		element.drawing = false
		mouseUpEvent := event.(tomo.EventMouseUp)
		mousePos := image.Pt (
			mouseUpEvent.X,
			mouseUpEvent.Y)
		element.core.PushRegion (artist.Line (
			element.core, element.color, 1,
			element.lastMousePos, mousePos))
		element.lastMousePos = mousePos

	case tomo.EventMouseMove:
		mouseMoveEvent := event.(tomo.EventMouseMove)
		mousePos := image.Pt (
			mouseMoveEvent.X,
			mouseMoveEvent.Y)
		element.core.PushRegion (artist.Line (
			element.core, element.color, 1,
			element.lastMousePos, mousePos))
		element.lastMousePos = mousePos
	}
	return
}
