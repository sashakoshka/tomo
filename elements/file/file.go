package fileElements

import "time"
import "io/fs"
import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// File displays an interactive visual representation of a file within any
// file system.
type File struct {
	*core.Core
	*core.FocusableCore
	core core.CoreControl
	focusableControl core.FocusableCoreControl
	
	config config.Wrapped
	theme  theme.Wrapped

	lastClick  time.Time
	pressed    bool
	iconID     theme.Icon
	filesystem fs.StatFS
	location   string
	selected   bool
	
	onChoose   func ()
}

// NewFile creates a new file element. If within is nil, it will use the OS file
// system
func NewFile (
	location string,
	within fs.StatFS,
) (
	element *File,
	err error,
) {
	element = &File { }
	element.theme.Case = theme.C("files", "file")
	element.Core, element.core = core.NewCore(element, element.drawAll)
	element.FocusableCore,
	element.focusableControl = core.NewFocusableCore(element.core, element.drawAndPush)
	err = element.SetLocation(location, within)
	return
}

// Location returns the file's location and filesystem.
func (element *File) Location () (string, fs.StatFS) {
	return element.location, element.filesystem
}

// SetLocation sets the file's location and filesystem. If within is nil, it
// will use the OS file system.
func (element *File) SetLocation (
	location string,
	within fs.StatFS,
) error {
	if within == nil {
		within = defaultFS { }
	}
	element.location   = location
	element.filesystem = within
	return element.Update()
}

// Update refreshes the element to match the file it represents.
func (element *File) Update () error {
	info, err := element.filesystem.Stat(element.location)

	if err != nil {
		element.iconID = theme.IconError
	} else if info.IsDir() {
		element.iconID = theme.IconDirectory
	} else {
		// TODO: choose icon based on file mime type
		element.iconID = theme.IconFile
	}

	element.updateMinimumSize()
	element.drawAndPush()
	return err
}

func (element *File) Selected () bool {
	return element.selected
}

func (element *File) SetSelected (selected bool) {
	if element.selected == selected { return }
	element.selected = selected
	element.drawAndPush()
}

func (element *File) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }
	if key == input.KeyEnter {
		element.pressed = true
		element.drawAndPush()
	}
}

func (element *File) HandleKeyUp(key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		element.drawAndPush()
		if !element.Enabled() { return }
		if element.onChoose != nil {
			element.onChoose()
		}
	}
}

func (element *File) OnChoose (callback func ()) {
	element.onChoose = callback
}

func (element *File) HandleMouseDown (x, y int, button input.Button) {
	if !element.Enabled() { return }
	if !element.Focused() { element.Focus() }
	if button != input.ButtonLeft { return }
	element.pressed = true
	element.drawAndPush()
}

func (element *File) HandleMouseUp (x, y int, button input.Button) {
	if button != input.ButtonLeft { return }
	element.pressed = false
	within := image.Point { x, y }.
		In(element.Bounds())
	if time.Since(element.lastClick) < time.Second / 2 {
		if element.Enabled() && within && element.onChoose != nil {
			element.onChoose()
		}
	} else {
		element.lastClick = time.Now()
	}
	element.drawAndPush()
}

// SetTheme sets the element's theme.
func (element *File) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawAndPush()
}

// SetConfig sets the element's configuration.
func (element *File) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.drawAndPush()
}

func (element *File) state () theme.State {
	return theme.State {
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
		Pressed:  element.pressed,
		On:       element.selected,
	}
}

func (element *File) icon () artist.Icon {
	return element.theme.Icon(element.iconID, theme.IconSizeLarge)
}

func (element *File) updateMinimumSize () {
	padding := element.theme.Padding(theme.PatternButton)
	icon := element.icon()
	if icon == nil {
		element.core.SetMinimumSize (
			padding.Horizontal(),
			padding.Vertical())
	} else {
		bounds := padding.Inverse().Apply(icon.Bounds())
		element.core.SetMinimumSize(bounds.Dx(), bounds.Dy())
	}
}

func (element *File) drawAndPush () {
	if element.core.HasImage() {
		element.drawAll()
		element.core.DamageAll()
	}
}

func (element *File) drawAll () {
	// background
	state  := element.state()
	bounds := element.Bounds()
	sink   := element.theme.Sink(theme.PatternButton)
	element.theme.
		Pattern(theme.PatternButton, state).
		Draw(element.core, bounds)

	// icon
	icon := element.icon()
	if icon != nil {
		iconBounds := icon.Bounds()
		offset := image.Pt (
			(bounds.Dx() - iconBounds.Dx()) / 2,
			(bounds.Dy() - iconBounds.Dy()) / 2)
		if element.pressed {
			offset = offset.Add(sink)
		}
		icon.Draw (
			element.core,
			element.theme.Color (
				theme.ColorForeground, state),
			bounds.Min.Add(offset))
	}
}
