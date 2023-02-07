package core

import "image"
import "image/color"
import "golang.org/x/image/font"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Core is a struct that implements some core functionality common to most
// widgets. It is meant to be embedded directly into a struct.
type Core struct {
	canvas canvas.Canvas

	metrics struct {
		minimumWidth  int
		minimumHeight int
	}

	config config.Config
	theme  theme.Theme
	c theme.Case

	handleSizeChange    func ()
	handleConfigChange  func ()
	handleThemeChange   func ()
	onMinimumSizeChange func ()
	onDamage func (region canvas.Canvas)
}

// NewCore creates a new element core and its corresponding control.
func NewCore (
	handleSizeChange   func (),
	handleConfigChange func (),
	handleThemeChange  func (),
	c theme.Case,
) (
	core *Core,
	control CoreControl,
) {
	core = &Core {
		handleSizeChange:   handleSizeChange,
		handleConfigChange: handleConfigChange,
		handleThemeChange:  handleThemeChange,
		c: c,
	}
	control = CoreControl { core: core }
	return
}

// ColorModel fulfills the draw.Image interface.
func (core *Core) ColorModel () (model color.Model) {
	return color.RGBAModel
}

// ColorModel fulfills the draw.Image interface.
func (core *Core) At (x, y int) (pixel color.Color) {
	if core.canvas == nil { return }
	return core.canvas.At(x, y)
}

// ColorModel fulfills the draw.Image interface.
func (core *Core) Bounds () (bounds image.Rectangle) {
	if core.canvas == nil { return }
	return core.canvas.Bounds()
}

// ColorModel fulfills the draw.Image interface.
func (core *Core) Set (x, y int, c color.Color) () {
	if core.canvas == nil { return }
	core.canvas.Set(x, y, c)
}

// Buffer fulfills the canvas.Canvas interface.
func (core *Core) Buffer () (data []color.RGBA, stride int) {
	if core.canvas == nil { return }
	return core.canvas.Buffer()
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
	if core.handleSizeChange != nil {
		core.handleSizeChange()
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

// SetConfig fulfills the elements.Configurable interface. This should not need
// to be overridden.
func (core *Core) SetConfig (config config.Config) {
	core.config = config
	if core.handleConfigChange != nil {
		core.handleConfigChange()
	}
}

// SetTheme fulfills the elements.Themeable interface. This should not need
// to be overridden.
func (core *Core) SetTheme (theme theme.Theme) {
	core.theme = theme
	if core.handleThemeChange != nil {
		core.handleThemeChange()
	}
}

// CoreControl is a struct that can exert control over a Core struct. It can be
// used as a canvas. It must not be directly embedded into an element, but
// instead kept as a private member. When a Core struct is created, a
// corresponding CoreControl struct is linked to it and returned alongside it.
type CoreControl struct {
	core *Core
}

// HasImage returns true if the core has an allocated image buffer, and false if
// it doesn't.
func (control CoreControl) HasImage () (has bool) {
	return control.core.canvas != nil && !control.core.canvas.Bounds().Empty()
}

// DamageRegion pushes the selected region of pixels to the parent element. This
// does not need to be called when responding to a resize event.
func (control CoreControl) DamageRegion (bounds image.Rectangle) {
	if control.core.onDamage != nil {
		control.core.onDamage(canvas.Cut(control.core, bounds))
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

// Config returns the current configuration.
func (control CoreControl) Config () (config.Config) {
	return control.core.config
}

// Theme returns the current theme.
func (control CoreControl) Theme () (theme.Theme) {
	return control.core.theme
}

// FontFace is like Theme.FontFace, but it automatically applies the correct
// case.
func (control CoreControl) FontFace (
	style theme.FontStyle,
	size  theme.FontSize,
) (
	face font.Face,
) {
	return control.core.theme.FontFace(style, size, control.core.c)
}

// Icon is like Theme.Icon, but it automatically applies the correct case.
func (control CoreControl) Icon (name string) (artist.Pattern) {
	return control.core.theme.Icon(name, control.core.c)
}

// Pattern is like Theme.Pattern, but it automatically applies the correct case.
func (control CoreControl) Pattern (
	id    theme.Pattern,
	state theme.PatternState,
) (
	pattern artist.Pattern,
) {
	return control.core.theme.Pattern(id, control.core.c, state)
}

// Inset is like Theme.Inset, but it automatically applies the correct case.
func (control CoreControl) Inset (id theme.Pattern) (inset theme.Inset) {
	return control.core.theme.Inset(id, control.core.c)
}

// Sink is like Theme.Sink, but it automatically applies the correct case.
func (control CoreControl) Sink (id theme.Pattern) (offset image.Point) {
	return control.core.theme.Sink(id, control.core.c)
}
