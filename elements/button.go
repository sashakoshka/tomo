package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"

// Button is a clickable button.
type Button struct {
	entity tomo.FocusableEntity
	drawer textdraw.Drawer

	enabled bool
	pressed bool
	text    string
	
	config config.Wrapped
	theme  theme.Wrapped

	showText bool
	hasIcon  bool
	iconId   tomo.Icon
	
	onClick func ()
}

// NewButton creates a new button with the specified label text.
func NewButton (text string) (element *Button) {
	element = &Button { showText: true, enabled: true }
	element.entity = tomo.NewEntity(element).(tomo.FocusableEntity)
	element.theme.Case = tomo.C("tomo", "button")
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.SetText(text)
	return
}

// Entity returns this element's entity.
func (element *Button) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Button) Draw (destination canvas.Canvas) {
	state   := element.state()
	bounds  := element.entity.Bounds()
	pattern := element.theme.Pattern(tomo.PatternButton, state)

	pattern.Draw(destination, bounds)
	
	foreground := element.theme.Color(tomo.ColorForeground, state)
	sink       := element.theme.Sink(tomo.PatternButton)
	margin     := element.theme.Margin(tomo.PatternButton)
	
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
		icon := element.theme.Icon(element.iconId, tomo.IconSizeSmall) 
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

// OnClick sets the function to be called when the button is clicked.
func (element *Button) OnClick (callback func ()) {
	element.onClick = callback
}

// Focus gives this element input focus.
func (element *Button) Focus () {
	if !element.entity.Focused() { element.entity.Focus() }
}

// Enabled returns whether this button is enabled or not.
func (element *Button) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this button can be clicked or not.
func (element *Button) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}

// SetText sets the button's label text.
func (element *Button) SetText (text string) {
	if element.text == text { return }
	element.text = text
	element.drawer.SetText([]rune(text))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetIcon sets the icon of the button. Passing theme.IconNone removes the
// current icon if it exists.
func (element *Button) SetIcon (id tomo.Icon) {
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
func (element *Button) ShowText (showText bool) {
	if element.showText == showText { return }
	element.showText = showText
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetTheme sets the element's theme.
func (element *Button) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawer.SetFace (element.theme.FontFace (
		tomo.FontStyleRegular,
		tomo.FontSizeNormal))
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *Button) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Button) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *Button) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	element.Focus()
	if button != input.ButtonLeft { return }
	element.pressed = true
	element.entity.Invalidate()
}

func (element *Button) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft { return }
	element.pressed = false
	within := image.Point { x, y }.In(element.entity.Bounds())
	if element.Enabled() && within && element.onClick != nil {
		element.onClick()
	}
	element.entity.Invalidate()
}

func (element *Button) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }
	if key == input.KeyEnter {
		element.pressed = true
		element.entity.Invalidate()
	}
}

func (element *Button) HandleKeyUp(key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.entity.Invalidate()
		if !element.Enabled() { return }
		if element.onClick != nil {
			element.onClick()
		}
	}
}

func (element *Button) updateMinimumSize () {
	padding := element.theme.Padding(tomo.PatternButton)
	margin  := element.theme.Margin(tomo.PatternButton)

	textBounds  := element.drawer.LayoutBounds()
	minimumSize := textBounds.Sub(textBounds.Min)
	
	if element.hasIcon {
		icon := element.theme.Icon(element.iconId, tomo.IconSizeSmall) 
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
	
	minimumSize = padding.Inverse().Apply(minimumSize)
	element.entity.SetMinimumSize(minimumSize.Dx(), minimumSize.Dy())
}

func (element *Button) state () tomo.State {
	return tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
		Pressed:  element.pressed,
	}
}
