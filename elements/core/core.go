package core

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"

// Core is a struct that implements some core functionality common to most
// widgets. It is meant to be embedded directly into a struct.
type Core struct {
	canvas tomo.BasicCanvas
	parent tomo.Element

	metrics struct {
		minimumWidth  int
		minimumHeight int
	}

	onMinimumSizeChange func ()
	onDamage func (region tomo.Canvas)
}

// NewCore creates a new element core and its corresponding control.
func NewCore (parent tomo.Element) (core *Core, control CoreControl) {
	core    = &Core { parent: parent }
	control = CoreControl { core: core }
	return
}

// ColorModel fulfills the draw.Image interface.
func (core *Core) ColorModel () (model color.Model) {
	return color.RGBAModel
}

// ColorModel fulfills the draw.Image interface.
func (core *Core) At (x, y int) (pixel color.Color) {
	return core.canvas.At(x, y)
}

// ColorModel fulfills the draw.Image interface.
func (core *Core) Bounds () (bounds image.Rectangle) {
	return core.canvas.Bounds()
}

// ColorModel fulfills the draw.Image interface.
func (core *Core) Set (x, y int, c color.Color) () {
	core.canvas.Set(x, y, c)
}

// Buffer fulfills the tomo.Canvas interface.
func (core *Core) Buffer () (data []color.RGBA, stride int) {
	return core.canvas.Buffer()
}

// MinimumSize fulfils the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) MinimumSize () (width, height int) {
	return core.metrics.minimumWidth, core.metrics.minimumHeight
}

// OnDamage fulfils the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) OnDamage (callback func (region tomo.Canvas)) {
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
	tomo.BasicCanvas
	core *Core
}

// HasImage returns true if the core has an allocated image buffer, and false if
// it doesn't.
func (control CoreControl) HasImage () (has bool) {
	return !control.Bounds().Empty()
}

// DamageRegion pushes the selected region of pixels to the parent element. This
// does not need to be called when responding to a resize event.
func (control CoreControl) DamageRegion (bounds image.Rectangle) {
	if control.core.onDamage != nil {
		control.core.onDamage(tomo.Cut(control, bounds))
	}
}

// DamageAll pushes all pixels to the parent element. This does not need to be
// called when responding to a resize event.
func (control CoreControl) DamageAll () {
	control.DamageRegion(control.Bounds())
}

// AllocateCanvas resizes the canvas, constraining the width and height so that
// they are not less than the specified minimum width and height.
func (control *CoreControl) AllocateCanvas (width, height int) {
	control.core.canvas = tomo.NewBasicCanvas(width, height)
	control.BasicCanvas = control.core.canvas
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

	// if there is an image buffer, and the current size is less
	// than this new minimum size, send core.parent a resize event.
	if control.HasImage() {
		bounds := control.Bounds()
		imageWidth,
		imageHeight,
		constrained := control.ConstrainSize(bounds.Dx(), bounds.Dy())
		if constrained {
			core.parent.Resize(imageWidth, imageHeight)
		}
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
