package testing

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/config"
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
	
	config config.Config
	theme  theme.Theme
	c      theme.Case
}

// NewMouse creates a new mouse test element.
func NewMouse () (element *Mouse) {
	element = &Mouse { c: theme.C("testing", "mouse") }
	element.Core, element.core = core.NewCore(element.draw)
	element.core.SetMinimumSize(32, 32)
	element.color = artist.NewUniform(color.Black)
	return
}

// SetTheme sets the element's theme.
func (element *Mouse) SetTheme (new theme.Theme) {
	element.theme = new
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *Mouse) SetConfig (new config.Config) {
	element.config = new
	element.redo()
}

func (element *Mouse) redo () {
	if !element.core.HasImage() { return }
	element.draw()
	element.core.DamageAll()
}

func (element *Mouse) draw () {
	bounds := element.Bounds()
	pattern := element.theme.Pattern (
		theme.PatternAccent,
		element.c,
		theme.PatternState { })
	artist.FillRectangle(element, pattern, bounds)
	artist.StrokeRectangle (
		element,
		artist.NewUniform(color.Black), 1,
		bounds)
	artist.Line (
		element, artist.NewUniform(color.White), 1,
		bounds.Min.Add(image.Pt(1, 1)),
		bounds.Min.Add(image.Pt(bounds.Dx() - 2, bounds.Dy() - 2)))
	artist.Line (
		element, artist.NewUniform(color.White), 1,
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
	element.core.DamageRegion (artist.Line (
		element, element.color, 1,
		element.lastMousePos, mousePos))
	element.lastMousePos = mousePos
}

func (element *Mouse) HandleMouseMove (x, y int) {
	if !element.drawing { return }
	mousePos := image.Pt(x, y)
	element.core.DamageRegion (artist.Line (
		element, element.color, 1,
		element.lastMousePos, mousePos))
	element.lastMousePos = mousePos
}

func (element *Mouse) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }
