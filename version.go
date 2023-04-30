package tomo

import "fmt"

// Version represents a semantic version number.
type Version [3]int

// TODO: when 1.0 is released, remove the notices. remember to update
// CurrentVersion too!

// CurrentVersion returns the current Tomo/Nasin version. Note that until 1.0 is
// released, this does not mean much.
func CurrentVersion () Version {
	return Version { 0, 0, 0 }
}

// CompatibleABI returns whether or not two versions are compatible on a binary
// level. Note that until 1.0 is released, this does not mean much.
func (version Version) CompatibleABI (other Version) bool {
	return version[0] == other[0] && version[1] == other[1]
}

// CompatibleAPI returns whether or not two versions are compatible on a source
// code level. Note that until 1.0 is released, this does not mean much.
func (version Version) CompatibleAPI (other Version) bool {
	return version[0] == other[0]
}

// String returns a string representation of the version.
func (version Version) String () string {
	return fmt.Sprint(version[0], ".", version[1], ".", version[2])
}
