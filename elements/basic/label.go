package basicElements

import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

var labelCase = theme.C("basic", "label")

// Label is a simple text box.
type Label struct {
	*core.Core
	core core.CoreControl

	wrap   bool
	text   string
	drawer artist.TextDrawer
	
	onFlexibleHeightChange func ()
}

// NewLabel creates a new label. If wrap is set to true, the text inside will be
// wrapped.
func NewLabel (text string, wrap bool) (element *Label) {
	element = &Label {  }
	element.Core, element.core = core.NewCore(element.handleResize)
	face := theme.FontFaceRegular()
	element.drawer.SetFace(face)
	element.SetWrap(wrap)
	element.SetText(text)
	return
}

func (element *Label) handleResize () {
	bounds := element.Bounds()
	if element.wrap {
		element.drawer.SetMaxWidth(bounds.Dx())
		element.drawer.SetMaxHeight(bounds.Dy())
	}
	element.draw()
	return
}

// FlexibleHeightFor returns the reccomended height for this element based on
// the given width in order to allow the text to wrap properly.
func (element *Label) FlexibleHeightFor (width int) (height int) {
	if element.wrap {
		return element.drawer.ReccomendedHeightFor(width)
	} else {
		_, height = element.MinimumSize()
		return
	}
}

// OnFlexibleHeightChange sets a function to be called when the parameters
// affecting this element's flexible height are changed.
func (element *Label) OnFlexibleHeightChange (callback func ()) {
	element.onFlexibleHeightChange = callback
}

// SetText sets the label's text.
func (element *Label) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
	element.updateMinimumSize()
	
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
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
		element.core.DamageAll()
	}
}

func (element *Label) updateMinimumSize () {
	if element.wrap {
		em := element.drawer.Em().Round()
		if em < 1 { em = theme.Padding() }
		element.core.SetMinimumSize (
			em, element.drawer.LineHeight().Round())
		if element.onFlexibleHeightChange != nil {
			element.onFlexibleHeightChange()
		}
	} else {
		bounds := element.drawer.LayoutBounds()
		element.core.SetMinimumSize(bounds.Dx(), bounds.Dy())
	}
}

func (element *Label) draw () {
	bounds := element.Bounds()

	pattern, _ := theme.BackgroundPattern(theme.PatternState {
		Case: labelCase,
	})
	artist.FillRectangle(element, pattern, bounds)

	textBounds := element.drawer.LayoutBounds()

	foreground, _ := theme.ForegroundPattern (theme.PatternState {
		Case: labelCase,
	})
	element.drawer.Draw (element, foreground, bounds.Min.Sub(textBounds.Min))
}
