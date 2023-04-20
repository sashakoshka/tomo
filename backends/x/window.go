package x

import "image"
import "errors"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/ewmh"
import "github.com/jezek/xgbutil/icccm"
import "github.com/jezek/xgbutil/xprop"
import "github.com/jezek/xgbutil/xevent"
import "github.com/jezek/xgbutil/xwindow"
import "github.com/jezek/xgbutil/keybind"
import "github.com/jezek/xgbutil/mousebind"
import "github.com/jezek/xgbutil/xgraphics"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/data"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

type mainWindow struct { *window }
type menuWindow struct { *window }
type window struct {
	system
	
	backend *Backend
	xWindow *xwindow.Window
	xCanvas *xgraphics.Image

	title, application string

	modalParent *window
	hasModal    bool
	shy         bool

	selectionRequest *selectionRequest
	selectionClaim   *selectionClaim

	metrics struct {
		bounds image.Rectangle
	}

	onClose func ()
}

func (backend *Backend) NewWindow (
	bounds image.Rectangle,
) (
	output tomo.MainWindow,
	err error,
) {
	if backend == nil { panic("nil backend") }
	window, err := backend.newWindow(bounds, false)
	
	output = mainWindow { window }
	return output, err
}

func (backend *Backend) newWindow (
	bounds   image.Rectangle,
	override bool,
) (
	output *window,
	err error,
) {
	if bounds.Dx() == 0 { bounds.Max.X = bounds.Min.X + 8 }
	if bounds.Dy() == 0 { bounds.Max.Y = bounds.Min.Y + 8 }
	
	window := &window { backend: backend }

	window.system.initialize()
	window.system.pushFunc = window.pasteAndPush
	window.theme.Case = tomo.C("tomo", "window")

	window.xWindow, err = xwindow.Generate(backend.connection)
	if err != nil { return }

	if override {
		err = window.xWindow.CreateChecked (
			backend.connection.RootWin(),
			bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy(),
			xproto.CwOverrideRedirect, 1)
	} else {
		err = window.xWindow.CreateChecked (
			backend.connection.RootWin(),
			bounds.Min.X, bounds.Min.Y, bounds.Dx(), bounds.Dy(), 0)
	}
	if err != nil { return }
	
	err = window.xWindow.Listen (
		xproto.EventMaskExposure,
		xproto.EventMaskStructureNotify,
		xproto.EventMaskPropertyChange,
		xproto.EventMaskPointerMotion,
		xproto.EventMaskKeyPress,
		xproto.EventMaskKeyRelease,
		xproto.EventMaskButtonPress,
		xproto.EventMaskButtonRelease)
	if err != nil { return }

	window.xWindow.WMGracefulClose (func (xWindow *xwindow.Window) {
		window.Close()
	})

	xevent.ExposeFun(window.handleExpose).
		Connect(backend.connection, window.xWindow.Id)
	xevent.ConfigureNotifyFun(window.handleConfigureNotify).
		Connect(backend.connection, window.xWindow.Id)
	xevent.KeyPressFun(window.handleKeyPress).
		Connect(backend.connection, window.xWindow.Id)
	xevent.KeyReleaseFun(window.handleKeyRelease).
		Connect(backend.connection, window.xWindow.Id)
	xevent.ButtonPressFun(window.handleButtonPress).
		Connect(backend.connection, window.xWindow.Id)
	xevent.ButtonReleaseFun(window.handleButtonRelease).
		Connect(backend.connection, window.xWindow.Id)
	xevent.MotionNotifyFun(window.handleMotionNotify).
		Connect(backend.connection, window.xWindow.Id)
	xevent.SelectionNotifyFun(window.handleSelectionNotify).
		Connect(backend.connection, window.xWindow.Id)
	xevent.PropertyNotifyFun(window.handlePropertyNotify).
		Connect(backend.connection, window.xWindow.Id)
	xevent.SelectionClearFun(window.handleSelectionClear).
		Connect(backend.connection, window.xWindow.Id)
	xevent.SelectionRequestFun(window.handleSelectionRequest).
		Connect(backend.connection, window.xWindow.Id)

	window.SetTheme(backend.theme)
	window.SetConfig(backend.config)
	
	window.metrics.bounds = bounds
	window.setMinimumSize(8, 8)

	window.reallocateCanvas()

	backend.windows[window.xWindow.Id] = window

	output = window
	return
}

func (window *window) Window () tomo.Window {
	return window
}

func (window *window) Adopt (child tomo.Element) {
	// disown previous child
	if window.child != nil {
		window.child.unlink()
		window.child = nil
	}

	// adopt new child
	if child != nil {
		childEntity, ok := child.Entity().(*entity)
		if ok && childEntity != nil {
			window.child = childEntity
			childEntity.setWindow(window)
			window.setMinimumSize (
				childEntity.minWidth,
				childEntity.minHeight)
			window.resizeChildToFit()
		}
	}
}

