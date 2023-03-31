package tomo

var backend Backend

// Run initializes a backend, calls the callback function, and begins the event
// loop in that order. This function does not return until Stop() is called, or
// the backend experiences a fatal error.
func Run (callback func ()) (err error) {
	backend, err = instantiateBackend()
	if err != nil { return }
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
// MainWindow. If the window could not be created, an error is returned
// explaining why.
func NewWindow (width, height int) (window MainWindow, err error) {
	assertBackend()
	return backend.NewWindow(width, height)
}

// SetTheme sets the theme of all open windows.
func SetTheme (theme Theme) {
	backend.SetTheme(theme)
}

// SetConfig sets the configuration of all open windows.
func SetConfig (config Config) {
	backend.SetConfig(config)
}

func assertBackend () {
	if backend == nil { panic("no backend is running") }
}
