package elements

import "io"
import "time"
import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/textmanip"
import "git.tebibyte.media/sashakoshka/tomo/fixedutil"
import "git.tebibyte.media/sashakoshka/tomo/artist/shapes"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// TextBox is a single-line text input.
type TextBox struct {
	lastClick   time.Time
	dragging    int
	dot         textmanip.Dot
	scroll      int
	placeholder string
	text        []rune
	
	placeholderDrawer textdraw.Drawer
	valueDrawer       textdraw.Drawer
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onKeyDown func (key input.Key, modifiers input.Modifiers) (handled bool)
	onChange  func ()
	onEnter   func ()
	onScrollBoundsChange func ()
}

// NewTextBox creates a new text box with the specified placeholder text, and
// a value. When the value is empty, the placeholder will be displayed in gray
// text.
func NewTextBox (placeholder, value string) (element *TextBox) {
	element = &TextBox { }
	element.theme.Case = tomo.C("tomo", "textBox")
	element.Core, element.core = core.NewCore(element, element.handleResize)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore (element.core, func () {
		if element.core.HasImage () {
			element.draw()
			element.core.DamageAll()
		}
	})
	element.placeholder = placeholder
	element.placeholderDrawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.valueDrawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.placeholderDrawer.SetText([]rune(placeholder))
	element.updateMinimumSize()
	element.SetValue(value)
	return
}

func (element *TextBox) handleResize () {
	element.scrollToCursor()
	element.draw()
	if parent, ok := element.core.Parent().(tomo.ScrollableParent); ok {
		parent.NotifyScrollBoundsChange(element)
	}
}

func (element *TextBox) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	if !element.Focused() { element.Focus() }

	if button == input.ButtonLeft {
		runeIndex := element.atPosition(image.Pt(x, y))
		if runeIndex == -1 { return }
		
		if time.Since(element.lastClick) < element.config.DoubleClickDelay() {
			element.dragging = 2
			element.dot = textmanip.WordAround(element.text, runeIndex)
		} else {
			element.dragging = 1
			element.dot = textmanip.EmptyDot(runeIndex)
			element.lastClick = time.Now()
		}
		
		element.redo()
	}
}

func (element *TextBox) HandleMotion (x, y int) {
	if !element.Enabled() { return }

	switch element.dragging {
	case 1:
		runeIndex := element.atPosition(image.Pt(x, y))
		if runeIndex > -1 {
			element.dot.End = runeIndex
			element.redo()
		}
		
	case 2:
		runeIndex := element.atPosition(image.Pt(x, y))
		if runeIndex > -1 {
			if runeIndex < element.dot.Start {
				element.dot.End =
					runeIndex -
					textmanip.WordToLeft (
						element.text,
						runeIndex)
			} else {
				element.dot.End =
					runeIndex +
					textmanip.WordToRight (
						element.text,
						runeIndex)
			}
			element.redo()
		}
	}
}

func (element *TextBox) textOffset () image.Point {
	padding     := element.theme.Padding(tomo.PatternInput)
	bounds      := element.Bounds()
	innerBounds := padding.Apply(bounds)
	textHeight  := element.valueDrawer.LineHeight().Round()
	return bounds.Min.Add (image.Pt (
		padding[artist.SideLeft] - element.scroll,
		padding[artist.SideTop] + (innerBounds.Dy() - textHeight) / 2))
}

func (element *TextBox) atPosition (position image.Point) int {
	offset := element.textOffset()
	textBoundsMin := element.valueDrawer.LayoutBounds().Min
	return element.valueDrawer.AtPosition (
		fixedutil.Pt(position.Sub(offset).Add(textBoundsMin)))
}

func (element *TextBox) HandleMouseUp (x, y int, button input.Button) {
	if button == input.ButtonLeft {
		element.dragging = 0
	}
}