func (window *window) SetTitle (title string) {
	window.title = title
	ewmh.WmNameSet (
		window.backend.connection,
		window.xWindow.Id,
		title)
	icccm.WmNameSet (
		window.backend.connection,
		window.xWindow.Id,
		title)
	icccm.WmIconNameSet (
		window.backend.connection,
		window.xWindow.Id,
		title)
}

func (window *window) SetApplicationName (name string) {
	window.application = name
	icccm.WmClassSet (
		window.backend.connection,
		window.xWindow.Id,
		&icccm.WmClass {
			Instance: name,
			Class:    name,
		})
}

func (window *window) SetIcon (sizes []image.Image) {
	wmIcons := []ewmh.WmIcon { }
	
	for _, icon := range sizes {
		width  := icon.Bounds().Max.X
		height := icon.Bounds().Max.Y
		wmIcon := ewmh.WmIcon {
			Width:  uint(width),
			Height: uint(height),
			Data:   make ([]uint, width * height),
		}

		// manually convert image data beacuse of course we have to do
		// this
		index := 0
		for y := 0; y < height; y ++ {
		for x := 0; x < width;  x ++ {
			r, g, b, a := icon.At(x, y).RGBA()
			r >>= 8
			g >>= 8
			b >>= 8
			a >>= 8
			wmIcon.Data[index] =
				(uint(a) << 24) |
				(uint(r) << 16) |
				(uint(g) << 8)  |
				(uint(b) << 0)
			index ++
		}}
		
		wmIcons = append(wmIcons, wmIcon)
	}

	ewmh.WmIconSet (
		window.backend.connection,
		window.xWindow.Id,
		wmIcons)
}

func (window *window) NewModal (bounds image.Rectangle) (tomo.Window, error) {
	modal, err := window.backend.newWindow (
		bounds.Add(window.metrics.bounds.Min), false)
	icccm.WmTransientForSet (
		window.backend.connection,
		modal.xWindow.Id,
		window.xWindow.Id)
	ewmh.WmStateSet (
		window.backend.connection,
		modal.xWindow.Id,
		[]string { "_NET_WM_STATE_MODAL" })
	modal.modalParent = window
	window.hasModal   = true
	modal.inheritProperties(window)
	return modal, err
}

func (window *window) NewMenu (bounds image.Rectangle) (tomo.MenuWindow, error) {
	menu, err := window.backend.newWindow (
		bounds.Add(window.metrics.bounds.Min), true)
	menu.shy = true
	icccm.WmTransientForSet (
		window.backend.connection,
		menu.xWindow.Id,
		window.xWindow.Id)
	menu.setType("POPUP_MENU")
	menu.inheritProperties(window)
	return menuWindow { window: menu }, err
}

func (window mainWindow) NewPanel (bounds image.Rectangle) (tomo.Window, error) {
	panel, err := window.backend.newWindow (
		bounds.Add(window.metrics.bounds.Min), false)
	if err != nil { return nil, err }
	panel.setClientLeader(window.window)
	window.setClientLeader(window.window)
	icccm.WmTransientForSet (
		window.backend.connection,
		panel.xWindow.Id,
		window.xWindow.Id)
	panel.setType("UTILITY")
	panel.inheritProperties(window.window)
	return panel, err
}

func (window menuWindow) Pin () {
	// TODO take off override redirect
	// TODO turn off shy
	// TODO set window type to MENU
	// TODO iungrab keyboard and mouse
}

func (window *window) Show () {
	if window.child == nil {
		window.xCanvas.For (func (x, y int) xgraphics.BGRA {
			return xgraphics.BGRA { }
		})

		window.pushRegion(window.xCanvas.Bounds())
	}

	window.xWindow.Map()
	if window.shy { window.grabInput() }
}

func (window *window) Hide () {
	window.xWindow.Unmap()
	if window.shy { window.ungrabInput() }
}

func (window *window) Copy (data data.Data) {
	selectionAtom, err := xprop.Atm(window.backend.connection, clipboardName)
	if err != nil { return }
	window.selectionClaim = window.claimSelection(selectionAtom, data)
}

func (window *window) Paste (callback func (data.Data, error), accept ...data.Mime) {
	// Follow:
	// https://tronche.com/gui/x/icccm/sec-2.html#s-2.4
	die := func (err error) { callback(nil, err) }
	if window.selectionRequest != nil {
		// TODO: add the request to a queue and take care of it when the
		// current selection has completed
		die(errors.New("there is already a selection request"))
		return
	}

	propertyName := "TOMO_SELECTION"
	selectionAtom, err := xprop.Atm(window.backend.connection, clipboardName)
	if err != nil { die(err); return }
	propertyAtom, err := xprop.Atm(window.backend.connection, propertyName)
	if err != nil { die(err); return }

	window.selectionRequest = window.newSelectionRequest (
		selectionAtom, propertyAtom, callback, accept...)
	if !window.selectionRequest.open() { window.selectionRequest = nil }
	return
}

