package fileElements

import "io/fs"
import "path/filepath"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"

type ReadDirStatFS interface {
	fs.ReadDirFS
	fs.StatFS
}

type DirectoryView struct {
	*basicElements.List

	filesystem ReadDirStatFS
	location string
	onChoose func (file string)
}

func NewDirectoryView (
	location string,
	within ReadDirStatFS,
) (
	element *DirectoryView,
	err error,
) {
	element = &DirectoryView {
		List: basicElements.NewList(),
	}
	err = element.SetLocation(location, within)
	return
}

func (element *DirectoryView) Location () (string, fs.ReadDirFS) {
	return element.location, element.filesystem
}

func (element *DirectoryView) SetLocation (
	location string,
	within ReadDirStatFS,
) error {
	if within == nil {
		within = defaultFS { }
	}
	element.location   = location
	element.filesystem = within
	return element.Update()
}

func (element *DirectoryView) Update () error {
	entries, err := element.filesystem.ReadDir(element.location)

	listEntries := make([]basicElements.ListEntry, len(entries))
	for index, entry := range entries {
		filePath := filepath.Join(element.location, entry.Name())
		listEntries[index] = basicElements.NewListEntry (
			entry.Name(),
			func () {
				filePath := filePath
				if element.onChoose != nil {
					element.onChoose(filePath)
				}
			})
	}
	element.Clear()
	element.Append(listEntries...)
	
	return err
}

func (element *DirectoryView) OnChoose (callback func (file string)) {
	element.onChoose = callback
}
