package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// Checkbox is a toggle-able checkbox with a label.
type Checkbox struct {
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

// NewCheckbox creates a new cbeckbox with the specified label text.
func NewCheckbox (text string, checked bool) (element *Checkbox) {
	element = &Checkbox { checked: checked, enabled: true }
	element.entity = tomo.NewEntity(element).(tomo.FocusableEntity)
	element.theme.Case = tomo.C("tomo", "checkbox")
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.SetText(text)
	return
}

// Entity returns this element's entity.
func (element *Checkbox) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Checkbox) Draw (destination canvas.Canvas) {
	bounds := element.entity.Bounds()
	boxBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dy()).Add(bounds.Min)

	state := tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
		Pressed:  element.pressed,
		On:       element.checked,
	}

	element.entity.DrawBackground(destination)
		
	pattern := element.theme.Pattern(tomo.PatternButton, state)
	pattern.Draw(destination, boxBounds)

	textBounds := element.drawer.LayoutBounds()
	margin := element.theme.Margin(tomo.PatternBackground)
	offset := bounds.Min.Add(image.Point {
		X: bounds.Dy() + margin.X,
	})

	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X

	foreground := element.theme.Color(tomo.ColorForeground, state)
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

// SetTheme sets the element's theme.
func (element *Checkbox) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *Checkbox) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Checkbox) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	element.Focus()
	element.pressed = true
	element.entity.Invalidate()
}

func (element *Checkbox) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft || !element.pressed { return }

	element.pressed = false
	within := image.Point { x, y }.In(element.entity.Bounds())
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
		margin := element.theme.Margin(tomo.PatternBackground)
		element.entity.SetMinimumSize (
			textBounds.Dy() + margin.X + textBounds.Dx(),
			textBounds.Dy())
	}
}
