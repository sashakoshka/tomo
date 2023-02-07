package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Checkbox is a toggle-able checkbox with a label.
type Checkbox struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl
	drawer artist.TextDrawer

	pressed bool
	checked bool
	text    string
	
	onToggle func ()
}

// NewCheckbox creates a new cbeckbox with the specified label text.
func NewCheckbox (text string, checked bool) (element *Checkbox) {
	element = &Checkbox { checked: checked }
	element.Core, element.core = core.NewCore (
		element.draw,
		element.redo,
		element.redo,
		theme.C("basic", "checkbox"))
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore(element.redo)
	element.SetText(text)
	return
}

func (element *Checkbox) redo () {
	element.drawer.SetFace (
		element.core.FontFace(theme.FontStyleRegular,
		theme.FontSizeNormal))
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Checkbox) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	element.Focus()
	element.pressed = true
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Checkbox) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft || !element.pressed { return }

	element.pressed = false
	within := image.Point { x, y }.
		In(element.Bounds())
	if within {
		element.checked = !element.checked
	}
	
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
	if within && element.onToggle != nil {
		element.onToggle()
	}
}

func (element *Checkbox) HandleMouseMove (x, y int) { }
func (element *Checkbox) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Checkbox) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter {
		element.pressed = true
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
	}
}

func (element *Checkbox) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.checked = !element.checked
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
		if element.onToggle != nil {
			element.onToggle()
		}
	}
}

// OnToggle sets the function to be called when the checkbox is toggled.
func (element *Checkbox) OnToggle (callback func ()) {
	element.onToggle = callback
}

// Value reports whether or not the checkbox is currently checked.
func (element *Checkbox) Value () (checked bool) {
	return element.checked
}

// SetEnabled sets whether this checkbox can be toggled or not.
func (element *Checkbox) SetEnabled (enabled bool) {
	element.focusableControl.SetEnabled(enabled)
}

// SetText sets the checkbox's label text.
func (element *Checkbox) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
	textBounds := element.drawer.LayoutBounds()
	
	if text == "" {
		element.core.SetMinimumSize(textBounds.Dy(), textBounds.Dy())
	} else {
		element.core.SetMinimumSize (
			textBounds.Dy() + element.core.Config().Padding() + textBounds.Dx(),
			textBounds.Dy())
	}
	
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Checkbox) draw () {
	bounds := element.Bounds()
	boxBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dy()).Add(bounds.Min)

	state := theme.PatternState {
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
		Pressed:  element.pressed,
		On:       element.checked,
	}

	backgroundPattern := element.core.Pattern(theme.PatternBackground, state)
	artist.FillRectangle(element, backgroundPattern, bounds)

	pattern := element.core.Pattern (theme.PatternButton, state)
	artist.FillRectangle(element, pattern, boxBounds)

	textBounds := element.drawer.LayoutBounds()
	offset := bounds.Min.Add(image.Point {
		X: bounds.Dy() + element.core.Config().Padding(),
	})

	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground := element.core.Pattern(theme.PatternForeground, state)
	element.drawer.Draw(element, foreground, offset)
}
