package x

import "git.tebibyte.media/sashakoshka/tomo"

import "github.com/jezek/xgbutil"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/xevent"
import "github.com/jezek/xgbutil/keybind"
import "github.com/jezek/xgbutil/mousebind"

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

	theme  tomo.Theme
	config tomo.Config

	windows map[xproto.Window] *window

	open bool
}

// NewBackend instantiates an X backend.
func NewBackend () (output tomo.Backend, err error) {
	backend := &Backend {
		windows:   map[xproto.Window] *window { },
		doChannel: make(chan func (), 32),
		open:      true,
	}
	
	// connect to X
	backend.connection, err = xgbutil.NewConn()
	if err != nil { return }
	backend.initializeKeymapInformation()

	keybind.Initialize(backend.connection)
	mousebind.Initialize(backend.connection)

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
		for _, window := range backend.windows {
			window.system.afterEvent()
		}
	}
}

// Stop gracefully closes the connection and stops the event loop.
func (backend *Backend) Stop () {
	backend.assert()
	if !backend.open { return }
	backend.open = false
	
	toClose := []*window { }
	for _, window := range backend.windows {
		toClose = append(toClose, window)
	}
	for _, window := range toClose {
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

// SetTheme sets the theme of all open windows.
func (backend *Backend) SetTheme (theme tomo.Theme) {
	backend.assert()
	backend.theme = theme
	for _, window := range backend.windows {
		window.SetTheme(theme)
	}
}

// SetConfig sets the configuration of all open windows.
func (backend *Backend) SetConfig (config tomo.Config) {
	backend.assert()
	backend.config = config
	for _, window := range backend.windows {
		window.SetConfig(config)
	}
} 

func (backend *Backend) assert () {
	if backend == nil { panic("nil backend") }
}
