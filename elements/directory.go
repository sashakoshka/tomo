package elements

import "image"
import "path/filepath"
import "git.tebibyte.media/sashakoshka/tomo"

// TODO: base on flow implementation of list. also be able to switch to a table
// variant for a more information dense view.

type historyEntry struct {
	location string
	filesystem ReadDirStatFS
}

// Directory displays a list of files within a particular directory and
// file system.
type Directory struct {
	*List
	history      []historyEntry
	historyIndex int
	onChoose     func (file string)
}

// NewDirectory creates a new directory view. If within is nil, it will use
// the OS file system.
func NewDirectory (
	location string,
	within ReadDirStatFS,
) (
	element *Directory,
	err error,
) {
	element = &Directory {
		List: NewList(8),
	}
	err = element.SetLocation(location, within)
	return
}

// Location returns the directory's location and filesystem.
func (element *Directory) Location () (string, ReadDirStatFS) {
	if len(element.history) < 1 { return "", nil }
	current := element.history[element.historyIndex]
	return current.location, current.filesystem
}

// SetLocation sets the directory's location and filesystem. If within is nil,
// it will use the OS file system.
func (element *Directory) SetLocation (
	location string,
	within ReadDirStatFS,
) error {
	if within == nil {
		within = defaultFS { }
	}
	element.scroll = image.Point { }

	if element.history != nil {
		element.historyIndex ++
	}
	element.history = append (
		element.history[:element.historyIndex],
		historyEntry { location, within })
	return element.Update()
}

// Backward goes back a directory in history
func (element *Directory) Backward () (bool, error) {
	if element.historyIndex > 1 {
		element.historyIndex --
		return true, element.Update()
	} else {
		return false, nil
	}
}

// Forward goes forward a directory in history
func (element *Directory) Forward () (bool, error) {
	if element.historyIndex < len(element.history) - 1 {
		element.historyIndex ++
		return true, element.Update()
	} else {
		return false, nil
	}
}

// Update refreshes the directory's contents.
func (element *Directory) Update () error {
	location, filesystem := element.Location()
	entries, err := filesystem.ReadDir(location)

	children := make([]tomo.Element, len(entries))
	for index, entry := range entries {
		filePath := filepath.Join(location, entry.Name())
		file, _ := NewFile(filePath, filesystem)
		file.OnChoose (func () {
			if element.onChoose != nil {
				element.onChoose(filePath)
			}
		})
		
		children[index] = file
	}
	
	element.DisownAll()
	element.Adopt(children...)
	return err
}

// OnChoose sets a function to be called when the user double-clicks a file or
// sub-directory within the directory view.
func (element *Directory) OnChoose (callback func (file string)) {
	element.onChoose = callback
}
