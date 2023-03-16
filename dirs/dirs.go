// Package dirs provides access to standard system and user directories.
package dirs

import "os"
import "strings"
import "path/filepath"

var homeDirectory string
var configHome    string
var configDirs    []string
var dataHome      string
var dataDirs      []string
var cacheHome     string

func init () {
	var err error
	homeDirectory, err = os.UserHomeDir()
	if err != nil {
		panic("could not get user home directory: " + err.Error())
	}

	configHome = os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		configHome = filepath.Join(homeDirectory, "/.config/")
	}
	
	configDirsString := os.Getenv("XDG_CONFIG_DIRS")
	if configDirsString == "" {
		configDirsString = "/etc/xdg/"
	}
	configDirs = append(strings.Split(configDirsString, ":"), configHome)

	dataHome = os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(homeDirectory, "/.local/share/")
	}
	
	dataDirsString := os.Getenv("XDG_CONFIG_DIRS")
	if dataDirsString == "" {
		dataDirsString = "/usr/local/share/:/usr/share/"
	}
	configDirs = append(strings.Split(configDirsString, ":"), configHome)

	cacheHome = os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		cacheHome = filepath.Join(homeDirectory, "/.cache/")
	}
}

// ConfigHome returns the path to the directory where user configuration files
// should be stored.
func ConfigHome (name string) (home string) {
	return filepath.Join(configHome, name)
}

// ConfigDirs returns all paths where configuration files might exist.
func ConfigDirs (name string) (dirs []string) {
	dirs = make([]string, len(configDirs))
	for index, dir := range configDirs {
		dirs[index] = filepath.Join(dir, name)
	}
	return
}

// DataHome returns the path to the directory where user data should be stored.
func DataHome (name string) (home string) {
	return filepath.Join(dataHome, name)
}

// DataDirs returns all paths where data files might exist.
func DataDirs (name string) (dirs []string) {
	dirs = make([]string, len(dataDirs))
	for index, dir := range dataDirs {
		dirs[index] = filepath.Join(dir, name)
	}
	return
}

// CacheHome returns the path to the directory where user cache files should be
// stored.
func CacheHome (name string) (home string) {
	return filepath.Join(cacheHome, name)
}
