package tomo

import "errors"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/elements"

// Backend represents a connection to a display server, or something similar.
// It is capable of managing an event loop, and creating windows.
type Backend interface {
	// Run runs the backend's event loop. It must block until the backend
	// experiences a fatal error, or Stop() is called.
	Run () (err error)

	// Stop stops the backend's event loop.
	Stop ()

	// Do executes the specified callback within the main thread as soon as
	// possible. This method must be safe to call from other threads.
	Do (callback func ())

	// NewWindow creates a new window with the specified width and height,
	// and returns a struct representing it that fulfills the MainWindow
	// interface.
	NewWindow (width, height int) (window elements.MainWindow, err error)

	// Copy puts data into the clipboard.
	Copy (data.Data)

	// Paste returns the data currently in the clipboard.
	Paste (accept []data.Mime) (data.Data)
	
	// SetTheme sets the theme of all open windows.
	SetTheme (theme.Theme)
	
	// SetConfig sets the configuration of all open windows.
	SetConfig (config.Config)
}

// BackendFactory represents a function capable of constructing a backend
// struct. Any connections should be initialized within this function. If there
// any errors encountered during this process, the function should immediately
// stop, clean up any resources, and return an error.
type BackendFactory func () (backend Backend, err error)

// RegisterBackend registers a backend factory. When an application calls
// tomo.Run(), the first registered backend that does not throw an error will be
// used.
func RegisterBackend (factory BackendFactory) {
	factories = append(factories, factory)
}

var factories []BackendFactory

func instantiateBackend () (backend Backend, err error) {
	// find a suitable backend
	for _, factory := range factories {
		backend, err = factory()
		if err == nil && backend != nil { return }
	}

	// if none were found, but there was no error produced, produce an
	// error
	if err == nil {
		err = errors.New("no available backends")
	}

	return
}
