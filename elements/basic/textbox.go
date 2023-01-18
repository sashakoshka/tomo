package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/textmanip"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type TextBox struct {
	*core.Core
	core core.CoreControl
	
	enabled  bool
	selected bool

	cursor int
	placeholder string
	text        []rune
	placeholderDrawer artist.TextDrawer
	valueDrawer       artist.TextDrawer
}

func NewTextBox (placeholder, text string) (element *TextBox) {
	element = &TextBox { enabled: true }
	element.Core, element.core = core.NewCore(element)
	element.placeholderDrawer.SetFace(theme.FontFaceRegular())
	element.valueDrawer.SetFace(theme.FontFaceRegular())
	element.placeholder = placeholder
	element.placeholderDrawer.SetText([]rune(placeholder))
	element.updateMinimumSize()
	element.SetText(text)
	return
}

func (element *TextBox) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.draw()
}

func (element *TextBox) HandleMouseDown (x, y int, button tomo.Button) {
	element.Select()
}

func (element *TextBox) HandleMouseUp (x, y int, button tomo.Button) { }
func (element *TextBox) HandleMouseMove (x, y int) { }
func (element *TextBox) HandleScroll (x, y int, deltaX, deltaY float64) { }

func (element *TextBox) HandleKeyDown (
	key tomo.Key,
	modifiers tomo.Modifiers,
	repeated bool,
) {
	altered := true
	switch {
	case key == tomo.KeyBackspace:
		if len(element.text) < 1 { break }
		element.text, element.cursor = textmanip.Backspace (
			element.text,
			element.cursor,
			modifiers.Control)
			
	case key == tomo.KeyDelete:
		if len(element.text) < 1 { break }
		element.text, element.cursor = textmanip.Delete (
			element.text,
			element.cursor,
			modifiers.Control)
			
	case key == tomo.KeyLeft:
		element.cursor = textmanip.MoveLeft (
			element.text,
			element.cursor,
			modifiers.Control)
			
	case key == tomo.KeyRight:
		element.cursor = textmanip.MoveRight (
			element.text,
			element.cursor,
			modifiers.Control)
		
	case key.Printable():
		element.text, element.cursor = textmanip.Type (
			element.text,
			element.cursor,
			rune(key))
			
	default:
		altered = false
	}

	if altered {
		element.valueDrawer.SetText(element.text)
	}
	
	if altered && element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

func (element *TextBox) HandleKeyUp(key tomo.Key, modifiers tomo.Modifiers) { }

func (element *TextBox) Selected () (selected bool) {
	return element.selected
}

func (element *TextBox) Select () {
	element.core.RequestSelection()
}

func (element *TextBox) HandleSelection (
	direction tomo.SelectionDirection,
) (
	accepted bool,
) {
	direction = direction.Canon()
	if !element.enabled { return false }
	if element.selected && direction != tomo.SelectionDirectionNeutral {
		return false
	}
	
	element.selected = true
	if element.core.HasImage() {
		element.draw()
		element.core.PushAll()
	}
	return true
}

func (element *TextBox) HandleDeselection () {
	element.selected = false
	if element.core.HasImage() {
		element.draw()
		element.core.PushAll()
	}
}

func (element *TextBox) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

func (element *TextBox) SetPlaceholder (placeholder string) {
	if element.placeholder == placeholder { return }
	
	element.placeholder = placeholder
	element.placeholderDrawer.SetText([]rune(placeholder))
	
	element.updateMinimumSize()
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

func (element *TextBox) updateMinimumSize () {
	textBounds := element.placeholderDrawer.LayoutBounds()
	element.core.SetMinimumSize (
		textBounds.Dx() +
		theme.Padding() * 2,
		element.placeholderDrawer.LineHeight().Round() +
		theme.Padding() * 2)
}

func (element *TextBox) SetText (text string) {
	// if element.text == text { return }

	element.text = []rune(text)
	element.valueDrawer.SetText(element.text)
	if element.cursor > element.valueDrawer.Length() {
		element.cursor = element.valueDrawer.Length()
	}
	
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

func (element *TextBox) draw () {
	bounds := element.core.Bounds()

	artist.FillRectangle (
		element.core,
		theme.InputPattern (
			element.enabled,
			element.Selected()),
		bounds)
		
	innerBounds := bounds
	innerBounds.Min.X += theme.Padding()
	innerBounds.Min.Y += theme.Padding()
	innerBounds.Max.X -= theme.Padding()
	innerBounds.Max.Y -= theme.Padding()

	if len(element.text) == 0 && !element.selected {
		// draw placeholder
		textBounds := element.placeholderDrawer.LayoutBounds()
		offset := image.Point {
			X: theme.Padding(),
			Y: theme.Padding(),
		}
		foreground := theme.ForegroundPattern(false)
		element.placeholderDrawer.Draw (
			element.core,
			foreground,
			offset.Sub(textBounds.Min))
	} else {
		// draw input value
		textBounds := element.valueDrawer.LayoutBounds()
		offset := image.Point {
			X: theme.Padding(),
			Y: theme.Padding(),
		}
		foreground := theme.ForegroundPattern(element.enabled)
		element.valueDrawer.Draw (
			element.core,
			foreground,
			offset.Sub(textBounds.Min))

		if element.selected {
			// cursor
			cursorPosition := element.valueDrawer.PositionOf (
				element.cursor)
			artist.Line (
				element.core,
				theme.ForegroundPattern(true), 1,
				cursorPosition.Add(offset),
				image.Pt (
					cursorPosition.X,
					cursorPosition.Y + element.valueDrawer.
					LineHeight().Round()).Add(offset))
		}
	}
}
