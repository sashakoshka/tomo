package elements

import "image"
import "tomo"
import "tomo/input"
import "art"
import "tomo/textdraw"

var checkboxCase = tomo.C("tomo", "checkbox")

// Checkbox is a toggle-able checkbox with a label.
type Checkbox struct {
	entity tomo.Entity
	drawer textdraw.Drawer

	enabled bool
	pressed bool
	checked bool
	text    string
	
	onToggle func ()
}

// NewCheckbox creates a new cbeckbox with the specified label text.
func NewCheckbox (text string, checked bool) (element *Checkbox) {
	element = &Checkbox { checked: checked, enabled: true }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.drawer.SetFace (element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal,
		checkboxCase))
	element.SetText(text)
	return
}

// Entity returns this element's entity.
func (element *Checkbox) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Checkbox) Draw (destination art.Canvas) {
	bounds := element.entity.Bounds()
	boxBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dy()).Add(bounds.Min)

	state := tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
		Pressed:  element.pressed,
		On:       element.checked,
	}

	element.entity.DrawBackground(destination)
		
	pattern := element.entity.Theme().Pattern(tomo.PatternButton, state, checkboxCase)
	pattern.Draw(destination, boxBounds)

	textBounds := element.drawer.LayoutBounds()
	margin := element.entity.Theme().Margin(tomo.PatternBackground, checkboxCase)
	offset := bounds.Min.Add(image.Point {
		X: bounds.Dy() + margin.X,
	})

	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground := element.entity.Theme().Color(tomo.ColorForeground, state, checkboxCase)
	element.drawer.Draw(destination, foreground, offset)
}

// OnToggle sets the function to be called when the checkbox is toggled.
func (element *Checkbox) OnToggle (callback func ()) {
	element.onToggle = callback
}

// Value reports whether or not the checkbox is currently checked.
func (element *Checkbox) Value () (checked bool) {
	return element.checked
}

// Focus gives this element input focus.
func (element *Checkbox) Focus () {
	if !element.entity.Focused() { element.entity.Focus() }
}

// Enabled returns whether this checkbox is enabled or not.
func (element *Checkbox) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this checkbox can be toggled or not.
func (element *Checkbox) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}

// SetText sets the checkbox's label text.
func (element *Checkbox) SetText (text string) {
	if element.text == text { return }
	element.text = text
	element.drawer.SetText([]rune(text))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Checkbox) HandleThemeChange () {
	element.drawer.SetFace (element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal,
		checkboxCase))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Checkbox) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *Checkbox) HandleMouseDown (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if !element.Enabled() { return }
	element.Focus()
	element.pressed = true
	element.entity.Invalidate()
}

func (element *Checkbox) HandleMouseUp (
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

func (element *Checkbox) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter {
		element.pressed = true
		element.entity.Invalidate()
	}
}

func (element *Checkbox) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.checked = !element.checked
		element.entity.Invalidate()
		if element.onToggle != nil {
			element.onToggle()
		}
	}
}

func (element *Checkbox) updateMinimumSize () {
	textBounds := element.drawer.LayoutBounds()
	if element.text == "" {
		element.entity.SetMinimumSize(textBounds.Dy(), textBounds.Dy())
	} else {
		margin := element.entity.Theme().Margin(tomo.PatternBackground, checkboxCase)
		element.entity.SetMinimumSize (
			textBounds.Dy() + margin.X + textBounds.Dx(),
			textBounds.Dy())
	}
}
