package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type Button struct {
	*core.Core
	core core.CoreControl
	
	pressed  bool
	enabled  bool
	onClick func ()

	text   string
	drawer artist.TextDrawer
}

func NewButton (text string) (element *Button) {
	element = &Button { enabled: true }
	element.Core, element.core = core.NewCore(element)
	element.drawer.SetFace(theme.FontFaceRegular())
	element.core.SetSelectable(true)
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

func (element *Button) OnClick (callback func ()) {
	element.onClick = callback
}

func (element *Button) Select () {
	element.core.Select()
}

func (element *Button) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.core.SetSelectable(enabled)
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
			element.Selected()),
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
