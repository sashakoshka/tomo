package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/textmanip"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// TextBox is a single-line text input.
type TextBox struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl
	
	cursor int
	scroll int
	placeholder string
	text        []rune
	
	placeholderDrawer artist.TextDrawer
	valueDrawer       artist.TextDrawer
	
	theme  theme.Theme
	config config.Config
	c theme.Case
	
	onKeyDown func (key input.Key, modifiers input.Modifiers) (handled bool)
	onChange  func ()
	onScrollBoundsChange func ()
}

// NewTextBox creates a new text box with the specified placeholder text, and
// a value. When the value is empty, the placeholder will be displayed in gray
// text.
func NewTextBox (placeholder, value string) (element *TextBox) {
	element = &TextBox { c: theme.C("basic", "textBox") }
	element.Core, element.core = core.NewCore(element.handleResize)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore (func () {
		if element.core.HasImage () {
			element.draw()
			element.core.DamageAll()
		}
	})
	element.placeholder = placeholder
	element.placeholderDrawer.SetText([]rune(placeholder))
	element.updateMinimumSize()
	element.SetValue(value)
	return
}

func (element *TextBox) handleResize () {
	element.scrollToCursor()
	element.draw()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *TextBox) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	if !element.Focused() { element.Focus() }
}

func (element *TextBox) HandleMouseUp (x, y int, button input.Button) { }
func (element *TextBox) HandleMouseMove (x, y int) { }
func (element *TextBox) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *TextBox) HandleKeyDown(key input.Key, modifiers input.Modifiers) {
	if element.onKeyDown != nil && element.onKeyDown(key, modifiers) {
		return
	}

	scrollMemory := element.scroll
	altered     := true
	textChanged := false
	switch {
	case key == input.KeyBackspace:
		if len(element.text) < 1 { break }
		element.text, element.cursor = textmanip.Backspace (
			element.text,
			element.cursor,
			modifiers.Control)
		textChanged = true
			
	case key == input.KeyDelete:
		if len(element.text) < 1 { break }
		element.text, element.cursor = textmanip.Delete (
			element.text,
			element.cursor,
			modifiers.Control)
		textChanged = true
			
	case key == input.KeyLeft:
		element.cursor = textmanip.MoveLeft (
			element.text,
			element.cursor,
			modifiers.Control)
			
	case key == input.KeyRight:
		element.cursor = textmanip.MoveRight (
			element.text,
			element.cursor,
			modifiers.Control)
		
	case key.Printable():
		element.text, element.cursor = textmanip.Type (
			element.text,
			element.cursor,
			rune(key))
		textChanged = true
			
	default:
		altered = false
	}

	if textChanged {
		element.runOnChange()
		element.valueDrawer.SetText(element.text)
	}

	if altered {
		element.scrollToCursor()
	}

	if (textChanged || scrollMemory != element.scroll) &&
		element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
	
	if altered {
		element.redo()
	}
}

func (element *TextBox) HandleKeyUp(key input.Key, modifiers input.Modifiers) { }

func (element *TextBox) SetPlaceholder (placeholder string) {
	if element.placeholder == placeholder { return }
	
	element.placeholder = placeholder
	element.placeholderDrawer.SetText([]rune(placeholder))
	
	element.updateMinimumSize()
	element.redo()
}

func (element *TextBox) SetValue (text string) {
	// if element.text == text { return }

	element.text = []rune(text)
	element.runOnChange()
	element.valueDrawer.SetText(element.text)
	if element.cursor > element.valueDrawer.Length() {
		element.cursor = element.valueDrawer.Length()
	}
	element.scrollToCursor()
	element.redo()
}

func (element *TextBox) Value () (value string) {
	return string(element.text)
}

func (element *TextBox) Filled () (filled bool) {
	return len(element.text) > 0
}

func (element *TextBox) OnKeyDown (
	callback func (key input.Key, modifiers input.Modifiers) (handled bool),
) {
	element.onKeyDown = callback
}

func (element *TextBox) OnChange (callback func ()) {
	element.onChange = callback
}

// ScrollContentBounds returns the full content size of the element.
func (element *TextBox) ScrollContentBounds () (bounds image.Rectangle) {
	bounds = element.valueDrawer.LayoutBounds()
	return bounds.Sub(bounds.Min)
}

