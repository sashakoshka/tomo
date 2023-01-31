package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

var switchCase = theme.C("basic", "switch")

// Switch is a toggle-able on/off switch with an optional label. It is
// functionally identical to Checkbox, but plays a different semantic role.
type Switch struct {
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

// NewSwitch creates a new switch with the specified label text.
func NewSwitch (text string, on bool) (element *Switch) {
	element = &Switch { checked: on, text: text }
	element.Core, element.core = core.NewCore(element.draw)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore (func () {
		if element.core.HasImage () {
			element.draw()
			element.core.DamageAll()
		}
	})
	element.drawer.SetFace(theme.FontFaceRegular())
	element.drawer.SetText([]rune(text))
	element.calculateMinimumSize()
	return
}

func (element *Switch) HandleMouseDown (x, y int, button tomo.Button) {
	if !element.Enabled() { return }
	element.Focus()
	element.pressed = true
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Switch) HandleMouseUp (x, y int, button tomo.Button) {
	if button != tomo.ButtonLeft || !element.pressed { return }

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

func (element *Switch) HandleKeyDown (key tomo.Key, modifiers tomo.Modifiers) {
	if key == tomo.KeyEnter {
		element.pressed = true
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
	}
}

func (element *Switch) HandleKeyUp (key tomo.Key, modifiers tomo.Modifiers) {
	if key == tomo.KeyEnter && element.pressed {
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
	element.calculateMinimumSize()
	
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Switch) calculateMinimumSize () {
	textBounds := element.drawer.LayoutBounds()
	lineHeight := element.drawer.LineHeight().Round()
	
	if element.text == "" {
		element.core.SetMinimumSize(lineHeight * 2, lineHeight)
	} else {
		element.core.SetMinimumSize (
			lineHeight * 2 + theme.Padding() + textBounds.Dx(),
			lineHeight)
	}
}

func (element *Switch) draw () {
	bounds := element.core.Bounds()
	handleBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dy())
	gutterBounds := image.Rect(0, 0, bounds.Dy() * 2, bounds.Dy())
	backgroundPattern, _ := theme.BackgroundPattern(theme.PatternState {
		Case: switchCase,
	})
	artist.FillRectangle ( element.core, backgroundPattern, bounds)

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

	gutterPattern, _ := theme.GutterPattern(theme.PatternState {
		Case: switchCase,
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
		Pressed:  element.pressed,
	})
	artist.FillRectangle(element.core, gutterPattern, gutterBounds)
	
	handlePattern, _ := theme.HandlePattern(theme.PatternState {
		Case: switchCase,
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
		Pressed:  element.pressed,
	})
	artist.FillRectangle(element.core, handlePattern, handleBounds)

	textBounds := element.drawer.LayoutBounds()
	offset := image.Point {
		X: bounds.Dy() * 2 + theme.Padding(),
	}

	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground, _ := theme.ForegroundPattern (theme.PatternState {
		Case: switchCase,
		Disabled: !element.Enabled(),
	})
	element.drawer.Draw(element.core, foreground, offset)
}
