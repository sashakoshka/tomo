package tomo

import "os"
import "io"
import "path/filepath"
import "git.tebibyte.media/sashakoshka/tomo/dirs"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/elements"

var backend Backend

// Run initializes a backend, calls the callback function, and begins the event
// loop in that order. This function does not return until Stop() is called, or
// the backend experiences a fatal error.
func Run (callback func ()) (err error) {
	backend, err = instantiateBackend()
	if err != nil { return }
	config := parseConfig()
	backend.SetConfig(config)
	backend.SetTheme(parseTheme(config.ThemePath()))
	if callback != nil { callback() }
	err = backend.Run()
	backend = nil
	return
}

// Stop gracefully stops the event loop and shuts the backend down. Call this
// before closing your application.
func Stop () {
	if backend != nil { backend.Stop() }
}

// Do executes the specified callback within the main thread as soon as
// possible. This function can be safely called from other threads.
func Do (callback func ()) {
	assertBackend()
	backend.Do(callback)
}

// NewWindow creates a new window using the current backend, and returns it as a
// Window. If the window could not be created, an error is returned explaining
// why. If this function is called without a running backend, an error is
// returned as well.
func NewWindow (width, height int) (window elements.Window, err error) {
	assertBackend()
	return backend.NewWindow(width, height)
}

// Copy puts data into the clipboard.
func Copy (data data.Data) {
	assertBackend()
	backend.Copy(data)
}

// Paste returns the data currently in the clipboard. This method may
// return nil.
func Paste (accept []data.Mime) (data.Data) {
	assertBackend()
	return backend.Paste(accept)
}

// SetTheme sets the theme of all open windows.
func SetTheme (theme theme.Theme) {
	backend.SetTheme(theme)
}

// SetConfig sets the configuration of all open windows.
func SetConfig (config config.Config) {
	backend.SetConfig(config)
}

func parseConfig () (config.Config) {
	return parseMany [config.Config] (
		dirs.ConfigDirs("tomo/tomo.conf"),
		config.Parse,
		config.Default { })
}

func parseTheme (path string) (theme.Theme) {
	if path == "" { return theme.Default { } }
	path = filepath.Join(path, "tomo")
	
	// find all tomo pattern graph files in the directory
	directory, err := os.Open(path)
	if err != nil { return theme.Default { } }
	names, _ := directory.Readdirnames(0)
	paths := []string { }
	for _, name := range names {
		if filepath.Ext(name) == ".tpg" {
			paths = append(paths, filepath.Join(path, name))
		}
	}

	// parse them
	return parseMany [theme.Theme] (
		paths,
		theme.Parse,
		theme.Default { })
}

func parseMany [OBJECT any] (
	paths []string,
	parser func (...io.Reader) OBJECT,
	fallback OBJECT,
) (
	object OBJECT,
) {
	// convert all paths into readers
	sources := []io.Reader { }
	for _, path := range paths {
		file, err := os.Open(path)
		if err != nil { continue }
		sources = append(sources, file)
		defer file.Close()
	}
	
	if sources == nil {
		// if there are no readers, return the fallback object
		return fallback
	} else {
		// if there are readers, parse them
		return parser(sources...)
	}
}

func assertBackend () {
	if backend == nil { panic("no backend is running") }
}
