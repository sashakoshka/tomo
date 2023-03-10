package core

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

// Core is a struct that implements some core functionality common to most
// widgets. It is meant to be embedded directly into a struct.
type Core struct {
	canvas canvas.Canvas

	metrics struct {
		minimumWidth  int
		minimumHeight int
	}

	drawSizeChange      func ()
	onMinimumSizeChange func ()
	onDamage func (region canvas.Canvas)
}

// NewCore creates a new element core and its corresponding control.
func NewCore (
	drawSizeChange func (),
) (
	core *Core,
	control CoreControl,
) {
	core = &Core {
		drawSizeChange: drawSizeChange,
	}
	control = CoreControl { core: core }
	return
}

// Bounds fulfills the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) Bounds () (bounds image.Rectangle) {
	if core.canvas == nil { return }
	return core.canvas.Bounds()
}

// MinimumSize fulfils the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) MinimumSize () (width, height int) {
	return core.metrics.minimumWidth, core.metrics.minimumHeight
}

// DrawTo fulfills the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) DrawTo (canvas canvas.Canvas) {
	core.canvas = canvas
	if core.drawSizeChange != nil && core.canvas != nil {
		core.drawSizeChange()
	}
}

// OnDamage fulfils the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) OnDamage (callback func (region canvas.Canvas)) {
	core.onDamage = callback
}

// OnMinimumSizeChange fulfils the tomo.Element interface. This should not need
// to be overridden.
func (core *Core) OnMinimumSizeChange (callback func ()) {
	core.onMinimumSizeChange = callback
}

// CoreControl is a struct that can exert control over a Core struct. It can be
// used as a canvas. It must not be directly embedded into an element, but
// instead kept as a private member. When a Core struct is created, a
// corresponding CoreControl struct is linked to it and returned alongside it.
type CoreControl struct {
	core *Core
}

// ColorModel fulfills the draw.Image interface.
func (control CoreControl) ColorModel () (model color.Model) {
	return color.RGBAModel
}

// At fulfills the draw.Image interface.
func (control CoreControl) At (x, y int) (pixel color.Color) {
	if control.core.canvas == nil { return }
	return control.core.canvas.At(x, y)
}

// Bounds fulfills the draw.Image interface.
func (control CoreControl) Bounds () (bounds image.Rectangle) {
	if control.core.canvas == nil { return }
	return control.core.canvas.Bounds()
}

// Set fulfills the draw.Image interface.
func (control CoreControl) Set (x, y int, c color.Color) () {
	if control.core.canvas == nil { return }
	control.core.canvas.Set(x, y, c)
}

// Buffer fulfills the canvas.Canvas interface.
func (control CoreControl) Buffer () (data []color.RGBA, stride int) {
	if control.core.canvas == nil { return }
	return control.core.canvas.Buffer()
}

// HasImage returns true if the core has an allocated image buffer, and false if
// it doesn't.
func (control CoreControl) HasImage () (has bool) {
	return control.core.canvas != nil && !control.core.canvas.Bounds().Empty()
}

// DamageRegion pushes the selected region of pixels to the parent element. This
// does not need to be called when responding to a resize event.
func (control CoreControl) DamageRegion (regions ...image.Rectangle) {
	if control.core.canvas == nil { return }
	if control.core.onDamage != nil {
		for _, region := range regions {
			control.core.onDamage (
				canvas.Cut(control.core.canvas, region))
		}
	}
}

// DamageAll pushes all pixels to the parent element. This does not need to be
// called when redrawing in response to a change in size.
func (control CoreControl) DamageAll () {
	control.DamageRegion(control.core.Bounds())
}

// SetMinimumSize sets the minimum size of this element, notifying the parent
// element in the process.
func (control CoreControl) SetMinimumSize (width, height int) {
	core := control.core
	if width == core.metrics.minimumWidth &&
		height == core.metrics.minimumHeight {
		return
	}

	core.metrics.minimumWidth  = width
	core.metrics.minimumHeight = height
	if control.core.onMinimumSizeChange != nil {
		control.core.onMinimumSizeChange()
	}
}

// ConstrainSize contstrains the specified width and height to the minimum width
// and height, and returns wether or not anything ended up being constrained.
func (control CoreControl) ConstrainSize (
	inWidth, inHeight int,
) (
	outWidth, outHeight int,
	constrained bool,
) {
	core := control.core
	outWidth  = inWidth
	outHeight = inHeight
	if outWidth < core.metrics.minimumWidth {
		outWidth = core.metrics.minimumWidth
		constrained = true
	}
	if outHeight < core.metrics.minimumHeight {
		outHeight = core.metrics.minimumHeight
		constrained = true
	}
	return
}
