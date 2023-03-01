package basicElements

import "image"
// import "runtime/debug"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/shatter"
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
	
	onClick func ()
}

// NewButton creates a new button with the specified label text.
func NewButton (text string) (element *Button) {
	element = &Button { }
	element.theme.Case = theme.C("basic", "button")
	element.Core, element.core = core.NewCore(element.drawAll)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore (func () {
		element.drawAndPush(true)
	})
	element.SetText(text)
	return
}

func (element *Button) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	if !element.Focused() { element.Focus() }
	if button != input.ButtonLeft { return }
	element.pressed = true
	element.drawAndPush(true)
}

func (element *Button) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft { return }
	element.pressed = false
	within := image.Point { x, y }.
		In(element.Bounds())
	if element.Enabled() && within && element.onClick != nil {
		element.onClick()
	}
	element.drawAndPush(true)
}

func (element *Button) HandleMouseMove (x, y int) { }
func (element *Button) HandleMouseScroll (x, y int, deltaX, deltaY float64) { }

func (element *Button) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }
	if key == input.KeyEnter {
		element.pressed = true
		element.drawAndPush(true)
	}
}

func (element *Button) HandleKeyUp(key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.drawAndPush(true)
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
	element.drawAndPush(false)
}

// SetTheme sets the element's theme.
func (element *Button) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawer.SetFace (element.theme.FontFace (
		theme.FontStyleRegular,
		theme.FontSizeNormal))
	element.updateMinimumSize()
	element.drawAndPush(false)
}

// SetConfig sets the element's configuration.
func (element *Button) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.drawAndPush(false)
}

func (element *Button) updateMinimumSize () {
	textBounds := element.drawer.LayoutBounds()
	padding    := element.theme.Padding(theme.PatternButton)
	minimumSize := padding.Inverse().Apply(textBounds)
	element.core.SetMinimumSize(minimumSize.Dx(), minimumSize.Dy())
}

func (element *Button) drawAndPush (partial bool) {
	if element.core.HasImage () {
		if partial {
			element.core.DamageRegion (append (
				element.drawBackground(true),
				element.drawText(true))...)
		} else {
			element.drawAll()
			element.core.DamageAll()
		}
	}
}

func (element *Button) state () theme.State {
	return theme.State {
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
		Pressed:  element.pressed,
	}
}

func (element *Button) drawBackground (partial bool) []image.Rectangle {
	state   := element.state()
	bounds  := element.Bounds()
	pattern := element.theme.Pattern(theme.PatternButton, state)
	static  := element.theme.Hints(theme.PatternButton).StaticInset

	if partial && static != (artist.Inset { }) {
		tiles := shatter.Shatter(bounds, static.Apply(bounds))
		artist.Draw(element.core, pattern, tiles...)
		return tiles
	} else {
		pattern.Draw(element.core, bounds)
		return []image.Rectangle { bounds }
	}
}

func (element *Button) drawText (partial bool) image.Rectangle {
	state      := element.state()
	bounds     := element.Bounds()
	foreground := element.theme.Color(theme.ColorForeground, state)
	sink       := element.theme.Sink(theme.PatternButton)

	textBounds := element.drawer.LayoutBounds()
	offset := image.Point {
		X: bounds.Min.X + (bounds.Dx() - textBounds.Dx()) / 2,
		Y: bounds.Min.Y + (bounds.Dy() - textBounds.Dy()) / 2,
	}
	offset.Y -= textBounds.Min.Y
	offset.X -= textBounds.Min.X
	region := textBounds.Union(textBounds.Add(sink)).Add(offset)
	
	if element.pressed {
		offset = offset.Add(sink)
	}

	if partial {
		pattern := element.theme.Pattern(theme.PatternButton, state)
		pattern.Draw(element.core, region)
	}
	
	element.drawer.Draw(element.core, foreground, offset)
	return region
}

func (element *Button) drawAll () {
	element.drawBackground(false)
	element.drawText(false)
}
