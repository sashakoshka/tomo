package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

var buttonCase = theme.C("basic", "button")

// Button is a clickable button.
type Button struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl
	drawer artist.TextDrawer

	pressed bool
	text    string
	
	onClick func ()
}

// NewButton creates a new button with the specified label text.
func NewButton (text string) (element *Button) {
	element = &Button { }
	element.Core, element.core = core.NewCore(element.draw)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore (func () {
		if element.core.HasImage () {
			element.draw()
			element.core.DamageAll()
		}
	})
	element.drawer.SetFace(theme.FontFaceRegular())
	element.SetText(text)
	return
}

func (element *Button) HandleMouseDown (x, y int, button tomo.Button) {
	if !element.Enabled()  { return }
	if !element.Focused() { element.Focus() }
	if button != tomo.ButtonLeft { return }
	element.pressed = true
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Button) HandleMouseUp (x, y int, button tomo.Button) {
	if button != tomo.ButtonLeft { return }
	element.pressed = false
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}

	within := image.Point { x, y }.
		In(element.Bounds())
		
	if !element.Enabled() { return }
	if within && element.onClick != nil {
		element.onClick()
	}
}

func (element *Button) HandleMouseMove (x, y int) { }
func (element *Button) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Button) HandleKeyDown (key tomo.Key, modifiers tomo.Modifiers) {
	if !element.Enabled() { return }
	if key == tomo.KeyEnter {
		element.pressed = true
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
	}
}

func (element *Button) HandleKeyUp(key tomo.Key, modifiers tomo.Modifiers) {
	if key == tomo.KeyEnter && element.pressed {
		element.pressed = false
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
		if !element.Enabled() { return }
		if element.onClick != nil {
			element.onClick()
		}
	}
}

// OnClick sets the function to be called when the button is clicked.
func (element *Button) OnClick (callback func ()) {
	element.onClick = callback
}

// SetEnabled sets whether this button can be clicked or not.
func (element *Button) SetEnabled (enabled bool) {
	element.focusableControl.SetEnabled(enabled)
}

// SetText sets the button's label text.
func (element *Button) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
	textBounds := element.drawer.LayoutBounds()
	_, inset := theme.ButtonPattern(theme.PatternState { Case: buttonCase })
	minimumSize := inset.Inverse().Apply(textBounds).Inset(-theme.Padding())
	element.core.SetMinimumSize(minimumSize.Dx(), minimumSize.Dy())
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Button) draw () {
	bounds := element.Bounds()

	pattern, inset := theme.ButtonPattern(theme.PatternState {
		Case: buttonCase,
		Disabled: !element.Enabled(),
		Focused: element.Focused(),
		Pressed:  element.pressed,
	})

	artist.FillRectangle(element, pattern, bounds)
		
	innerBounds := inset.Apply(bounds)

	textBounds := element.drawer.LayoutBounds()
	offset := image.Point {
		X: innerBounds.Min.X + (innerBounds.Dx() - textBounds.Dx()) / 2,
		Y: innerBounds.Min.X + (innerBounds.Dy() - textBounds.Dy()) / 2,
	}

	// account for the fact that the bounding rectangle will be shifted over
	// due to the bounds origin being at the baseline of the first line
	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground, _ := theme.ForegroundPattern (theme.PatternState {
		Case: buttonCase,
		Disabled: !element.Enabled(),
	})
	element.drawer.Draw(element, foreground, offset)
}
