package elements

import "golang.org/x/image/math/fixed"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// Label is a simple text box.
type Label struct {
	entity tomo.FlexibleEntity
	
	align  textdraw.Align
	wrap   bool
	text   string
	drawer textdraw.Drawer

	forcedColumns int
	forcedRows    int
	minHeight     int
	
	config config.Wrapped
	theme  theme.Wrapped
}

// NewLabel creates a new label. If wrap is set to true, the text inside will be
// wrapped.
func NewLabel (text string, wrap bool) (element *Label) {
	element = &Label { }
	element.theme.Case = tomo.C("tomo", "label")
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.SetWrap(wrap)
	element.SetText(text)
	return
}

// Bind binds this element to an entity.
func (element *Label) Bind (entity tomo.Entity) {
	if entity == nil { element.entity = nil; return }
	element.entity = entity.(tomo.FlexibleEntity)
	element.updateMinimumSize()
}

// EmCollapse forces a minimum width and height upon the label. The width is
// measured in emspaces, and the height is measured in lines. If a zero value is
// given for a dimension, its minimum will be determined by the label's content.
// If the label's content is greater than these dimensions, it will be truncated
// to fit.
func (element *Label) EmCollapse (columns int, rows int) {
	element.forcedColumns = columns
	element.forcedRows    = rows
	if element.entity == nil { return }
	element.updateMinimumSize()
}

// FlexibleHeightFor returns the reccomended height for this element based on
// the given width in order to allow the text to wrap properly.
func (element *Label) FlexibleHeightFor (width int) (height int) {
	if element.wrap {
		return element.drawer.ReccomendedHeightFor(width)
	} else {
		return element.minHeight
	}
}

// SetText sets the label's text.
func (element *Label) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
	if element.entity == nil { return }
	element.updateMinimumSize()
	element.entity.Invalidate()
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
	if element.entity == nil { return }
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetAlign sets the alignment method of the label.
func (element *Label) SetAlign (align textdraw.Align) {
	if align == element.align { return }
	element.align = align
	element.drawer.SetAlign(align)
	if element.entity == nil { return }
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetTheme sets the element's theme.
func (element *Label) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	if element.entity == nil { return }
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *Label) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	if element.entity == nil { return }
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Label) Draw (destination canvas.Canvas) {
	if element.entity == nil { return }
	
	bounds := element.entity. Bounds()
	
	if element.wrap {
		element.drawer.SetMaxWidth(bounds.Dx())
		element.drawer.SetMaxHeight(bounds.Dy())
	}
	
	element.entity.DrawBackground(destination, bounds)

	textBounds := element.drawer.LayoutBounds()
	foreground := element.theme.Color (
		tomo.ColorForeground,
		tomo.State { })
	element.drawer.Draw(destination, foreground, bounds.Min.Sub(textBounds.Min))
}

func (element *Label) updateMinimumSize () {
	var width, height int
	
	if element.wrap {
		em := element.drawer.Em().Round()
		if em < 1 {
			em = element.theme.Padding(tomo.PatternBackground)[0]
		}
		width, height = em, element.drawer.LineHeight().Round()
		element.entity.NotifyFlexibleHeightChange()
	} else {
		bounds := element.drawer.LayoutBounds()
		width, height = bounds.Dx(), bounds.Dy()
	}

	if element.forcedColumns > 0 {
		width =
			element.drawer.Em().
			Mul(fixed.I(element.forcedColumns)).Floor()
	}

	if element.forcedRows > 0 {
		height =
			element.drawer.LineHeight().
			Mul(fixed.I(element.forcedRows)).Floor()
	}

	element.minHeight = height
	element.entity.SetMinimumSize(width, height)
}
