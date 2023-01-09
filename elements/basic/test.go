package basic

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Test is a simple element that can be used as a placeholder.
type Test struct {
	*Core
	core CoreControl
}

// NewTest creates a new test element.
func NewTest () (element *Test) {
	element = &Test { }
	element.Core, element.core = NewCore(element)
	element.core.SetMinimumSize(32, 32)
	return
}

func (element *Test) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		for y := 0; y < resizeEvent.Height; y ++ {
		for x := 0; x < resizeEvent.Width;  x ++ {
			pixel := color.RGBA {
				R: 0x40, G: 0x80, B: 0x90, A: 0xFF,
			}
			element.core.SetRGBA (x, y, pixel)
		}}
		artist.Line (
			element.core, artist.NewUniform(color.White), 1,
			image.Pt(0, 0),
			image.Pt(resizeEvent.Width, resizeEvent.Height))
		artist.Line (
			element.core, artist.NewUniform(color.White), 1,
			image.Pt(0, resizeEvent.Height),
			image.Pt(resizeEvent.Width, 0))
	
	default:
	}
	return
}

func (element *Test) Selectable () (selectable bool) {
	return
}

func (element *Test) AdvanceSelection (direction int) (ok bool) {
	return
}
