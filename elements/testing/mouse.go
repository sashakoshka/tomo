package testing

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Mouse is an element capable of testing mouse input. When the mouse is clicked
// and dragged on it, it draws a trail.
type Mouse struct {
	*core.Core
	core core.CoreControl
	drawing      bool
	color        artist.Pattern
	lastMousePos image.Point
}

// NewMouse creates a new mouse test element.
func NewMouse () (element *Mouse) {
	element = &Mouse { }
	element.Core, element.core = core.NewCore(element)
	element.core.SetMinimumSize(32, 32)
	element.color = artist.NewUniform(color.Black)
	return
}

func (element *Mouse) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		artist.FillRectangle (
			element.core,
			theme.AccentImage(),
			element.Bounds())
		artist.StrokeRectangle (
			element.core,
			artist.NewUniform(color.Black), 1,
			element.Bounds())
		artist.Line (
			element.core, artist.NewUniform(color.White), 3,
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
		if !element.drawing { return }
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
