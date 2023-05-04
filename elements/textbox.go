package elements

import "io"
import "time"
import "image"
import "tomo"
import "tomo/data"
import "tomo/input"
import "art"
import "tomo/textdraw"
import "tomo/textmanip"
import "tomo/fixedutil"
import "art/shapes"

var textBoxCase = tomo.C("tomo", "textBox")

// TextBox is a single-line text input.
type TextBox struct {
	entity tomo.Entity
	
	enabled     bool
	lastClick   time.Time
	dragging    int
	dot         textmanip.Dot
	scroll      int
	placeholder string
	text        []rune
	
	placeholderDrawer textdraw.Drawer
	valueDrawer       textdraw.Drawer
	
	onKeyDown func (key input.Key, modifiers input.Modifiers) (handled bool)
	onChange  func ()
	onEnter   func ()
	onScrollBoundsChange func ()
}

// NewTextBox creates a new text box with the specified placeholder text, and
// a value. When the value is empty, the placeholder will be displayed in gray
// text.
func NewTextBox (placeholder, value string) (element *TextBox) {
	element = &TextBox { enabled: true }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.placeholder = placeholder
	element.placeholderDrawer.SetFace (element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal, textBoxCase))
	element.valueDrawer.SetFace (element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal, textBoxCase))
	element.placeholderDrawer.SetText([]rune(placeholder))
	element.updateMinimumSize()
	element.SetValue(value)
	return
}

// Entity returns this element's entity.
func (element *TextBox) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *TextBox) Draw (destination art.Canvas) {
	bounds := element.entity.Bounds()

	state := element.state()
	pattern := element.entity.Theme().Pattern(tomo.PatternInput, state, textBoxCase)
	padding := element.entity.Theme().Padding(tomo.PatternInput, textBoxCase)
	innerCanvas := art.Cut(destination, padding.Apply(bounds))
	pattern.Draw(destination, bounds)
	offset := element.textOffset()

	if element.entity.Focused() && !element.dot.Empty() {
		// draw selection bounds
		accent := element.entity.Theme().Color(tomo.ColorAccent, state, textBoxCase)
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
		foreground := element.entity.Theme().Color (
			tomo.ColorForeground,
			tomo.State { Disabled: true }, textBoxCase)
		element.placeholderDrawer.Draw (
			innerCanvas,
			foreground,
			offset.Sub(textBounds.Min))
	} else {
		// draw input value
		textBounds := element.valueDrawer.LayoutBounds()
		foreground := element.entity.Theme().Color(tomo.ColorForeground, state, textBoxCase)
		element.valueDrawer.Draw (
			innerCanvas,
			foreground,
			offset.Sub(textBounds.Min))
	}
	
	if element.entity.Focused() && element.dot.Empty() {
		// draw cursor
		foreground := element.entity.Theme().Color(tomo.ColorForeground, state, textBoxCase)
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

// Layout causes the element to perform a layout operation.
func (element *TextBox) Layout () {
	element.scrollToCursor()
}

func (element *TextBox) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *TextBox) HandleMouseDown  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if !element.Enabled() { return }
	element.Focus()

	switch button {
	case input.ButtonLeft:
		runeIndex := element.atPosition(position)
		if runeIndex == -1 { return }
		
		if time.Since(element.lastClick) < element.entity.Config().DoubleClickDelay() {
			element.dragging = 2
			element.dot = textmanip.WordAround(element.text, runeIndex)
		} else {
			element.dragging = 1
			element.dot = textmanip.EmptyDot(runeIndex)
			element.lastClick = time.Now()
		}
		
		element.entity.Invalidate()
	case input.ButtonRight:
		element.contextMenu(position)
	}
}

func (element *TextBox) HandleMouseUp  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if button == input.ButtonLeft {
		element.dragging = 0
	}
}

func (element *TextBox) HandleMotion (position image.Point) {
	if !element.Enabled() { return }

	switch element.dragging {
	case 1:
		runeIndex := element.atPosition(position)
		if runeIndex > -1 {
			element.dot.End = runeIndex
			element.entity.Invalidate()
		}
		
	case 2:
		runeIndex := element.atPosition(position)
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
			element.entity.Invalidate()
		}
	}
}

