package x

import "image"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/ewmh"
import "github.com/jezek/xgbutil/icccm"
import "github.com/jezek/xgbutil/xevent"
import "github.com/jezek/xgbutil/xwindow"
import "github.com/jezek/xgbutil/xgraphics"
import "git.tebibyte.media/sashakoshka/tomo"

type Window struct {
	backend *Backend
	xWindow *xwindow.Window
	xCanvas *xgraphics.Image
	child   tomo.Element
	onClose func ()
	skipChildDrawCallback bool

	metrics struct {
		width  int
		height int
	}
}

func (backend *Backend) NewWindow (
	width, height int,
) (
	output tomo.Window,
	err error,
) {
	if backend == nil { panic("nil backend") }

	window := &Window { backend: backend }

	window.xWindow, err = xwindow.Generate(backend.connection)
	if err != nil { return }
	window.xWindow.Create (
		backend.connection.RootWin(),
		0, 0, width, height, 0)
	err = window.xWindow.Listen (
		xproto.EventMaskStructureNotify,
		xproto.EventMaskPointerMotion,
		xproto.EventMaskKeyPress,
		xproto.EventMaskKeyRelease,
		xproto.EventMaskButtonPress,
		xproto.EventMaskButtonRelease)
	if err != nil { return }

	window.xWindow.WMGracefulClose (func (xWindow *xwindow.Window) {
		window.Close()
	})
	
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
	
	window.metrics.width  = width
	window.metrics.height = height
	window.childMinimumSizeChangeCallback(8, 8)

	window.reallocateCanvas()

	backend.windows[window.xWindow.Id] = window
	output = window
	return
}

func (window *Window) Adopt (child tomo.Element) {
	if window.child != nil {
		child.SetParentHooks (tomo.ParentHooks { })
		if previousChild, ok := window.child.(tomo.Selectable); ok {
			if previousChild.Selected() {
				previousChild.HandleDeselection()
			}
		}
	}
	window.child = child
	if child != nil {
		child.SetParentHooks (tomo.ParentHooks {
			Draw: window.childDrawCallback,
			MinimumSizeChange: window.childMinimumSizeChangeCallback,
			SelectionRequest: window.childSelectionRequestCallback,
		})
		
		window.resizeChildToFit()
	}
	window.childMinimumSizeChangeCallback(child.MinimumSize())
}

func (window *Window) Child () (child tomo.Element) {
	child = window.child
	return
}

func (window *Window) SetTitle (title string) {
	ewmh.WmNameSet (
		window.backend.connection,
		window.xWindow.Id,
		title)
}

func (window *Window) SetIcon (sizes []image.Image) {
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

func (window *Window) Show () {
	if window.child == nil {
		window.xCanvas.For (func (x, y int) xgraphics.BGRA {
			return xgraphics.BGRA { }
		})

		window.pushRegion(window.xCanvas.Bounds())
	}
	
	window.xWindow.Map()
}

func (window *Window) Hide () {
	window.xWindow.Unmap()
}

func (window *Window) Close () {
	delete(window.backend.windows, window.xWindow.Id)
	if window.onClose != nil { window.onClose() }
	xevent.Detach(window.xWindow.X, window.xWindow.Id)
	window.xWindow.Destroy()
}

func (window *Window) OnClose (callback func ()) {
	window.onClose = callback
}

func (window *Window) reallocateCanvas () {
	if window.xCanvas != nil {
		window.xCanvas.Destroy()
	}
	window.xCanvas = xgraphics.New (
		window.backend.connection,
		image.Rect (
			0, 0,
			window.metrics.width,
			window.metrics.height))
	
	window.xCanvas.XSurfaceSet(window.xWindow.Id)
}

func (window *Window) redrawChildEntirely () {
	data, stride := window.child.Buffer()
	window.xCanvas.For (func (x, y int) (c xgraphics.BGRA) {
		rgba := data[x + y * stride]
		c.R, c.G, c.B, c.A = rgba.R, rgba.G, rgba.B, rgba.A
		return
	})
	
	window.pushRegion(window.xCanvas.Bounds())
}

func (window *Window) resizeChildToFit () {
	window.skipChildDrawCallback = true
	if child, ok := window.child.(tomo.Expanding); ok {
		minimumHeight := child.MinimumHeightFor(window.metrics.width)
		_, minimumWidth := child.MinimumSize()
		window.childMinimumSizeChangeCallback (
			minimumWidth, minimumHeight)
	} else {
		window.child.Resize (
			window.metrics.width,
			window.metrics.height)
			window.redrawChildEntirely()
	}
	window.skipChildDrawCallback = false
}

func (window *Window) childDrawCallback (region tomo.Canvas) {
	if window.skipChildDrawCallback { return }

	data, stride := region.Buffer()
	bounds := region.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x ++ {
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
		rgba := data[x + y * stride]
		window.xCanvas.SetBGRA (x, y, xgraphics.BGRA {
			R: rgba.R,
			G: rgba.G,
			B: rgba.B,
			A: rgba.A,
		})
	}}

	window.pushRegion(region.Bounds())
}

func (window *Window) childMinimumSizeChangeCallback (width, height int) {
	icccm.WmNormalHintsSet (
		window.backend.connection,
		window.xWindow.Id,
		&icccm.NormalHints {
			Flags:     icccm.SizeHintPMinSize,
			MinWidth:  uint(width),
			MinHeight: uint(height),
		})
	newWidth  := window.metrics.width
	newHeight := window.metrics.height
	if newWidth  < width  { newWidth  = width  }
	if newHeight < height { newHeight = height }
	if newWidth != window.metrics.width ||
		newHeight != window.metrics.height {
		window.xWindow.Resize(newWidth, newHeight)
	}
}

func (window *Window) childSelectionRequestCallback () (granted bool) {
	if child, ok := window.child.(tomo.Selectable); ok {
		child.HandleSelection(tomo.SelectionDirectionNeutral)
	}
	return true
}

func (window *Window) childSelectionMotionRequestCallback (
	direction tomo.SelectionDirection,
) (
	granted bool,
) {
	if child, ok := window.child.(tomo.Selectable); ok {
		if !child.HandleSelection(direction) {
			child.HandleDeselection()
		}
		return true
	}
	return true
}

func (window *Window) pushRegion (region image.Rectangle) {
	if window.xCanvas == nil { panic("whoopsie!!!!!!!!!!!!!!") }
	image, ok := window.xCanvas.SubImage(region).(*xgraphics.Image)
	if ok {
		image.XDraw()
		window.xCanvas.XPaint(window.xWindow.Id)
	}
}
