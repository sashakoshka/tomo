package fileElements

import "os"
import "io/fs"

// ReadDirStatFS is a combination of fs.ReadDirFS and fs.StatFS. It is the
// minimum filesystem needed to satisfy a directory view.
type ReadDirStatFS interface {
	fs.ReadDirFS
	fs.StatFS
}

type defaultFS struct { }

func (defaultFS) Open (name string) (fs.File, error) {
	return os.Open(name)
}

func (defaultFS) ReadDir (name string) ([]fs.DirEntry, error) {
	return os.ReadDir(name)
}

func (defaultFS) Stat (name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
