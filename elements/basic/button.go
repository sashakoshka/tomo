package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"

type Button struct {
	*Core
	core CoreControl
	
	pressed  bool
	enabled  bool
	selected bool
	onClick func ()

	text   string
	drawer artist.TextDrawer
}

func NewButton (text string) (element *Button) {
	element = &Button { enabled: true }
	element.Core, element.core = NewCore(element)
	element.drawer.SetFace(theme.FontFaceRegular())
	element.SetText(text)
	return
}

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

	case tomo.EventSelect:
		element.selected = true

	case tomo.EventDeselect:
		element.selected = false
	// TODO: handle selection events, and the enter key
	}
	return
}

func (element *Button) OnClick (callback func ()) {
	element.onClick = callback
}

func (element *Button) AdvanceSelection (direction int) (ok bool) {
	wasSelected := element.selected
	element.selected = false
	if element.core.HasImage() && wasSelected {
		element.draw()
		element.core.PushAll()
	}
	return
}

func (element *Button) Selectable () (selectable bool) {
	return true
}

func (element *Button) Select () {
	element.core.Select()
}

func (element *Button) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	if element.core.HasImage () {
		element.draw()
		element.core.PushAll()
	}
}

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

	artist.ChiseledRectangle (
		element.core,
		theme.RaisedProfile (
			element.pressed,
			element.enabled,
			element.selected),
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

	foreground := theme.ForegroundImage()
	if !element.enabled {
		foreground = theme.DisabledForegroundImage()
	}

	element.drawer.Draw(element.core, foreground, offset)
}
