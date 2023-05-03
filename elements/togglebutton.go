package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"

var toggleButtonCase = tomo.C("tomo", "toggleButton")

// ToggleButton is a togglable button.
type ToggleButton struct {
	entity tomo.Entity
	drawer textdraw.Drawer

	enabled bool
	pressed bool
	on      bool
	text    string

	showText bool
	hasIcon  bool
	iconId   tomo.Icon
	
	onToggle func ()
}

// NewToggleButton creates a new toggle button with the specified label text.
func NewToggleButton (text string, on bool) (element *ToggleButton) {
	element = &ToggleButton {
		showText: true,
		enabled:  true,
		on:       on,
	}
	element.entity = tomo.GetBackend().NewEntity(element)
	element.drawer.SetFace (element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal,
		toggleButtonCase))
	element.SetText(text)
	return
}

// Entity returns this element's entity.
func (element *ToggleButton) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *ToggleButton) Draw (destination artist.Canvas) {
	state   := element.state()
	bounds  := element.entity.Bounds()
	pattern := element.entity.Theme().Pattern(tomo.PatternButton, state, toggleButtonCase)
	
	lampPattern := element.entity.Theme().Pattern(tomo.PatternLamp, state, toggleButtonCase)
	lampPadding := element.entity.Theme().Padding(tomo.PatternLamp, toggleButtonCase).Horizontal()
	lampBounds  := bounds
	lampBounds.Max.X = lampBounds.Min.X + lampPadding
	bounds.Min.X += lampPadding

	pattern.Draw(destination, bounds)
	lampPattern.Draw(destination, lampBounds)
	
	foreground := element.entity.Theme().Color(tomo.ColorForeground, state, toggleButtonCase)
	sink       := element.entity.Theme().Sink(tomo.PatternButton, toggleButtonCase)
	margin     := element.entity.Theme().Margin(tomo.PatternButton, toggleButtonCase)
	
	offset := image.Pt (
		bounds.Dx() / 2,
		bounds.Dy() / 2).Add(bounds.Min)

	if element.showText {
		textBounds := element.drawer.LayoutBounds()
		offset.X -= textBounds.Dx() / 2
		offset.Y -= textBounds.Dy() / 2
		offset.Y -= textBounds.Min.Y
		offset.X -= textBounds.Min.X
	}

	if element.hasIcon {
		icon := element.entity.Theme().Icon(element.iconId, tomo.IconSizeSmall, toggleButtonCase) 
		if icon != nil {
			iconBounds := icon.Bounds()
			addedWidth := iconBounds.Dx()
			iconOffset := offset

			if element.showText {
				addedWidth += margin.X
			}
			
			iconOffset.X -= addedWidth / 2
			iconOffset.Y =
				bounds.Min.Y +
				(bounds.Dy() -
				iconBounds.Dy()) / 2
			if element.pressed {
				iconOffset = iconOffset.Add(sink)
			}
			offset.X += addedWidth / 2

			icon.Draw(destination, foreground, iconOffset)
		}
	}

	if element.showText {
		if element.pressed {
			offset = offset.Add(sink)
		}
		element.drawer.Draw(destination, foreground, offset)
	}
}

// OnToggle sets the function to be called when the button is toggled.
func (element *ToggleButton) OnToggle (callback func ()) {
	element.onToggle = callback
}

// Value reports whether or not the button is currently on.
func (element *ToggleButton) Value () (on bool) {
	return element.on
}

// Focus gives this element input focus.
func (element *ToggleButton) Focus () {
	if !element.entity.Focused() { element.entity.Focus() }
}

// Enabled returns whether this button is enabled or not.
func (element *ToggleButton) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this button can be toggled or not.
func (element *ToggleButton) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}

// SetText sets the button's label text.
func (element *ToggleButton) SetText (text string) {
	if element.text == text { return }
	element.text = text
	element.drawer.SetText([]rune(text))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetIcon sets the icon of the button. Passing theme.IconNone removes the
// current icon if it exists.
func (element *ToggleButton) SetIcon (id tomo.Icon) {
	if id == tomo.IconNone {
		element.hasIcon = false
	} else {
		if element.hasIcon && element.iconId == id { return }
		element.hasIcon = true
		element.iconId = id
	}
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// ShowText sets whether or not the button's text will be displayed.
func (element *ToggleButton) ShowText (showText bool) {
	if element.showText == showText { return }
	element.showText = showText
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *ToggleButton) HandleThemeChange () {
	element.drawer.SetFace (element.entity.Theme().FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal, toggleButtonCase))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *ToggleButton) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *ToggleButton) HandleMouseDown (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if !element.Enabled() { return }
	element.Focus()
	if button != input.ButtonLeft { return }
	element.pressed = true
	element.entity.Invalidate()
}

func (element *ToggleButton) HandleMouseUp (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if button != input.ButtonLeft { return }
	element.pressed = false
	within := position.In(element.entity.Bounds())
	if element.Enabled() && within {
		element.on = !element.on
		if element.onToggle != nil {
			element.onToggle()
		}
	}
	element.entity.Invalidate()
}

func (element *ToggleButton) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }
	if key == input.KeyEnter {
		element.pressed = true
		element.entity.Invalidate()
	}
}

func (element *ToggleButton) HandleKeyUp(key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.entity.Invalidate()
		if !element.Enabled() { return }
		element.on = !element.on
		if element.onToggle != nil {
			element.onToggle()
		}
	}
}

func (element *ToggleButton) updateMinimumSize () {
	padding     := element.entity.Theme().Padding(tomo.PatternButton, toggleButtonCase)
	margin      := element.entity.Theme().Margin(tomo.PatternButton, toggleButtonCase)
	lampPadding := element.entity.Theme().Padding(tomo.PatternLamp, toggleButtonCase)

	textBounds  := element.drawer.LayoutBounds()
	minimumSize := textBounds.Sub(textBounds.Min)
	
	if element.hasIcon {
		icon := element.entity.Theme().Icon(element.iconId, tomo.IconSizeSmall, toggleButtonCase) 
		if icon != nil {
			bounds := icon.Bounds()
			if element.showText {
				minimumSize.Max.X += bounds.Dx()
				minimumSize.Max.X += margin.X
			} else {
				minimumSize.Max.X = bounds.Dx()
			}
		}
	}

	minimumSize.Max.X += lampPadding.Horizontal()
	minimumSize = padding.Inverse().Apply(minimumSize)
	element.entity.SetMinimumSize(minimumSize.Dx(), minimumSize.Dy())
}

func (element *ToggleButton) state () tomo.State {
	return tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
		Pressed:  element.pressed,
		On:       element.on,
	}
}
