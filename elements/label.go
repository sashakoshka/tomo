package elements

import "image"
import "golang.org/x/image/math/fixed"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/input"
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

// NewLabel creates a new label.
func NewLabel (text string) (element *Label) {
	element = &Label { }
	element.theme.Case = tomo.C("tomo", "label")
	element.entity = tomo.NewEntity(element).(tomo.FlexibleEntity)
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.SetText(text)
	return
}

// NewLabelWrapped creates a new label with text wrapping on.
func NewLabelWrapped (text string) (element *Label) {
	element = NewLabel(text)
	element.SetWrap(true)
	return
}

// Entity returns this element's entity.
func (element *Label) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Label) Draw (destination canvas.Canvas) {
	bounds := element.entity.Bounds()
	
	if element.wrap {
		element.drawer.SetMaxWidth(bounds.Dx())
		element.drawer.SetMaxHeight(bounds.Dy())
	}
	
	element.entity.DrawBackground(destination)

	textBounds := element.drawer.LayoutBounds()
	foreground := element.theme.Color (
		tomo.ColorForeground,
		tomo.State { })
	element.drawer.Draw(destination, foreground, bounds.Min.Sub(textBounds.Min))
}

// Copy copies the label's textto the clipboard.
func (element *Label) Copy () {
	window := element.entity.Window()
	if window != nil {
		window.Copy(data.Bytes(data.MimePlain, []byte(element.text)))
	}
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
		return element.minHeight
	}
}

// SetText sets the label's text.
func (element *Label) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
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
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetAlign sets the alignment method of the label.
func (element *Label) SetAlign (align textdraw.Align) {
	if align == element.align { return }
	element.align = align
	element.drawer.SetAlign(align)
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
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *Label) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Label) HandleMouseDown  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if button == input.ButtonRight {
		element.contextMenu(position)
	}
}

func (element *Label) HandleMouseUp  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) { }

func (element *Label) contextMenu (position image.Point) {
	window := element.entity.Window()
	menu, err := window.NewMenu(image.Rectangle { position, position })
	if err != nil { return }

	closeAnd := func (callback func ()) func () {
		return func () { callback(); menu.Close() }
	}
	
	copyButton := NewButton("Copy")
	copyButton.ShowText(false)
	copyButton.SetIcon(tomo.IconCopy)
	copyButton.OnClick(closeAnd(element.Copy))

	menu.Adopt (NewHBox (
		SpaceNone,
		copyButton,
	))
	copyButton.Focus()
	menu.Show()
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
