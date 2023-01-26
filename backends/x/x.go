package x

import "git.tebibyte.media/sashakoshka/tomo"

import "github.com/jezek/xgbutil"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/xevent"

// Backend is an instance of an X backend.
type Backend struct {
	connection *xgbutil.XUtil

	doChannel chan(func ())

	modifierMasks struct {
		capsLock   uint16
		shiftLock  uint16
		numLock    uint16
		modeSwitch uint16

		alt   uint16
		meta  uint16
		super uint16
		hyper uint16
	}

	windows map[xproto.Window] *Window
}

// NewBackend instantiates an X backend.
func NewBackend () (output tomo.Backend, err error) {
	backend := &Backend {
		windows: map[xproto.Window] *Window { },
		doChannel: make(chan func (), 0),
	}
	
	// connect to X
	backend.connection, err = xgbutil.NewConn()
	if err != nil { return }
	backend.initializeKeymapInformation()

	output = backend
	return
}

// Run runs the backend's event loop. This method will not exit until Stop() is
// called, or the backend experiences a fatal error.
func (backend *Backend) Run () (err error) {
	backend.assert()
	pingBefore,
	pingAfter,
	pingQuit := xevent.MainPing(backend.connection)
	for {
		select {
		case <- pingBefore:
			<- pingAfter
		case callback := <- backend.doChannel:
			callback()
		case <- pingQuit:
			return
		}
	}
}

// Stop gracefully closes the connection and stops the event loop.
func (backend *Backend) Stop () {
	backend.assert()
	for _, window := range backend.windows {
		window.Close()
	}
	xevent.Quit(backend.connection)
	backend.connection.Conn().Close()
}

// Do executes the specified callback within the main thread as soon as
// possible. This function can be safely called from other threads.
func (backend *Backend) Do (callback func ()) {
	backend.assert()
	backend.doChannel <- callback
}

// Copy puts data into the clipboard. This method is not yet implemented and
// will do nothing!
func (backend *Backend) Copy (data tomo.Data) {
	backend.assert()
	// TODO
}

// Paste returns the data currently in the clipboard. This method may
// return nil. This method is not yet implemented and will do nothing!
func (backend *Backend) Paste () (data tomo.Data) {
	backend.assert()
	// TODO
	return
}

func (backend *Backend) assert () {
	if backend == nil { panic("nil backend") }
}

func init () {
	tomo.RegisterBackend(NewBackend)
}
