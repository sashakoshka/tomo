package core

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"

// Core is a struct that implements some core functionality common to most
// widgets. It is meant to be embedded directly into a struct.
type Core struct {
	canvas *image.RGBA
	parent tomo.Element

	metrics struct {
		minimumWidth  int
		minimumHeight int
	}

	selectable bool
	selected   bool
	hooks tomo.ParentHooks
}

// NewCore creates a new element core and its corresponding control.
func NewCore (parent tomo.Element) (core *Core, control CoreControl) {
	core    = &Core { parent: parent }
	control = CoreControl { core: core }
	return
}

func (core Core) ColorModel () (model color.Model) {
	return color.RGBAModel
}

func (core Core) At (x, y int) (pixel color.Color) {
	if core.canvas == nil { return color.RGBA { } }
	pixel = core.canvas.At(x, y)
	return
}

func (core Core) RGBAAt (x, y int) (pixel color.RGBA) {
	if core.canvas == nil { return color.RGBA { } }
	pixel = core.canvas.RGBAAt(x, y)
	return
}

func (core Core) Bounds () (bounds image.Rectangle) {
	if core.canvas != nil { bounds = core.canvas.Bounds() }
	return
}

func (core Core) Selectable () (selectable bool) {
	return core.selectable
}

func (core Core) Selected () (selected bool) {
	return core.selected
}

func (core Core) AdvanceSelection (direction int) (ok bool) {
	return
}

func (core *Core) SetParentHooks (hooks tomo.ParentHooks) {
	core.hooks = hooks
}

func (core Core) MinimumSize () (width, height int) {
	return core.metrics.minimumWidth, core.metrics.minimumHeight
}

// CoreControl is a struct that can exert control over a control struct. It can
// be used as a canvas. It must not be directly embedded into an element, but
// instead kept as a private member.
type CoreControl struct {
	*image.RGBA
	core *Core
}

func (control CoreControl) HasImage () (has bool) {
	has = control.RGBA != nil
	return
}

func (control CoreControl) Select () (granted bool) {
	return control.core.hooks.RunSelectionRequest()
}

func (control CoreControl) SetSelected (selected bool) {
	if !control.core.selectable { return }
	control.core.selected = selected
}

func (control CoreControl) SetSelectable (selectable bool) {
	if control.core.selectable == selectable { return }
	control.core.selectable = selectable
	if !selectable { control.core.selected = false }
	control.core.hooks.RunSelectabilityChange(selectable)
}

func (control CoreControl) PushRegion (bounds image.Rectangle) {
	control.core.hooks.RunDraw(control.SubImage(bounds).(*image.RGBA))
}

func (control CoreControl) PushAll () {
	control.PushRegion(control.Bounds())
}

func (control *CoreControl) AllocateCanvas (width, height int) {
	core := control.core
	width, height, _ = control.ConstrainSize(width, height)
	core.canvas  = image.NewRGBA(image.Rect (0, 0, width, height))
	control.RGBA = core.canvas
}

func (control CoreControl) SetMinimumSize (width, height int) {
	core := control.core
	if width == core.metrics.minimumWidth &&
		height == core.metrics.minimumHeight {
		return
	}

	core.metrics.minimumWidth  = width
	core.metrics.minimumHeight = height
	core.hooks.RunMinimumSizeChange(width, height)

	// if there is an image buffer, and the current size is less
	// than this new minimum size, send core.parent a resize event.
	if control.HasImage() {
		bounds := control.Bounds()
		imageWidth,
		imageHeight,
		constrained := control.ConstrainSize (
			bounds.Dx(),
			bounds.Dy())
		if constrained {
			core.parent.Handle (tomo.EventResize {
				Width:  imageWidth,
				Height: imageHeight,
			})
		}
	}
}

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
