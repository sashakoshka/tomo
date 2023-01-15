package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Label is a simple text box.
type Label struct {
	*core.Core
	core core.CoreControl

	wrap   bool
	text   string
	drawer artist.TextDrawer
}

// NewLabel creates a new label. If wrap is set to true, the text inside will be
// wrapped.
func NewLabel (text string, wrap bool) (element *Label) {
	element = &Label {  }
	element.Core, element.core = core.NewCore(element)
	face := theme.FontFaceRegular()
	element.drawer.SetFace(face)
	element.SetWrap(wrap)
	element.SetText(text)
	return
}

// Handle handles and event.
func (element *Label) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		if element.wrap {
			element.drawer.SetMaxWidth (resizeEvent.Width)
			element.drawer.SetMaxHeight(resizeEvent.Height)
		}
		element.draw()
	}
	return
}

// SetText sets the label's text.
func (element *Label) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText(text)
	element.updateMinimumSize()
	
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

// SetWrap sets wether or not the label's text wraps. If the text is set to
// wrap, the element will have a minimum size of a single character and
// automatically wrap its text. If the text is set to not wrap, the element will
// have a minimum size that fits its text.
func (element *Label) SetWrap (wrap bool) {
	if wrap == element.wrap { return }
	if !wrap {
		element.drawer.SetMaxWidth(0)
		element.drawer.SetMaxHeight(0)
	}
	element.wrap = wrap
	element.updateMinimumSize()
	
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

func (element *Label) updateMinimumSize () {
	if element.wrap {
		em := element.drawer.Em().Round()
		if em < 1 { em = theme.Padding() }
		element.core.SetMinimumSize (
			em, element.drawer.LineHeight().Round())
	} else {
		bounds := element.drawer.LayoutBounds()
		element.core.SetMinimumSize(bounds.Dx(), bounds.Dy())
	}
}

func (element *Label) draw () {
	bounds := element.core.Bounds()

	artist.FillRectangle (
		element.core,
		theme.BackgroundPattern(),
		bounds)

	textBounds := element.drawer.LayoutBounds()

	foreground := theme.ForegroundPattern(true)
	element.drawer.Draw (element.core, foreground, image.Point {
		X: 0 - textBounds.Min.X,
		Y: 0 - textBounds.Min.Y,
	})
}
