package fileElements

import "os"
import "io/fs"

type defaultFS struct { }

func (defaultFS) Open (name string) (fs.File, error) {
	return os.Open(name)
}

func (defaultFS) ReadDir (name string) ([]DirEntry, error) {
	return os.ReadDir(name)
}

func (defaultFS) Stat (name string) (FileInfo, error) {
	return os.Stat(name)
}
