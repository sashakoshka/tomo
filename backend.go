package tomo

import "image"

// Backend represents a connection to a display server, or something similar.
// It is capable of managing an event loop, and creating windows.
type Backend interface {
	// Run runs the backend's event loop. It must block until the backend
	// experiences a fatal error, or Stop() is called.
	Run () error

	// Stop stops the backend's event loop.
	Stop ()

	// Do executes the specified callback within the main thread as soon as
	// possible. This method must be safe to call from other threads.
	Do (callback func ())

	// NewEntity creates a new entity for the specified element.
	NewEntity (owner Element) Entity

	// NewWindow creates a new window within the specified bounding
	// rectangle. The position on screen may be overridden by the backend or
	// operating system.
	NewWindow (bounds image.Rectangle) (MainWindow, error)
	
	// SetTheme sets the theme of all open windows.
	SetTheme (Theme)
	
	// SetConfig sets the configuration of all open windows.
	SetConfig (Config)
}

var backend Backend

// GetBackend returns the currently running backend.
func GetBackend () Backend {
	return backend
}

// SetBackend sets the currently running backend. The backend can only be set
// once—if there already is one then this function will do nothing.
func SetBackend (b Backend) {
	if backend != nil { return }
	backend = b
}

// Bounds creates a rectangle from an x, y, width, and height.
func Bounds (x, y, width, height int) image.Rectangle {
	return image.Rect(x, y, x + width, y + height)
}
