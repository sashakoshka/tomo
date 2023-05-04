package testing

import "image"
import "tomo"
import "tomo/input"
import "art"
import "art/shapes"
import "art/artutil"

var mouseCase = tomo.C("tomo", "mouse")

// Mouse is an element capable of testing mouse input. When the mouse is clicked
// and dragged on it, it draws a trail.
type Mouse struct {
	entity       tomo.Entity
	pressed      bool
	lastMousePos image.Point
}

// NewMouse creates a new mouse test element.
func NewMouse () (element *Mouse) {
	element = &Mouse { }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.entity.SetMinimumSize(32, 32)
	return
}

func (element *Mouse) Entity () tomo.Entity {
	return element.entity
}

func (element *Mouse) Draw (destination art.Canvas) {
	bounds := element.entity.Bounds()
	accent := element.entity.Theme().Color (
		tomo.ColorAccent,
		tomo.State { },
		mouseCase)
	shapes.FillColorRectangle(destination, accent, bounds)
	shapes.StrokeColorRectangle (
		destination,
		artutil.Hex(0x000000FF),
		bounds, 1)
	shapes.ColorLine (
		destination, artutil.Hex(0xFFFFFFFF), 1,
		bounds.Min.Add(image.Pt(1, 1)),
		bounds.Min.Add(image.Pt(bounds.Dx() - 2, bounds.Dy() - 2)))
	shapes.ColorLine (
		destination, artutil.Hex(0xFFFFFFFF), 1,
		bounds.Min.Add(image.Pt(1, bounds.Dy() - 2)),
		bounds.Min.Add(image.Pt(bounds.Dx() - 2, 1)))
	if element.pressed {
		midpoint := bounds.Min.Add(bounds.Max.Sub(bounds.Min).Div(2))
		shapes.ColorLine (
			destination, artutil.Hex(0x000000FF), 1,
			midpoint, element.lastMousePos)
	}
}

func (element *Mouse) HandleThemeChange (new tomo.Theme) {
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
