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
	return core.canvas.At(x, y)
}

func (core Core) Bounds () (bounds image.Rectangle) {
	return core.canvas.Bounds()
}

func (core Core) Set (x, y int, c color.Color) () {
	core.canvas.Set(x, y, c)
}

func (core Core) Buffer () (data []color.RGBA, stride int) {
	return core.canvas.Buffer()
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
	tomo.BasicCanvas
	core *Core
}

func (control CoreControl) HasImage () (empty bool) {
	return !control.Bounds().Empty()
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
	control.core.hooks.RunDraw(tomo.Cut(control, bounds))
}

func (control CoreControl) PushAll () {
	control.PushRegion(control.Bounds())
}

func (control *CoreControl) AllocateCanvas (width, height int) {
	core := control.core
	width, height, _ = control.ConstrainSize(width, height)
	core.canvas  = tomo.NewBasicCanvas(width, height)
	control.BasicCanvas = core.canvas
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
