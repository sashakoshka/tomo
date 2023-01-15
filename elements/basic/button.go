package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Button is a clickable button.
type Button struct {
	*core.Core
	core core.CoreControl
	
	pressed  bool
	enabled  bool
	onClick func ()

	text   string
	drawer artist.TextDrawer
}

// NewButton creates a new button with the specified label text.
func NewButton (text string) (element *Button) {
	element = &Button { enabled: true }
	element.Core, element.core = core.NewCore(element)
	element.drawer.SetFace(theme.FontFaceRegular())
	element.core.SetSelectable(true)
	element.SetText(text)
	return
}

// Handle handles an event.
func (element *Button) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		element.draw()

	case tomo.EventMouseDown:
		if !element.enabled { break }
		
		mouseDownEvent := event.(tomo.EventMouseDown)
		element.Select()
		if mouseDownEvent.Button != tomo.ButtonLeft { break }
		element.pressed = true
		if element.core.HasImage() {
			element.draw()
			element.core.PushAll()
		}

	case tomo.EventKeyDown:
		keyDownEvent := event.(tomo.EventKeyDown)
		if keyDownEvent.Key == tomo.KeyEnter {
			element.pressed = true
			if element.core.HasImage() {
				element.draw()
				element.core.PushAll()
			}
		}

	case tomo.EventMouseUp:
		if !element.enabled { break }
	
		mouseUpEvent := event.(tomo.EventMouseUp)
		if mouseUpEvent.Button != tomo.ButtonLeft { break }
		element.pressed = false
		if element.core.HasImage() {
			element.draw()
			element.core.PushAll()
		}

		within := image.Point { mouseUpEvent.X, mouseUpEvent.Y }.
			In(element.Bounds())
			
		if within && element.onClick != nil {
			element.onClick()
		}

	case tomo.EventKeyUp:
		keyDownEvent := event.(tomo.EventKeyUp)
		if keyDownEvent.Key == tomo.KeyEnter && element.pressed {
			element.pressed = false
			if element.core.HasImage() {
				element.draw()
				element.core.PushAll()
			}
			if element.onClick != nil {
				element.onClick()
			}
		}

	case tomo.EventSelect:
		element.core.SetSelected(true)
		if element.core.HasImage() {
			element.draw()
			element.core.PushAll()
		}

	case tomo.EventDeselect:
		element.core.SetSelected(false)
		if element.core.HasImage() {
			element.draw()
			element.core.PushAll()
		}
	}
	return
}

// OnClick sets the function to be called when the button is clicked.
func (element *Button) OnClick (callback func ()) {
	element.onClick = callback
}

// Select requests that this button's parent container send it a selection
// event.
func (element *Button) Select () {
	element.core.Select()
}

// SetEnabled sets whether this button can be clicked or not.
func (element *Button) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.core.SetSelectable(enabled)
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

// SetText sets the button's label text.
func (element *Button) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText(text)
	textBounds := element.drawer.LayoutBounds()
	element.core.SetMinimumSize (
		theme.Padding() * 2 + textBounds.Dx(),
		theme.Padding() * 2 + textBounds.Dy())
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

func (element *Button) draw () {
	bounds := element.core.Bounds()

	artist.FillRectangle (
		element.core,
		theme.ButtonPattern (
			element.enabled,
			element.Selected(),
			element.pressed),
		bounds)
		
	innerBounds := bounds
	innerBounds.Min.X += theme.Padding()
	innerBounds.Min.Y += theme.Padding()
	innerBounds.Max.X -= theme.Padding()
	innerBounds.Max.Y -= theme.Padding()

	textBounds := element.drawer.LayoutBounds()
	offset := image.Point {
		X: theme.Padding() + (innerBounds.Dx() - textBounds.Dx()) / 2,
		Y: theme.Padding() + (innerBounds.Dy() - textBounds.Dy()) / 2,
	}

	// account for the fact that the bounding rectangle will be shifted over
	// due to the bounds origin being at the baseline of the first line
	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	if element.pressed {
		offset = offset.Add(theme.SinkOffsetVector())
	}

	foreground := theme.ForegroundPattern(element.enabled)
	element.drawer.Draw(element.core, foreground, offset)
}
