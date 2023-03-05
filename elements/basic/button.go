package basicElements

import "image"
// import "runtime/debug"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
// import "git.tebibyte.media/sashakoshka/tomo/artist"
// import "git.tebibyte.media/sashakoshka/tomo/shatter"
import "git.tebibyte.media/sashakoshka/tomo/textdraw"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Button is a clickable button.
type Button struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl
	drawer textdraw.Drawer

	pressed bool
	text    string
	
	config config.Wrapped
	theme  theme.Wrapped

	showText bool
	hasIcon  bool
	iconId   theme.Icon
	
	onClick func ()
}

// NewButton creates a new button with the specified label text.
func NewButton (text string) (element *Button) {
	element = &Button { showText: true }
	element.theme.Case = theme.C("basic", "button")
	element.Core, element.core = core.NewCore(element.drawAll)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore(element.drawAndPush)
	element.SetText(text)
	return
}

func (element *Button) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	if !element.Focused() { element.Focus() }
	if button != input.ButtonLeft { return }
	element.pressed = true
	element.drawAndPush()
}

func (element *Button) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft { return }
	element.pressed = false
	within := image.Point { x, y }.
		In(element.Bounds())
	if element.Enabled() && within && element.onClick != nil {
		element.onClick()
	}
	element.drawAndPush()
}

func (element *Button) HandleMouseMove (x, y int) { }
func (element *Button) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Button) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }
	if key == input.KeyEnter {
		element.pressed = true
		element.drawAndPush()
	}
}

func (element *Button) HandleKeyUp(key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.drawAndPush()
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
	element.updateMinimumSize()
	element.drawAndPush()
}

// SetIcon sets the icon of the button.
func (element *Button) SetIcon (id theme.Icon) {
	if element.hasIcon && element.iconId == id { return }
	element.hasIcon = true
	element.iconId = id
	element.updateMinimumSize()
	element.drawAndPush()
}

// ClearIcon removes the button's icon, if it exists.
func (element *Button) ClearIcon () {
	if !element.hasIcon { return }
	element.hasIcon = false
	element.updateMinimumSize()
	element.drawAndPush()
}

// ShowText sets whether or not the button's text will be displayed.
func (element *Button) ShowText (showText bool) {
	if element.showText == showText { return }
	element.showText = showText
	element.updateMinimumSize()
	element.drawAndPush()
}

// SetTheme sets the element's theme.
func (element *Button) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawer.SetFace (element.theme.FontFace (
		theme.FontStyleRegular,
		theme.FontSizeNormal))
	element.updateMinimumSize()
	element.drawAndPush()
}

// SetConfig sets the element's configuration.
func (element *Button) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.drawAndPush()
}

func (element *Button) updateMinimumSize () {
	padding     := element.theme.Padding(theme.PatternButton)
	margin      := element.theme.Margin(theme.PatternButton)
	var minimumSize image.Rectangle

	if element.showText {
		minimumSize = element.drawer.LayoutBounds()
	}
	
	minimumSize = minimumSize.Sub(minimumSize.Min)
	
	if element.hasIcon {
		icon := element.theme.Icon(element.iconId, theme.IconSizeSmall) 
		if icon != nil {
			bounds := icon.Bounds()
			minimumSize.Max.X += bounds.Dx()
			if element.showText {
				minimumSize.Max.X += margin.X
			}
			if minimumSize.Max.Y < bounds.Dy() {
				minimumSize.Max.Y = bounds.Dy()
			}
		}
	}
	
	minimumSize = padding.Inverse().Apply(minimumSize)
	element.core.SetMinimumSize(minimumSize.Dx(), minimumSize.Dy())
}

func (element *Button) state () theme.State {
	return theme.State {
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
		Pressed:  element.pressed,
	}
}

func (element *Button) drawAndPush () {
	if element.core.HasImage () {
		element.drawAll()
		element.core.DamageAll()
	}
}

func (element *Button) drawAll () {
	element.drawBackground()
	element.drawText()
}

func (element *Button) drawBackground () []image.Rectangle {
	state   := element.state()
	bounds  := element.Bounds()
	pattern := element.theme.Pattern(theme.PatternButton, state)

	pattern.Draw(element.core, bounds)
	return []image.Rectangle { bounds }
}

func (element *Button) drawText () {
	state      := element.state()
	bounds     := element.Bounds()
	foreground := element.theme.Color(theme.ColorForeground, state)
	sink       := element.theme.Sink(theme.PatternButton)
	margin     := element.theme.Margin(theme.PatternButton)
	
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
		icon := element.theme.Icon(element.iconId, theme.IconSizeSmall) 
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

			icon.Draw(element.core, foreground, iconOffset)
		}
	}

	if element.showText {
		if element.pressed {
			offset = offset.Add(sink)
		}
		element.drawer.Draw(element.core, foreground, offset)
	}
}
