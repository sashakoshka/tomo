package testing

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/artist/shapes"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// Mouse is an element capable of testing mouse input. When the mouse is clicked
// and dragged on it, it draws a trail.
type Mouse struct {
	*core.Core
	core core.CoreControl
	drawing      bool
	lastMousePos image.Point
	
	config config.Wrapped
	theme  theme.Wrapped
}

// NewMouse creates a new mouse test element.
func NewMouse () (element *Mouse) {
	element = &Mouse { }
	element.theme.Case = tomo.C("tomo", "piano")
	element.Core, element.core = core.NewCore(element, element.draw)
	element.core.SetMinimumSize(32, 32)
	return
}

// SetTheme sets the element's theme.
func (element *Mouse) SetTheme (new tomo.Theme) {
	element.theme.Theme = new
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *Mouse) SetConfig (new tomo.Config) {
	element.config.Config = new
	element.redo()
}

func (element *Mouse) redo () {
	if !element.core.HasImage() { return }
	element.draw()
	element.core.DamageAll()
}

func (element *Mouse) draw () {
	bounds := element.Bounds()
	accent := element.theme.Color (
		tomo.ColorAccent,
		tomo.State { })
	shapes.FillColorRectangle(element.core, accent, bounds)
	shapes.StrokeColorRectangle (
		element.core,
		artist.Hex(0x000000FF),
		bounds, 1)
	shapes.ColorLine (
		element.core, artist.Hex(0xFFFFFFFF), 1,
		bounds.Min.Add(image.Pt(1, 1)),
		bounds.Min.Add(image.Pt(bounds.Dx() - 2, bounds.Dy() - 2)))
	shapes.ColorLine (
		element.core, artist.Hex(0xFFFFFFFF), 1,
		bounds.Min.Add(image.Pt(1, bounds.Dy() - 2)),
		bounds.Min.Add(image.Pt(bounds.Dx() - 2, 1)))
}

func (element *Mouse) HandleMouseDown (x, y int, button input.Button) {
	element.drawing = true
	element.lastMousePos = image.Pt(x, y)
}

func (element *Mouse) HandleMouseUp (x, y int, button input.Button) {
	element.drawing = false
	mousePos := image.Pt(x, y)
	element.core.DamageRegion (shapes.ColorLine (
		element.core, artist.Hex(0x000000FF), 1,
		element.lastMousePos, mousePos))
	element.lastMousePos = mousePos
}

func (element *Mouse) HandleMotion (x, y int) {
	if !element.drawing { return }
	mousePos := image.Pt(x, y)
	element.core.DamageRegion (shapes.ColorLine (
		element.core, artist.Hex(0x000000FF), 1,
		element.lastMousePos, mousePos))
	element.lastMousePos = mousePos
}