func (window *window) Close () {
	if window.onClose != nil { window.onClose() }
	if window.modalParent != nil {
		// we are a modal dialog, so unlock the parent
		window.modalParent.hasModal = false
	}
	window.Hide()
	window.Adopt(nil)
	delete(window.backend.windows, window.xWindow.Id)
	window.xWindow.Destroy()
}

func (window *window) OnClose (callback func ()) {
	window.onClose = callback
}

func (window *window) grabInput () {
	keybind.GrabKeyboard(window.backend.connection, window.xWindow.Id)
	mousebind.GrabPointer (
		window.backend.connection,
		window.xWindow.Id,
		window.backend.connection.RootWin(), 0)
}

func (window *window) ungrabInput () {
	keybind.UngrabKeyboard(window.backend.connection)
	mousebind.UngrabPointer(window.backend.connection)
}

func (window *window) inheritProperties (parent *window) {
	window.SetApplicationName(parent.application)
}

func (window *window) setType (ty string) error {
	return ewmh.WmWindowTypeSet (
		window.backend.connection,
		window.xWindow.Id,
		[]string { "_NET_WM_WINDOW_TYPE_" + ty })
}

func (window *window) setClientLeader (leader *window) error {
	hints, _ := icccm.WmHintsGet(window.backend.connection, window.xWindow.Id)
	if hints == nil {
		hints = &icccm.Hints { }
	}
	hints.Flags |= icccm.HintWindowGroup
	hints.WindowGroup = leader.xWindow.Id
	return icccm.WmHintsSet (
		window.backend.connection,
		window.xWindow.Id,
		hints)
}

func (window *window) reallocateCanvas () {
	window.canvas.Reallocate (
		window.metrics.bounds.Dx(),
		window.metrics.bounds.Dy())

	previousWidth, previousHeight := 0, 0
	if window.xCanvas != nil {
		previousWidth  = window.xCanvas.Bounds().Dx()
		previousHeight = window.xCanvas.Bounds().Dy()
	}
	
	newWidth  := window.metrics.bounds.Dx()
	newHeight := window.metrics.bounds.Dy()
	larger    := newWidth > previousWidth || newHeight > previousHeight
	smaller   := newWidth < previousWidth / 2 || newHeight < previousHeight / 2

	allocStep := 128
	
	if larger || smaller {
		if window.xCanvas != nil {
			window.xCanvas.Destroy()
		}
		window.xCanvas = xgraphics.New (
			window.backend.connection,
			image.Rect (
				0, 0,
				(newWidth  / allocStep + 1) * allocStep,
				(newHeight / allocStep + 1) * allocStep))
		window.xCanvas.CreatePixmap()
	}
	
}

func (window *window) pasteAndPush (region image.Rectangle) {
	window.paste(region)
	window.pushRegion(region)
}

func (window *window) paste (region image.Rectangle) {
	canvas := canvas.Cut(window.canvas, region)
	data, stride := canvas.Buffer()
	bounds := canvas.Bounds().Intersect(window.xCanvas.Bounds())

	dstStride := window.xCanvas.Stride
	dstData   := window.xCanvas.Pix
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
		srcYComponent := y * stride
		dstYComponent := y * dstStride
		for x := bounds.Min.X; x < bounds.Max.X; x ++ {
			rgba := data[srcYComponent + x]
			index := dstYComponent + x * 4
			dstData[index + 0] = rgba.B
			dstData[index + 1] = rgba.G
			dstData[index + 2] = rgba.R
			dstData[index + 3] = rgba.A
		}
	}
}

func (window *window) pushRegion (region image.Rectangle) {
	if window.xCanvas == nil { panic("whoopsie!!!!!!!!!!!!!!") }
	image, ok := window.xCanvas.SubImage(region).(*xgraphics.Image)
	if ok {
		image.XDraw()
		image.XExpPaint (
			window.xWindow.Id,
			image.Bounds().Min.X,
			image.Bounds().Min.Y)
	}
}

func (window *window) setMinimumSize (width, height int) {
	if width  < 8 { width  = 8 }
	if height < 8 { height = 8 }
	icccm.WmNormalHintsSet (
		window.backend.connection,
		window.xWindow.Id,
		&icccm.NormalHints {
			Flags:     icccm.SizeHintPMinSize,
			MinWidth:  uint(width),
			MinHeight: uint(height),
		})
	newWidth  := window.metrics.bounds.Dx()
	newHeight := window.metrics.bounds.Dy()
	if newWidth  < width  { newWidth  = width  }
	if newHeight < height { newHeight = height }
	if newWidth != window.metrics.bounds.Dx() ||
		newHeight != window.metrics.bounds.Dy() {
		window.xWindow.Resize(newWidth, newHeight)
	}
}
