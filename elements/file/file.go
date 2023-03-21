package fileElements

import "io/fs"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"

// File is a 
type File struct {
	*basicElements.Icon

	// we inherit from Icon directly because it is not our responsibility
	// to draw text. this will be the responsibility of the directory that
	// contains the file. we don't handle mouse events on the file label
	// text either because when the user clicks on that we want to rename
	// the file.

	filesystem fs.StatFS
	location string
	onChoose func ()
}

func NewFile (
	location string,
	within fs.StatFS,
) (
	element *File,
	err error,
) {
	element = &File {
		Icon: basicElements.NewIcon(theme.IconFile, theme.IconSizeLarge),
	}
	err = element.SetLocation(location, within)
	return
}

func (element *File) Location () (string, fs.StatFS) {
	return element.location, element.filesystem
}

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

func (element *File) Update () error {
	info, err := element.filesystem.Stat(element.location)
	if err != nil { return err }

	if info.IsDir() {
		element.SetIcon(theme.IconDirectory, theme.IconSizeLarge)
	} else {
		element.SetIcon(theme.IconFile, theme.IconSizeLarge)
	}
	
	return err
}

func (element *File) OnChoose (callback func ()) {
	element.onChoose = callback
}
