package basicElements

import "golang.org/x/image/math/fixed"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Label is a simple text box.
type Label struct {
	*core.Core
	core core.CoreControl

	wrap   bool
	text   string
	drawer textdraw.Drawer

	forcedColumns int
	forcedRows    int
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onFlexibleHeightChange func ()
}

// NewLabel creates a new label. If wrap is set to true, the text inside will be
// wrapped.
func NewLabel (text string, wrap bool) (element *Label) {
	element = &Label { }
	element.theme.Case = theme.C("basic", "label")
	element.Core, element.core = core.NewCore(element, element.handleResize)
	element.SetWrap(wrap)
	element.SetText(text)
	return
}

func (element *Label) redo () {
	face := element.theme.FontFace (
		theme.FontStyleRegular,
		theme.FontSizeNormal)
	element.drawer.SetFace(face)
	element.updateMinimumSize()
	bounds := element.Bounds()
	if element.wrap {
		element.drawer.SetMaxWidth(bounds.Dx())
		element.drawer.SetMaxHeight(bounds.Dy())
	}
	element.draw()
	element.core.DamageAll()
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

// EmCollapse forces a minimum width and height upon the label. The width is
// measured in emspaces, and the height is measured in lines. If a zero value is
// given for a dimension, its minimum will be determined by the label's content.
// If the label's content is greater than these dimensions, it will be truncated
// to fit.
func (element *Label) EmCollapse (columns int, rows int) {
	element.forcedColumns = columns
	element.forcedRows    = rows
	element.updateMinimumSize()
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

// SetTheme sets the element's theme.
func (element *Label) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawer.SetFace (element.theme.FontFace (
		theme.FontStyleRegular,
		theme.FontSizeNormal))
	element.updateMinimumSize()
	
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

// SetConfig sets the element's configuration.
func (element *Label) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Label) updateMinimumSize () {
	var width, height int
	
	if element.wrap {
		em := element.drawer.Em().Round()
		if em < 1 {
			em = element.theme.Padding(theme.PatternBackground)[0]
		}
		width, height = em, element.drawer.LineHeight().Round()
		if element.onFlexibleHeightChange != nil {
			element.onFlexibleHeightChange()
		}
	} else {
		bounds := element.drawer.LayoutBounds()
		width, height = bounds.Dx(), bounds.Dy()
	}

	if element.forcedColumns > 0 {
		width = int (
			element.drawer.Em().
			Mul(fixed.I(element.forcedColumns)))
	}

	if element.forcedRows > 0 {
		height = int (
			element.drawer.LineHeight().
			Mul(fixed.I(element.forcedRows)))
	}

	element.core.SetMinimumSize(width, height)
}

func (element *Label) draw () {
	bounds := element.Bounds()

	pattern := element.theme.Pattern (
		theme.PatternBackground,
		theme.State { })
	pattern.Draw(element.core, bounds)

	textBounds := element.drawer.LayoutBounds()

	foreground := element.theme.Color (
		theme.ColorForeground,
		theme.State { })
	element.drawer.Draw(element.core, foreground, bounds.Min.Sub(textBounds.Min))
}
