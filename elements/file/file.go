package fileElements

import "io/fs"
import "image"
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
	
	iconID     theme.Icon
	filesystem fs.StatFS
	location   string
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
	element.iconID = theme.IconError
	info, err := element.filesystem.Stat(element.location)
	if err != nil { return err }

	// TODO: choose icon based on file mime type
	if info.IsDir() {
		element.iconID = theme.IconDirectory
	} else {
		element.iconID = theme.IconFile
	}

	element.updateMinimumSize()
	element.drawAndPush()
	return err
}

func (element *File) OnChoose (callback func ()) {
	element.onChoose = callback
}

func (element *File) state () theme.State {
	return theme.State {
		Disabled: !element.Enabled(),
		Focused:  element.Focused(),
		// Pressed:  element.pressed,
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
		icon.Draw (
			element.core,
			element.theme.Color (
				theme.ColorForeground, state),
			bounds.Min.Add(offset))
	}
}
