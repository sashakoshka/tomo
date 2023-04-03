package core

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/shatter"

// Core is a struct that implements some core functionality common to most
// widgets. It is meant to be embedded directly into a struct.
type Core struct {
	canvas canvas.Canvas
	bounds image.Rectangle
	parent tomo.Parent
	outer  tomo.Element

	metrics struct {
		minimumWidth  int
		minimumHeight int
	}

	drawSizeChange func ()
	onDamage       func (region image.Rectangle)
}

// NewCore creates a new element core and its corresponding control given the
// element that it will be a part of. If outer is nil, this function will return
// nil.
func NewCore (
	outer tomo.Element,
	drawSizeChange func (),
) (
	core *Core,
	control CoreControl,
) {
	if outer == nil { return }
	core = &Core {
		outer:          outer,
		drawSizeChange: drawSizeChange,
	}
	control = CoreControl { core: core }
	return
}

// Bounds fulfills the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) Bounds () (bounds image.Rectangle) {
	if core.canvas == nil { return }
	return core.bounds
}

// MinimumSize fulfils the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) MinimumSize () (width, height int) {
	return core.metrics.minimumWidth, core.metrics.minimumHeight
}

// MinimumSize fulfils the tomo.Element interface. This should not need to be
// overridden, unless you want to detect when the element is parented or
// unparented.
func (core *Core) SetParent (parent tomo.Parent) {
	if parent != nil && core.parent != nil {
		panic("core.SetParent: element already has a parent")
	}

	core.parent = parent
}

// DrawTo fulfills the tomo.Element interface. This should not need to be
// overridden.
func (core *Core) DrawTo (
	canvas   canvas.Canvas,
	bounds   image.Rectangle,
	onDamage func (region image.Rectangle),
) {
	core.canvas   = canvas
	core.bounds   = bounds
	core.onDamage = onDamage
	if core.drawSizeChange != nil && core.canvas != nil {
		core.drawSizeChange()
	}
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

// Parent returns the element's parent.
func (control CoreControl) Parent () tomo.Parent {
	return control.core.parent
}

// DrawBackground fills the element's canvas with the parent's background
// pattern, if the parent supports it. If it is not supported, the fallback
// pattern will be used instead.
func (control CoreControl) DrawBackground (fallback artist.Pattern) {
	control.DrawBackgroundBounds(fallback, control.Bounds())
}

// DrawBackgroundBounds is like DrawBackground, but it takes in a bounding
// rectangle instead of using the element's bounds.
func (control CoreControl) DrawBackgroundBounds (
	fallback artist.Pattern,
	bounds image.Rectangle,
) {
	parent, ok := control.Parent().(tomo.BackgroundParent)
	if ok {
		parent.DrawBackground(bounds)
	} else if fallback != nil {
		fallback.Draw(canvas.Cut(control, bounds), control.Bounds())
	}
}

// DrawBackgroundBoundsShatter is like DrawBackgroundBounds, but uses the
// shattering algorithm to avoid drawing in areas specified by rocks.
func (control CoreControl) DrawBackgroundBoundsShatter (
	fallback artist.Pattern,
	bounds image.Rectangle,
	rocks ...image.Rectangle,
) {
	tiles := shatter.Shatter(bounds, rocks...)
	for _, tile := range tiles {
		control.DrawBackgroundBounds(fallback, tile)
	}
}

// Window returns the window containing the element.
func (control CoreControl) Window () tomo.Window {
	parent := control.Parent()
	if parent == nil {
		return nil
	} else {
		return parent.Window()
	}
}

// Outer returns the outer element given when the control was constructed.
func (control CoreControl) Outer () tomo.Element {
	return control.core.outer
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
			control.core.onDamage(region)
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
	if control.core.parent != nil {
		control.core.parent.NotifyMinimumSizeChange(control.core.outer)
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
