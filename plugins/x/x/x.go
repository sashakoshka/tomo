package x

import "git.tebibyte.media/sashakoshka/tomo"
import defaultTheme  "git.tebibyte.media/sashakoshka/tomo/default/theme"
import defaultConfig "git.tebibyte.media/sashakoshka/tomo/default/config"

import "github.com/jezek/xgbutil"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/xevent"
import "github.com/jezek/xgbutil/keybind"
import "github.com/jezek/xgbutil/mousebind"

type backend struct {
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
	backend := &backend {
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

func (backend *backend) Run () (err error) {
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

func (backend *backend) Stop () {
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

func (backend *backend) Do (callback func ()) {
	backend.assert()
	backend.doChannel <- callback
}

func (backend *backend) SetTheme (theme tomo.Theme) {
	backend.assert()
	if theme == nil {
		backend.theme = defaultTheme.Default { }
	} else {
		backend.theme = theme
	}
	for _, window := range backend.windows {
		window.handleThemeChange()
	}
}

func (backend *backend) SetConfig (config tomo.Config) {
	backend.assert()
	if config == nil {
		backend.config = defaultConfig.Default { }
	} else {
		backend.config = config
	}
	backend.config = config
	for _, window := range backend.windows {
		window.handleConfigChange()
	}
} 

func (backend *backend) assert () {
	if backend == nil { panic("nil backend") }
}
