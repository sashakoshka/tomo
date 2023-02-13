package x

import "image"
import "github.com/jezek/xgb/xproto"
import "github.com/jezek/xgbutil/ewmh"
import "github.com/jezek/xgbutil/icccm"
import "github.com/jezek/xgbutil/xevent"
import "github.com/jezek/xgbutil/xwindow"
import "github.com/jezek/xgbutil/xgraphics"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/elements"

type Window struct {
	backend *Backend
	xWindow *xwindow.Window
	xCanvas *xgraphics.Image
	canvas  canvas.BasicCanvas
	child   elements.Element
	onClose func ()
	skipChildDrawCallback bool

	theme  theme.Theme
	config config.Config

	metrics struct {
		width  int
		height int
	}
}

func (backend *Backend) NewWindow (
	width, height int,
) (
	output elements.Window,
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
		xproto.EventMaskExposure,
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

	window.SetTheme(backend.theme)
	window.SetConfig(backend.config)
	
	window.metrics.width  = width
	window.metrics.height = height
	window.childMinimumSizeChangeCallback(8, 8)

	window.reallocateCanvas()

	backend.windows[window.xWindow.Id] = window
	output = window
	return
}

func (window *Window) Adopt (child elements.Element) {
	// disown previous child
	if window.child != nil {
		window.child.OnDamage(nil)
		window.child.OnMinimumSizeChange(nil)
	}
	if previousChild, ok := window.child.(elements.Flexible); ok {
		previousChild.OnFlexibleHeightChange(nil)
	}
	if previousChild, ok := window.child.(elements.Focusable); ok {
		previousChild.OnFocusRequest(nil)
		previousChild.OnFocusMotionRequest(nil)
		if previousChild.Focused() {
			previousChild.HandleUnfocus()
		}
	}
	
	// adopt new child
	window.child = child
	if newChild, ok := child.(elements.Themeable); ok {
		newChild.SetTheme(window.theme)
	}
	if newChild, ok := child.(elements.Configurable); ok {
		newChild.SetConfig(window.config)
	}
	if newChild, ok := child.(elements.Flexible); ok {
		newChild.OnFlexibleHeightChange(window.resizeChildToFit)
	}
	if newChild, ok := child.(elements.Focusable); ok {
		newChild.OnFocusRequest(window.childSelectionRequestCallback)
	}
	if child != nil {
		child.OnDamage(window.childDrawCallback)
		child.OnMinimumSizeChange (func () {
			window.childMinimumSizeChangeCallback (
				child.MinimumSize())
		})
		if !window.childMinimumSizeChangeCallback(child.MinimumSize()) {
			window.resizeChildToFit()
			window.redrawChildEntirely()
		}
	}
}

func (window *Window) Child () (child elements.Element) {
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

func (window *Window) SetTheme (theme theme.Theme) {
	window.theme = theme
	if child, ok := window.child.(elements.Themeable); ok {
		child.SetTheme(theme)
	}
}

func (window *Window) SetConfig (config config.Config) {
	window.config = config
	if child, ok := window.child.(elements.Configurable); ok {
		child.SetConfig(config)
	}
}

func (window *Window) reallocateCanvas () {
	window.canvas.Reallocate(window.metrics.width, window.metrics.height)

	previousWidth, previousHeight := 0, 0
	if window.xCanvas != nil {
		previousWidth  = window.xCanvas.Bounds().Dx()
		previousHeight = window.xCanvas.Bounds().Dy()
	}
	
	newWidth  := window.metrics.width
	newHeight := window.metrics.height
	larger    := newWidth > previousWidth || newHeight > previousHeight
	smaller   := newWidth < previousWidth / 2 || newHeight < previousHeight / 2
	if larger || smaller {
		if window.xCanvas != nil {
			window.xCanvas.Destroy()
		}
		window.xCanvas = xgraphics.New (
			window.backend.connection,
			image.Rect (
				0, 0,
				(newWidth  / 64) * 64 + 64,
				(newHeight / 64) * 64 + 64))
		window.xCanvas.CreatePixmap()
	}
	
}

func (window *Window) redrawChildEntirely () {
	window.pushRegion(window.paste(window.canvas))
	
}

func (window *Window) resizeChildToFit () {
	window.skipChildDrawCallback = true
	if child, ok := window.child.(elements.Flexible); ok {
		minimumHeight := child.FlexibleHeightFor(window.metrics.width)
		minimumWidth, _ := child.MinimumSize()
		
		icccm.WmNormalHintsSet (
			window.backend.connection,
			window.xWindow.Id,
			&icccm.NormalHints {
				Flags:     icccm.SizeHintPMinSize,
				MinWidth:  uint(minimumWidth),
				MinHeight: uint(minimumHeight),
			})
				
		if window.metrics.height >= minimumHeight &&
			window.metrics.width >= minimumWidth {
			window.child.DrawTo(window.canvas)
		}
	} else {
		window.child.DrawTo(window.canvas)
	}
	window.skipChildDrawCallback = false
}

func (window *Window) childDrawCallback (region canvas.Canvas) {
	if window.skipChildDrawCallback { return }
	window.pushRegion(window.paste(region))
}

func (window *Window) paste (canvas canvas.Canvas) (updatedRegion image.Rectangle) {
	data, stride := canvas.Buffer()
	bounds := canvas.Bounds().Intersect(window.xCanvas.Bounds())
	for x := bounds.Min.X; x < bounds.Max.X; x ++ {
	for y := bounds.Min.Y; y < bounds.Max.Y; y ++ {
		rgba := data[x + y * stride]
		index := x * 4 + y * window.xCanvas.Stride
		window.xCanvas.Pix[index + 0] = rgba.B
		window.xCanvas.Pix[index + 1] = rgba.G
		window.xCanvas.Pix[index + 2] = rgba.R
		window.xCanvas.Pix[index + 3] = rgba.A
	}}

	return bounds
}

func (window *Window) childMinimumSizeChangeCallback (width, height int) (resized bool) {
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
		return true
	}

	return false
}

func (window *Window) childSelectionRequestCallback () (granted bool) {
	if _, ok := window.child.(elements.Focusable); ok {
		return true
	}
	return false
}

func (window *Window) childSelectionMotionRequestCallback (
	direction input.KeynavDirection,
) (
	granted bool,
) {
	if child, ok := window.child.(elements.Focusable); ok {
		if !child.HandleFocus(direction) {
			child.HandleUnfocus()
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
		image.XExpPaint (
			window.xWindow.Id,
			image.Bounds().Min.X,
			image.Bounds().Min.Y)
	}
}
