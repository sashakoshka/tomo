package testing

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist/shapes"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// Mouse is an element capable of testing mouse input. When the mouse is clicked
// and dragged on it, it draws a trail.
type Mouse struct {
	entity       tomo.Entity
	pressed      bool
	lastMousePos image.Point
	
	config config.Wrapped
	theme  theme.Wrapped
}

// NewMouse creates a new mouse test element.
func NewMouse () (element *Mouse) {
	element = &Mouse { }
	element.theme.Case = tomo.C("tomo", "mouse")
	element.entity = tomo.NewEntity(element)
	element.entity.SetMinimumSize(32, 32)
	return
}

func (element *Mouse) Entity () tomo.Entity {
	return element.entity
}

func (element *Mouse) Draw (destination canvas.Canvas) {
	bounds := element.entity.Bounds()
	accent := element.theme.Color (
		tomo.ColorAccent,
		tomo.State { })
	shapes.FillColorRectangle(destination, accent, bounds)
	shapes.StrokeColorRectangle (
		destination,
		artist.Hex(0x000000FF),
		bounds, 1)
	shapes.ColorLine (
		destination, artist.Hex(0xFFFFFFFF), 1,
		bounds.Min.Add(image.Pt(1, 1)),
		bounds.Min.Add(image.Pt(bounds.Dx() - 2, bounds.Dy() - 2)))
	shapes.ColorLine (
		destination, artist.Hex(0xFFFFFFFF), 1,
		bounds.Min.Add(image.Pt(1, bounds.Dy() - 2)),
		bounds.Min.Add(image.Pt(bounds.Dx() - 2, 1)))
	if element.pressed {
		midpoint := bounds.Min.Add(bounds.Max.Sub(bounds.Min).Div(2))
		shapes.ColorLine (
			destination, artist.Hex(0x000000FF), 1,
			midpoint, element.lastMousePos)
	}
}

// SetTheme sets the element's theme.
func (element *Mouse) SetTheme (new tomo.Theme) {
	element.theme.Theme = new
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *Mouse) SetConfig (new tomo.Config) {
	element.config.Config = new
	element.entity.Invalidate()
}

func (element *Mouse) HandleMouseDown (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	element.pressed = true
	element.lastMousePos = position
	element.entity.Invalidate()
}

func (element *Mouse) HandleMouseUp (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	element.pressed = false
	element.entity.Invalidate()
}

func (element *Mouse) HandleMotion (position image.Point) {
	if !element.pressed { return }
	element.lastMousePos = position
	element.entity.Invalidate()
}