func (element *TextBox) HandleKeyDown(key input.Key, modifiers input.Modifiers) {
	if element.onKeyDown != nil && element.onKeyDown(key, modifiers) {
		return
	}

	scrollMemory := element.scroll
	altered     := true
	textChanged := false
	switch {
	case key == input.KeyEnter:
		if element.onEnter != nil {
			element.onEnter()
		}
	
	case key == input.KeyBackspace:
		if len(element.text) < 1 { break }
		element.text, element.dot = textmanip.Backspace (
			element.text,
			element.dot,
			modifiers.Control)
		textChanged = true
			
	case key == input.KeyDelete:
		if len(element.text) < 1 { break }
		element.text, element.dot = textmanip.Delete (
			element.text,
			element.dot,
			modifiers.Control)
		textChanged = true
			
	case key == input.KeyLeft:
		if modifiers.Shift {
			element.dot = textmanip.SelectLeft (
				element.text,
				element.dot,
				modifiers.Control)
		} else {
			element.dot = textmanip.MoveLeft (
				element.text,
				element.dot,
				modifiers.Control)
		}
			
	case key == input.KeyRight:
		if modifiers.Shift {
			element.dot = textmanip.SelectRight (
				element.text,
				element.dot,
				modifiers.Control)
		} else {
			element.dot = textmanip.MoveRight (
				element.text,
				element.dot,
				modifiers.Control)
		}

	case key == 'a' && modifiers.Control:
		element.dot.Start = 0
		element.dot.End   = len(element.text)

	case key == 'x' && modifiers.Control:
		var lifted []rune
		element.text, element.dot, lifted = textmanip.Lift (
			element.text,
			element.dot)
		if lifted != nil {
			element.clipboardPut(lifted)
			textChanged = true
		}

	case key == 'c' && modifiers.Control:
		element.clipboardPut(element.dot.Slice(element.text))

	case key == 'v' && modifiers.Control:
		window := element.core.Window()
		if window == nil { break }
		window.Paste (func (d data.Data, err error) {
			if err != nil { return }
			reader, ok := d[data.MimePlain]
			if !ok { return }
			bytes, _ := io.ReadAll(reader)
			element.text, element.dot = textmanip.Type (
				element.text,
				element.dot,
				[]rune(string(bytes))...)
			element.notifyAsyncTextChange()
		})
		
	case key.Printable():
		element.text, element.dot = textmanip.Type (
			element.text,
			element.dot,
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

	if (textChanged || scrollMemory != element.scroll) {
		if parent, ok := element.core.Parent().(tomo.ScrollableParent); ok {
			parent.NotifyScrollBoundsChange(element)
		}
	}
	
	if altered {
		element.redo()
	}
}

func (element *TextBox) clipboardPut (text []rune) {
	window := element.core.Window()
	if window != nil {
		window.Copy(data.Bytes(data.MimePlain, []byte(string(text))))
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
	if element.dot.End > element.valueDrawer.Length() {
		element.dot = textmanip.EmptyDot(element.valueDrawer.Length())
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

func (element *TextBox) OnEnter (callback func ()) {
	element.onEnter = callback
}

func (element *TextBox) OnChange (callback func ()) {
	element.onChange = callback
}

// OnScrollBoundsChange sets a function to be called when the element's viewport
// bounds, content bounds, or scroll axes change.
func (element *TextBox) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
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
	padding := element.theme.Padding(tomo.PatternInput)
	return padding.Apply(element.Bounds()).Dx()
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
	if parent, ok := element.core.Parent().(tomo.ScrollableParent); ok {
		parent.NotifyScrollBoundsChange(element)
	}
}

// ScrollAxes returns the supported axes for scrolling.
func (element *TextBox) ScrollAxes () (horizontal, vertical bool) {
	return true, false
}

func (element *TextBox) runOnChange () {
	if element.onChange != nil {
		element.onChange()
	}
}

func (element *TextBox) scrollToCursor () {
	if !element.core.HasImage() { return }

	padding := element.theme.Padding(tomo.PatternInput)
	bounds  := padding.Apply(element.Bounds())
	bounds = bounds.Sub(bounds.Min)
	bounds.Max.X -= element.valueDrawer.Em().Round()
	cursorPosition := fixedutil.RoundPt (
		element.valueDrawer.PositionAt(element.dot.End))
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
func (element *TextBox) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	face := element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal)
	element.placeholderDrawer.SetFace(face)
	element.valueDrawer.SetFace(face)
	element.updateMinimumSize()
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *TextBox) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.redo()
}

func (element *TextBox) updateMinimumSize () {
	textBounds := element.placeholderDrawer.LayoutBounds()
	padding := element.theme.Padding(tomo.PatternInput)
	element.core.SetMinimumSize (
		padding.Horizontal() + textBounds.Dx(),
		padding.Vertical()   +
		element.placeholderDrawer.LineHeight().Round())
}

func (element *TextBox) notifyAsyncTextChange () {
	element.runOnChange()
	element.valueDrawer.SetText(element.text)
	element.scrollToCursor()
	if parent, ok := element.core.Parent().(tomo.ScrollableParent); ok {
		parent.NotifyScrollBoundsChange(element)
	}
	element.redo()
}

func (element *TextBox) redo () {
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *TextBox) draw () {
	bounds := element.Bounds()

	state := tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
	}
	pattern := element.theme.Pattern(tomo.PatternInput, state)
	padding := element.theme.Padding(tomo.PatternInput)
	innerCanvas := canvas.Cut(element.core, padding.Apply(bounds))
	pattern.Draw(element.core, bounds)
	offset := element.textOffset()

	if element.Focused() && !element.dot.Empty() {
		// draw selection bounds
		accent := element.theme.Color(tomo.ColorAccent,  state)
		canon := element.dot.Canon()
		foff  := fixedutil.Pt(offset)
		start := element.valueDrawer.PositionAt(canon.Start).Add(foff)
		end   := element.valueDrawer.PositionAt(canon.End).Add(foff)
		end.Y += element.valueDrawer.LineHeight()
		shapes.FillColorRectangle (
			innerCanvas,
			accent,
			image.Rectangle {
				fixedutil.RoundPt(start),
				fixedutil.RoundPt(end),
			})
	}

	if len(element.text) == 0 {
		// draw placeholder
		textBounds := element.placeholderDrawer.LayoutBounds()
		foreground := element.theme.Color (
			tomo.ColorForeground,
			tomo.State { Disabled: true })
		element.placeholderDrawer.Draw (
			innerCanvas,
			foreground,
			offset.Sub(textBounds.Min))
	} else {
		// draw input value
		textBounds := element.valueDrawer.LayoutBounds()
		foreground := element.theme.Color(tomo.ColorForeground, state)
		element.valueDrawer.Draw (
			innerCanvas,
			foreground,
			offset.Sub(textBounds.Min))
	}
	
	if element.Focused() && element.dot.Empty() {
		// draw cursor
		foreground := element.theme.Color(tomo.ColorForeground, state)
		cursorPosition := fixedutil.RoundPt (
			element.valueDrawer.PositionAt(element.dot.End))
		shapes.ColorLine (
			innerCanvas,
			foreground, 1,
			cursorPosition.Add(offset),
			image.Pt (
				cursorPosition.X,
				cursorPosition.Y + element.valueDrawer.
				LineHeight().Round()).Add(offset))
	}
}