// ScrollViewportBounds returns the size and position of the element's viewport
// relative to ScrollBounds.
func (element *TextBox) ScrollViewportBounds () (bounds image.Rectangle) {
	return image.Rect (
		element.scroll,
		0,
		element.scroll + element.scrollViewportWidth(),
		0)
}

func (element *TextBox) scrollViewportWidth () (width int) {
	return element.Bounds().Inset(element.config.Padding()).Dx()
}

// ScrollTo scrolls the viewport to the specified point relative to
// ScrollBounds.
func (element *TextBox) ScrollTo (position image.Point) {
	// constrain to minimum
	element.scroll = position.X
	if element.scroll < 0 { element.scroll = 0 }
	
	// constrain to maximum
	contentBounds := element.ScrollContentBounds()
	maxPosition   := contentBounds.Max.X - element.scrollViewportWidth()
	if element.scroll > maxPosition { element.scroll = maxPosition }

	element.redo()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

// ScrollAxes returns the supported axes for scrolling.
func (element *TextBox) ScrollAxes () (horizontal, vertical bool) {
	return true, false
}

func (element *TextBox) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
}

func (element *TextBox) runOnChange () {
	if element.onChange != nil {
		element.onChange()
	}
}

func (element *TextBox) scrollToCursor () {
	if !element.core.HasImage() { return }

	bounds := element.Bounds().Inset(element.config.Padding())
	bounds = bounds.Sub(bounds.Min)
	bounds.Max.X -= element.valueDrawer.Em().Round()
	cursorPosition := element.valueDrawer.PositionOf(element.cursor)
	cursorPosition.X -= element.scroll
	maxX := bounds.Max.X
	minX := maxX
	if cursorPosition.X > maxX {
		element.scroll += cursorPosition.X - maxX
	} else if cursorPosition.X < minX {
		element.scroll -= minX - cursorPosition.X
		if element.scroll < 0 { element.scroll = 0 }
	}
}

// SetTheme sets the element's theme.
func (element *TextBox) SetTheme (new theme.Theme) {
	element.theme = new
	face := element.theme.FontFace (
		theme.FontStyleRegular,
		theme.FontSizeNormal,
		element.c)
	element.placeholderDrawer.SetFace(face)
	element.valueDrawer.SetFace(face)
	element.updateMinimumSize()
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *TextBox) SetConfig (new config.Config) {
	element.config = new
	element.updateMinimumSize()
	element.redo()
}

func (element *TextBox) updateMinimumSize () {
	textBounds := element.placeholderDrawer.LayoutBounds()
	inset := element.theme.Inset(theme.PatternInput, element.c)
	element.core.SetMinimumSize (
		textBounds.Dx() +
		element.config.Padding() * 2 + inset[3] + inset[1],
		element.placeholderDrawer.LineHeight().Round() +
		element.config.Padding() * 2 + inset[0] + inset[2])
}

func (element *TextBox) redo () {
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *TextBox) draw () {
	bounds := element.Bounds()

	// FIXME: take index into account
	state := theme.PatternState {
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
	}
	pattern := element.theme.Pattern(theme.PatternSunken, element.c, state)
	artist.FillRectangle(element, pattern, bounds)

	if len(element.text) == 0 && !element.Focused() {
		// draw placeholder
		textBounds := element.placeholderDrawer.LayoutBounds()
		offset := bounds.Min.Add (image.Point {
			X: element.config.Padding(),
			Y: element.config.Padding(),
		})
		foreground := element.theme.Pattern (
			theme.PatternForeground, element.c,
			theme.PatternState { Disabled: true })
		element.placeholderDrawer.Draw (
			element,
			foreground,
			offset.Sub(textBounds.Min))
	} else {
		// draw input value
		textBounds := element.valueDrawer.LayoutBounds()
		offset := bounds.Min.Add (image.Point {
			X: element.config.Padding() - element.scroll,
			Y: element.config.Padding(),
		})
		foreground := element.theme.Pattern (
			theme.PatternForeground, element.c, state)
		element.valueDrawer.Draw (
			element,
			foreground,
			offset.Sub(textBounds.Min))

		if element.Focused() {
			// cursor
			cursorPosition := element.valueDrawer.PositionOf (
				element.cursor)
			artist.Line (
				element,
				foreground, 1,
				cursorPosition.Add(offset),
				image.Pt (
					cursorPosition.X,
					cursorPosition.Y + element.valueDrawer.
					LineHeight().Round()).Add(offset))
		}
	}
}