func (element *TextBox) textOffset () image.Point {
	padding     := element.entity.Theme().Padding(tomo.PatternInput, textBoxCase)
	bounds      := element.entity.Bounds()
	innerBounds := padding.Apply(bounds)
	textHeight  := element.valueDrawer.LineHeight().Round()
	return bounds.Min.Add (image.Pt (
		padding[art.SideLeft] - element.scroll,
		padding[art.SideTop] + (innerBounds.Dy() - textHeight) / 2))
}

func (element *TextBox) atPosition (position image.Point) int {
	offset := element.textOffset()
	textBoundsMin := element.valueDrawer.LayoutBounds().Min
	return element.valueDrawer.AtPosition (
		fixedutil.Pt(position.Sub(offset).Add(textBoundsMin)))
}

func (element *TextBox) HandleKeyDown(key input.Key, modifiers input.Modifiers) {
	if element.onKeyDown != nil && element.onKeyDown(key, modifiers) {
		return
	}

	scrollMemory := element.scroll
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
		element.scrollToCursor()
		element.entity.Invalidate()
			
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
		element.scrollToCursor()
		element.entity.Invalidate()

	case key == 'a' && modifiers.Control:
		element.dot.Start = 0
		element.dot.End   = len(element.text)
		element.scrollToCursor()
		element.entity.Invalidate()

	case key == 'x' && modifiers.Control: element.Cut()
	case key == 'c' && modifiers.Control: element.Copy()
	case key == 'v' && modifiers.Control: element.Paste()

	case key == input.KeyMenu:
		pos := fixedutil.RoundPt(element.valueDrawer.PositionAt(element.dot.End)).
			Add(element.textOffset())
		pos.Y += element.valueDrawer.LineHeight().Round()
		element.contextMenu(pos)
		
	case key.Printable():
		element.text, element.dot = textmanip.Type (
			element.text,
			element.dot,
			rune(key))
		textChanged = true
	}

	if textChanged {
		element.runOnChange()
		element.valueDrawer.SetText(element.text)
		element.scrollToCursor()
		element.entity.Invalidate()
	}

	if (textChanged || scrollMemory != element.scroll) {
		element.entity.NotifyScrollBoundsChange()
	}
}

// Cut cuts the selected text in the text box and places it in the clipboard.
func (element *TextBox) Cut () {
	var lifted []rune
	element.text, element.dot, lifted = textmanip.Lift (
		element.text,
		element.dot)
	if lifted != nil {
		element.clipboardPut(lifted)
		element.notifyAsyncTextChange()
	}
}

// Copy copies the selected text in the text box and places it in the clipboard.
func (element *TextBox) Copy () {
	element.clipboardPut(element.dot.Slice(element.text))
}

// Paste pastes text data from the clipboard into the text box.
func (element *TextBox) Paste () {
	window := element.entity.Window()
	if window == nil { return }
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
}

func (element *TextBox) HandleKeyUp(key input.Key, modifiers input.Modifiers) { }

// SetPlaceholder sets the element's placeholder text.
func (element *TextBox) SetPlaceholder (placeholder string) {
	if element.placeholder == placeholder { return }
	
	element.placeholder = placeholder
	element.placeholderDrawer.SetText([]rune(placeholder))
	
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetValue sets the input's value.
func (element *TextBox) SetValue (text string) {
	// if element.text == text { return }

	element.text = []rune(text)
	element.runOnChange()
	element.valueDrawer.SetText(element.text)
	if element.dot.End > element.valueDrawer.Length() {
		element.dot = textmanip.EmptyDot(element.valueDrawer.Length())
	}
	element.scrollToCursor()
	element.entity.Invalidate()
}

// Value returns the input's value.
func (element *TextBox) Value () (value string) {
	return string(element.text)
}

// Filled returns whether or not this element has a value.
func (element *TextBox) Filled () (filled bool) {
	return len(element.text) > 0
}

// OnKeyDown specifies a function to be called when a key is pressed within the
// text input.
func (element *TextBox) OnKeyDown (
	callback func (key input.Key, modifiers input.Modifiers) (handled bool),
) {
	element.onKeyDown = callback
}

// OnEnter specifies a function to be called when the enter key is pressed
// within this input.
func (element *TextBox) OnEnter (callback func ()) {
	element.onEnter = callback
}

// OnChange specifies a function to be called when the value of this input
// changes.
func (element *TextBox) OnChange (callback func ()) {
	element.onChange = callback
}

// OnScrollBoundsChange sets a function to be called when the element's viewport
// bounds, content bounds, or scroll axes change.
func (element *TextBox) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
}

