package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// Switch is a toggle-able on/off switch with an optional label. It is
// functionally identical to Checkbox, but plays a different semantic role.
type Switch struct {
	entity tomo.FocusableEntity
	drawer textdraw.Drawer

	enabled bool
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
		text:    text,
		enabled: true,
	}
	element.entity = tomo.NewEntity(element).(tomo.FocusableEntity)
	element.theme.Case = tomo.C("tomo", "switch")
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.drawer.SetText([]rune(text))
	element.updateMinimumSize()
	return
}

// Entity returns this element's entity.
func (element *Switch) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Switch) Draw (destination canvas.Canvas) {
	bounds := element.entity.Bounds()
	handleBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dy()).Add(bounds.Min)
	gutterBounds := image.Rect(0, 0, bounds.Dy() * 2, bounds.Dy()).Add(bounds.Min)

	state := tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
		Pressed:  element.pressed,
		On:       element.checked,
	}

	element.entity.DrawBackground(destination)

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
		tomo.PatternGutter, state)
	gutterPattern.Draw(destination, gutterBounds)
	
	handlePattern := element.theme.Pattern (
		tomo.PatternHandle, state)
	handlePattern.Draw(destination, handleBounds)

	textBounds := element.drawer.LayoutBounds()
	offset := bounds.Min.Add(image.Point {
		X: bounds.Dy() * 2 +
			element.theme.Margin(tomo.PatternBackground).X,
	})

	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground := element.theme.Color(tomo.ColorForeground, state)
	element.drawer.Draw(destination, foreground, offset)
}

func (element *Switch) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *Switch) HandleMouseDown  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if !element.Enabled() { return }
	element.Focus()
	element.pressed = true
	element.entity.Invalidate()
}

func (element *Switch) HandleMouseUp (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if button != input.ButtonLeft || !element.pressed { return }

	element.pressed = false
	within := position.In(element.entity.Bounds())
	if within {
		element.checked = !element.checked
	}
	
	element.entity.Invalidate()
	if within && element.onToggle != nil {
		element.onToggle()
	}
}

func (element *Switch) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter {
		element.pressed = true
		element.entity.Invalidate()
	}
}

func (element *Switch) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.checked = !element.checked
		element.entity.Invalidate()
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

// Focus gives this element input focus.
func (element *Switch) Focus () {
	if !element.entity.Focused() { element.entity.Focus() }
}

// Enabled returns whether this switch is enabled or not.
func (element *Switch) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this switch can be toggled or not.
func (element *Switch) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}

// SetText sets the checkbox's label text.
func (element *Switch) SetText (text string) {
	if element.text == text { return }

	element.text = text
	element.drawer.SetText([]rune(text))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetTheme sets the element's theme.
func (element *Switch) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *Switch) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Switch) updateMinimumSize () {
	textBounds := element.drawer.LayoutBounds()
	lineHeight := element.drawer.LineHeight().Round()
	
	if element.text == "" {
		element.entity.SetMinimumSize(lineHeight * 2, lineHeight)
	} else {
		element.entity.SetMinimumSize (
			lineHeight * 2 +
			element.theme.Margin(tomo.PatternBackground).X +
			textBounds.Dx(),
			lineHeight)
	}
}
