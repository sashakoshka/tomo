package basic

import "image"
import "image/color"
import "git.tebibyte.media/sashakoshka/tomo"

// Core is a struct that implements some core functionality common to most
// widgets. It is possible to embed this directly into a struct, but this is not
// reccomended as it exposes internal functionality.
type Core struct {
	*image.RGBA
	parent tomo.Element
	
	drawCallback func (region tomo.Image)
	minimumSizeChangeCallback func (width, height int)

	metrics struct {
		minimumWidth  int
		minimumHeight int
	}
}

// Core creates a new element core.
func NewCore (parent tomo.Element) (core Core) {
	core = Core { parent: parent }
	return
}

func (core Core) ColorModel () (model color.Model) {
	return color.RGBAModel
}

func (core Core) At (x, y int) (pixel color.Color) {
	if core.RGBA == nil { return color.RGBA { } }
	pixel = core.RGBA.At(x, y)
	return
}

func (core Core) RGBAAt (x, y int) (pixel color.RGBA) {
	if core.RGBA == nil { return color.RGBA { } }
	pixel = core.RGBA.RGBAAt(x, y)
	return
}

func (core Core) Bounds () (bounds image.Rectangle) {
	if core.RGBA != nil { bounds = core.RGBA.Bounds() }
	return
}

func (core *Core) SetDrawCallback (draw func (region tomo.Image)) {
	core.drawCallback = draw
}

func (core *Core) SetMinimumSizeChangeCallback (
	notify func (width, height int),
) {
	core.minimumSizeChangeCallback = notify
}

func (core Core) HasImage () (has bool) {
	has = core.RGBA != nil
	return
}

func (core Core) PushRegion (bounds image.Rectangle) {
	if core.drawCallback != nil {
		core.drawCallback(core.SubImage(bounds).
			(*image.RGBA))
	}
}

func (core Core) PushAll () {
	core.PushRegion(core.Bounds())
}

func (core *Core) AllocateCanvas (width, height int) {
	width, height, _ = core.ConstrainSize(width, height)
	core.RGBA = image.NewRGBA(image.Rect (0, 0, width, height))
}

func (core Core) MinimumWidth () (minimum int) {
	minimum = core.metrics.minimumWidth
	return
}

func (core Core) MinimumHeight () (minimum int) {
	minimum = core.metrics.minimumHeight
	return
}

func (core *Core) SetMinimumSize (width, height int) {
	if width != core.metrics.minimumWidth ||
		height != core.metrics.minimumHeight {

		core.metrics.minimumWidth  = width
		core.metrics.minimumHeight = height
		
		if core.minimumSizeChangeCallback != nil {
			core.minimumSizeChangeCallback(width, height)
		}

		// if there is an image buffer, and the current size is less
		// than this new minimum size, send core.parent a resize event.
		if core.HasImage() {
			bounds := core.Bounds()
			imageWidth,
			imageHeight,
			constrained := core.ConstrainSize (
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
}

func (core Core) ConstrainSize (
	inWidth, inHeight int,
) (
	outWidth, outHeight int,
	constrained bool,
) {
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
