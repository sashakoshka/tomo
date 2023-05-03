package elements

import "time"
import "io/fs"
import "image"
import "tomo"
import "tomo/input"
import "tomo/artist"

var fileCase = tomo.C("files", "file")

// File displays an interactive visual representation of a file within any
// file system.
type File struct {
	entity tomo.Entity

	lastClick  time.Time
	pressed    bool
	enabled    bool
	iconID     tomo.Icon
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
	element = &File { enabled: true }
	element.entity = tomo.GetBackend().NewEntity(element)
	err = element.SetLocation(location, within)
	return
}

// Entity returns this element's entity.
func (element *File) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *File) Draw (destination artist.Canvas) {
	// background
	state  := element.state()
	bounds := element.entity.Bounds()
	sink   := element.entity.Theme().Sink(tomo.PatternButton, fileCase)
	element.entity.Theme().
		Pattern(tomo.PatternButton, state, fileCase).
		Draw(destination, bounds)

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
			destination,
			element.entity.Theme().Color(tomo.ColorForeground, state, fileCase),
			bounds.Min.Add(offset))
	}
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
		element.iconID = tomo.IconError
	} else if info.IsDir() {
		element.iconID = tomo.IconDirectory
	} else {
		// TODO: choose icon based on file mime type
		element.iconID = tomo.IconFile
	}

	element.updateMinimumSize()
	element.entity.Invalidate()
	return err
}

func (element *File) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if !element.Enabled() { return }
	if key == input.KeyEnter {
		element.pressed = true
		element.entity.Invalidate()
	}
}

func (element *File) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	if key == input.KeyEnter && element.pressed {
		element.pressed = false
		if !element.Enabled() { return }
		element.entity.Invalidate()
		if element.onChoose != nil {
			element.onChoose()
		}
	}
}

func (element *File) HandleFocusChange () {
	element.entity.Invalidate()
}

func (element *File) HandleSelectionChange () {
	element.entity.Invalidate()
}

func (element *File) OnChoose (callback func ()) {
	element.onChoose = callback
}

// Focus gives this element input focus.
func (element *File) Focus () {
	if !element.entity.Focused() { element.entity.Focus() }
}

// Enabled returns whether this file is enabled or not.
func (element *File) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this file is enabled or not.
func (element *File) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
}

func (element *File) HandleMouseDown  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if !element.Enabled() { return }
	if !element.entity.Focused() { element.Focus() }
	if button != input.ButtonLeft { return }
	element.pressed = true
	element.entity.Invalidate()
}

func (element *File) HandleMouseUp  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if button != input.ButtonLeft { return }
	element.pressed = false
	within := position.In(element.entity.Bounds())
	if time.Since(element.lastClick) < element.entity.Config().DoubleClickDelay() {
		if element.Enabled() && within && element.onChoose != nil {
			element.onChoose()
		}
	} else {
		element.lastClick = time.Now()
	}
	element.entity.Invalidate()
}

func (element *File) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *File) state () tomo.State {
	return tomo.State {
		Disabled: !element.Enabled(),
		Focused:  element.entity.Focused(),
		Pressed:  element.pressed,
		On:       element.entity.Selected(),
	}
}

func (element *File) icon () artist.Icon {
	return element.entity.Theme().Icon(element.iconID, tomo.IconSizeLarge, fileCase)
}

func (element *File) updateMinimumSize () {
	padding := element.entity.Theme().Padding(tomo.PatternButton, fileCase)
	icon := element.icon()
	if icon == nil {
		element.entity.SetMinimumSize (
			padding.Horizontal(),
			padding.Vertical())
	} else {
		bounds := padding.Inverse().Apply(icon.Bounds())
		element.entity.SetMinimumSize(bounds.Dx(), bounds.Dy())
	}
}
