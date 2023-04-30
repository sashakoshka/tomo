package nasin

import "image"
import "errors"
import "git.tebibyte.media/sashakoshka/tomo"

// Application represents a Tomo/Nasin application.
type Application interface {
	Init () error
}

// Run initializes Tomo and Nasin, and runs the given application. This function
// will block until the application exits or a fatal error occurrs.
func Run (application Application) {
	loadPlugins()

	backend, err := instantiateBackend()
	if err != nil {
		println("nasin: cannot start application:", err.Error())
		return
	}
	tomo.SetBackend(backend)
	
	if application == nil { panic("nasin: nil application") }
	err = application.Init()
	if err != nil {
		println("nasin: backend exited with error:", err.Error())
		return
	}
	
	err = backend.Run()
	if err != nil {
		println("nasin: backend exited with error:", err.Error())
		return
	}
	return
}

// Stop stops the event loop
func Stop () {
	assertBackend()
	tomo.GetBackend().Stop()
}

// Do executes the specified callback within the main thread as soon as
// possible.
func Do (callback func()) {
	assertBackend()
	tomo.GetBackend().Do(callback)
}

// NewWindow creates a new window within the specified bounding rectangle. The
// position on screen may be overridden by the backend or operating system.
func NewWindow (bounds image.Rectangle) (tomo.MainWindow, error) {
	assertBackend()
	return tomo.GetBackend().NewWindow(bounds)
}

func assertBackend () {
	if tomo.GetBackend() == nil {
		panic("nasin: no running tomo backend")
	}
}

func instantiateBackend () (backend tomo.Backend, err error) {
	// find a suitable backend
	for _, factory := range factories {
		backend, err = factory()
		if err == nil && backend != nil { return }
	}

	// if none were found, but there was no error produced, produce an error
	if err == nil {
		return nil, errors.New("no available tomo backends")
	}

	return
}