// Focus gives this element input focus.
func (element *TextBox) Focus () {
	if !element.entity.Focused() { element.entity.Focus() }
}

// Enabled returns whether this label can be edited or not.
func (element *TextBox) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this label can be edited or not.
func (element *TextBox) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
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

	element.entity.Invalidate()
	element.entity.NotifyScrollBoundsChange()
}

// ScrollAxes returns the supported axes for scrolling.
func (element *TextBox) ScrollAxes () (horizontal, vertical bool) {
	return true, false
}

func (element *TextBox) HandleThemeChange () {
	face := element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal,
		textBoxCase)
	element.placeholderDrawer.SetFace(face)
	element.valueDrawer.SetFace(face)
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *TextBox) contextMenu (position image.Point) {
	window := element.entity.Window()
	menu, err := window.NewMenu(image.Rectangle { position, position })
	if err != nil { return }

	closeAnd := func (callback func ()) func () {
		return func () { callback(); menu.Close() }
	}

	cutButton := NewButton("Cut")
	cutButton.ShowText(false)
	cutButton.SetIcon(tomo.IconCut)
	cutButton.SetEnabled(!element.dot.Empty())
	cutButton.OnClick(closeAnd(element.Cut))
	
	copyButton := NewButton("Copy")
	copyButton.ShowText(false)
	copyButton.SetIcon(tomo.IconCopy)
	copyButton.SetEnabled(!element.dot.Empty())
	copyButton.OnClick(closeAnd(element.Copy))
	
	pasteButton := NewButton("Paste")
	pasteButton.ShowText(false)
	pasteButton.SetIcon(tomo.IconPaste)
	pasteButton.OnClick(closeAnd(element.Paste))

	menu.Adopt (NewHBox (
		SpaceNone,
		pasteButton,
		copyButton,
		cutButton,
	))
	pasteButton.Focus()
	menu.Show()
}

func (element *TextBox) runOnChange () {
	if element.onChange != nil {
		element.onChange()
	}
}

func (element *TextBox) scrollViewportWidth () (width int) {
	padding := element.entity.Theme().Padding(tomo.PatternInput, textBoxCase)
	return padding.Apply(element.entity.Bounds()).Dx()
}

func (element *TextBox) scrollToCursor () {
	padding := element.entity.Theme().Padding(tomo.PatternInput, textBoxCase)
	bounds  := padding.Apply(element.entity.Bounds())
	bounds = bounds.Sub(bounds.Min)
	bounds.Max.X -= element.valueDrawer.Em().Round()
	cursorPosition := fixedutil.RoundPt (
		element.valueDrawer.PositionAt(element.dot.End))
	cursorPosition.X -= element.scroll
	maxX := bounds.Max.X
	minX := maxX
	if cursorPosition.X > maxX {
		element.scroll += cursorPosition.X - maxX
		element.entity.NotifyScrollBoundsChange()
		element.entity.Invalidate()
	} else if cursorPosition.X < minX {
		element.scroll -= minX - cursorPosition.X
		if element.scroll < 0 { element.scroll = 0 }
		element.entity.NotifyScrollBoundsChange()
		element.entity.Invalidate()
	}
}

func (element *TextBox) updateMinimumSize () {
	textBounds := element.placeholderDrawer.LayoutBounds()
	padding := element.entity.Theme().Padding(tomo.PatternInput, textBoxCase)
	element.entity.SetMinimumSize (
		padding.Horizontal() + textBounds.Dx(),
		padding.Vertical()   +
		element.placeholderDrawer.LineHeight().Round())
}

func (element *TextBox) notifyAsyncTextChange () {
	element.runOnChange()
	element.valueDrawer.SetText(element.text)
	element.scrollToCursor()
	element.entity.Invalidate()
}

func (element *TextBox) clipboardPut (text []rune) {
	window := element.entity.Window()
	if window != nil {
		window.Copy(data.Bytes(data.MimePlain, []byte(string(text))))
	}
}

func (element *TextBox) state () tomo.State {
	return tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
	}
}
