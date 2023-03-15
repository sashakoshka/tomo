package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Switch is a toggle-able on/off switch with an optional label. It is
// functionally identical to Checkbox, but plays a different semantic role.
type Switch struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl
	drawer textdraw.Drawer

	pressed bool
	checked bool
	text    string
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onToggle func ()
}

// NewSwitch creates a new switch with the specified label text.
func NewSwitch (text string, on bool) (element *Switch) {
	element = &Switch {
		checked: on,
		text: text,
	}
	element.theme.Case = theme.C("basic", "switch")
	element.Core, element.core = core.NewCore(element, element.draw)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore(element.core, element.redo)
	element.drawer.SetText([]rune(text))
	element.updateMinimumSize()
	return
}

func (element *Switch) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	element.Focus()
	element.pressed = true
	element.redo()
}

func (element *Switch) HandleMouseUp (x, y int, button input.Button) {
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

func (element *Switch) HandleMouseMove (x, y int) { }
func (element *Switch) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Switch) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter {
		element.pressed = true
		element.redo()
	}
}

func (element *Switch) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.checked = !element.checked
		element.redo()
		if element.onToggle != nil {
			element.onToggle()
		}
	}
}

// OnToggle sets the function to be called when the switch is flipped.
func (element *Switch) OnToggle (callback func ()) {
	element.onToggle = callback
}

// Value reports whether or not the switch is currently on.
func (element *Switch) Value () (on bool) {
	return element.checked
}

// SetEnabled sets whether this switch can be flipped or not.
func (element *Switch) SetEnabled (enabled bool) {
	element.focusableControl.SetEnabled(enabled)
}

// SetText sets the checkbox's label text.
func (element *Switch) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
	element.updateMinimumSize()
	element.redo()
}

// SetTheme sets the element's theme.
func (element *Switch) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawer.SetFace (element.theme.FontFace (
		theme.FontStyleRegular,
		theme.FontSizeNormal))
	element.updateMinimumSize()
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *Switch) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.redo()
}

func (element *Switch) redo () {
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Switch) updateMinimumSize () {
	textBounds := element.drawer.LayoutBounds()
	lineHeight := element.drawer.LineHeight().Round()
	
	if element.text == "" {
		element.core.SetMinimumSize(lineHeight * 2, lineHeight)
	} else {
		element.core.SetMinimumSize (
			lineHeight * 2 +
			element.theme.Margin(theme.PatternBackground).X +
			textBounds.Dx(),
			lineHeight)
	}
}

func (element *Switch) draw () {
	bounds := element.Bounds()
	handleBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dy()).Add(bounds.Min)
	gutterBounds := image.Rect(0, 0, bounds.Dy() * 2, bounds.Dy()).Add(bounds.Min)

	state := theme.State {
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
		Pressed:  element.pressed,
	}
	backgroundPattern := element.theme.Pattern (
		theme.PatternBackground, state)
	backgroundPattern.Draw(element.core, bounds)

	if element.checked {
		handleBounds.Min.X += bounds.Dy()
		handleBounds.Max.X += bounds.Dy()
		if element.pressed {
			handleBounds.Min.X -= 2
			handleBounds.Max.X -= 2
		}
	} else {
		if element.pressed {
			handleBounds.Min.X += 2
			handleBounds.Max.X += 2
		}
	}

	gutterPattern := element.theme.Pattern (
		theme.PatternGutter, state)
	gutterPattern.Draw(element.core, gutterBounds)
	
	handlePattern := element.theme.Pattern (
		theme.PatternHandle, state)
	handlePattern.Draw(element.core, handleBounds)

	textBounds := element.drawer.LayoutBounds()
	offset := bounds.Min.Add(image.Point {
		X: bounds.Dy() * 2 +
			element.theme.Margin(theme.PatternBackground).X,
	})

	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground := element.theme.Color (
		theme.ColorForeground, state)
	element.drawer.Draw(element.core, foreground, offset)
}
