package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

var checkboxCase = theme.C("basic", "checkbox")

// Checkbox is a toggle-able checkbox with a label.
type Checkbox struct {
	*core.Core
	*core.SelectableCore
	core core.CoreControl
	selectableControl core.SelectableCoreControl
	drawer artist.TextDrawer

	pressed bool
	checked bool
	text    string
	
	onToggle func ()
}

// NewCheckbox creates a new cbeckbox with the specified label text.
func NewCheckbox (text string, checked bool) (element *Checkbox) {
	element = &Checkbox { checked: checked }
	element.Core, element.core = core.NewCore(element)
	element.SelectableCore,
	element.selectableControl = core.NewSelectableCore (func () {
		if element.core.HasImage () {
			element.draw()
			element.core.DamageAll()
		}
	})
	element.drawer.SetFace(theme.FontFaceRegular())
	element.SetText(text)
	return
}

// Resize changes this element's size.
func (element *Checkbox) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.draw()
}

func (element *Checkbox) HandleMouseDown (x, y int, button tomo.Button) {
	if !element.Enabled() { return }
	element.Select()
	element.pressed = true
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Checkbox) HandleMouseUp (x, y int, button tomo.Button) {
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

func (element *Checkbox) HandleMouseMove (x, y int) { }
func (element *Checkbox) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Checkbox) HandleKeyDown (key tomo.Key, modifiers tomo.Modifiers) {
	if key == tomo.KeyEnter {
		element.pressed = true
		if element.core.HasImage() {
			element.draw()
			element.core.DamageAll()
		}
	}
}

func (element *Checkbox) HandleKeyUp (key tomo.Key, modifiers tomo.Modifiers) {
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
	element.selectableControl.SetEnabled(enabled)
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
			textBounds.Dy() + theme.Padding() + textBounds.Dx(),
			textBounds.Dy())
	}
	
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Checkbox) draw () {
	bounds := element.core.Bounds()
	boxBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dy())

	backgroundPattern, _ := theme.BackgroundPattern(theme.PatternState {
		Case: checkboxCase,
	})
	artist.FillRectangle ( element.core, backgroundPattern, bounds)

	pattern, inset := theme.ButtonPattern(theme.PatternState {
		Case: checkboxCase,
		Disabled: !element.Enabled(),
		Selected: element.Selected(),
		Pressed:  element.pressed,
	})
	artist.FillRectangle(element.core, pattern, boxBounds)

	textBounds := element.drawer.LayoutBounds()
	offset := image.Point {
		X: bounds.Dy() + theme.Padding(),
	}

	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground, _ := theme.ForegroundPattern (theme.PatternState {
		Case: checkboxCase,
		Disabled: !element.Enabled(),
	})
	element.drawer.Draw(element.core, foreground, offset)
	
	if element.checked {
		checkBounds := inset.Apply(boxBounds).Inset(2)
		artist.FillRectangle(element.core, foreground, checkBounds)
	}
}
