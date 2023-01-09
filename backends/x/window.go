package x

import "image"
import "image/color"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/ewmh"
import "github.com/jezek/xgbutil/icccm"
import "github.com/jezek/xgbutil/xevent"
import "github.com/jezek/xgbutil/xwindow"
import "github.com/jezek/xgbutil/xgraphics"
import "git.tebibyte.media/sashakoshka/tomo"

type Window struct {
	backend                   *Backend
	xWindow                   *xwindow.Window
	xCanvas                   *xgraphics.Image
	child                     tomo.Element
	onClose                   func ()
	drawCallback              func (region tomo.Image)
	minimumSizeChangeCallback func (width, height int)
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

func (window *Window) ColorModel () (model color.Model) {
	return color.RGBAModel
}

func (window *Window) At (x, y int) (pixel color.Color) {
	pixel = window.xCanvas.At(x, y)
	return
}

func (window *Window) RGBAAt (x, y int) (pixel color.RGBA) {
	sourcePixel := window.xCanvas.At(x, y).(xgraphics.BGRA)
	pixel = color.RGBA {
		R: sourcePixel.R,
		G: sourcePixel.G,
		B: sourcePixel.B,
		A: sourcePixel.A,
	}
	return
}

func (window *Window) Bounds () (bounds image.Rectangle) {
	bounds.Max = image.Point {
		X: window.metrics.width,
		Y: window.metrics.height,
	}
	return
}

func (window *Window) Handle (event tomo.Event) () {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		// we will receive a resize event from X later which will be
		// handled by our event handler callbacks.
		if resizeEvent.Width < window.MinimumWidth() {
			resizeEvent.Width = window.MinimumWidth()
		}
		if resizeEvent.Height < window.MinimumHeight() {
			resizeEvent.Height = window.MinimumHeight()
		}
		window.xWindow.Resize(resizeEvent.Width, resizeEvent.Height)
	default:
		if window.child != nil { window.child.Handle(event) }
	}
	return
}

func (window *Window) SetDrawCallback (draw func (region tomo.Image)) {
	window.drawCallback = draw
}

func (window *Window) SetMinimumSizeChangeCallback (
	notify func (width, height int),
) {
	window.minimumSizeChangeCallback = notify
}

func (window *Window) Selectable () (selectable bool) {
	if window.child != nil { selectable = window.child.Selectable() }
	return
}

func (window *Window) MinimumWidth () (minimum int) {
	if window.child != nil { minimum = window.child.MinimumWidth() }
	minimum = 8
	return
}

func (window *Window) MinimumHeight () (minimum int) {
	if window.child != nil { minimum = window.child.MinimumHeight() }
	minimum = 8
	return
}

func (window *Window) Adopt (child tomo.Element) {
	if window.child != nil {
		window.child.SetDrawCallback(nil)
		window.child.SetMinimumSizeChangeCallback(nil)
	}
	window.child = child
	if child != nil {
		child.SetDrawCallback(window.childDrawCallback)
		child.SetMinimumSizeChangeCallback (
			window.childMinimumSizeChangeCallback)
		window.resizeChildToFit()
	}
	window.childMinimumSizeChangeCallback (
		child.MinimumWidth(),
		child.MinimumHeight())
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
	window.xCanvas.For (func (x, y int) (c xgraphics.BGRA) {
		rgba := window.child.RGBAAt(x, y)
		c.R, c.G, c.B, c.A = rgba.R, rgba.G, rgba.B, rgba.A
		return
	})
	
	window.pushRegion(window.xCanvas.Bounds())
}

func (window *Window) resizeChildToFit () {
	window.skipChildDrawCallback = true
	window.child.Handle(tomo.EventResize {
		Width:  window.metrics.width,
		Height: window.metrics.height,
	})
	window.skipChildDrawCallback = false
	window.redrawChildEntirely()
}

func (window *Window) childDrawCallback (region tomo.Image) {
	if window.skipChildDrawCallback { return }

	bounds := region.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x ++ {
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
		rgba := region.RGBAAt(x, y)
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

func (window *Window) pushRegion (region image.Rectangle) {
	if window.xCanvas == nil { panic("whoopsie!!!!!!!!!!!!!!") }
	image, ok := window.xCanvas.SubImage(region).(*xgraphics.Image)
	if ok {
		image.XDraw()
		window.xCanvas.XPaint(window.xWindow.Id)
	}
}
