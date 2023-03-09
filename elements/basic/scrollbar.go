package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type ScrollBar struct {
	*core.Core
	core core.CoreControl

	vertical bool
	enabled  bool
	dragging bool
	dragOffset int
	track image.Rectangle
	bar image.Rectangle

	contentBounds  image.Rectangle
	viewportBounds image.Rectangle
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onSlide   func ()
	onRelease func ()
}

func NewScrollBar (vertical bool) (element *ScrollBar) {
	element = &ScrollBar {
		vertical: vertical,
		enabled:  true,
	}
	if vertical {
		element.theme.Case = theme.C("basic", "scrollBarHorizontal")
	} else {
		element.theme.Case = theme.C("basic", "scrollBarVertical")
	}
	element.Core, element.core = core.NewCore(element.draw)
	element.updateMinimumSize()
	return
}

func (element *ScrollBar) HandleMouseDown (x, y int, button input.Button) {
	
}

func (element *ScrollBar) HandleMouseUp (x, y int, button input.Button) {
	
}

func (element *ScrollBar) HandleMouseMove (x, y int) {
	
}

func (element *ScrollBar) HandleMouseScroll (x, y int, deltaX, deltaY float64) {
	
}

// SetEnabled sets whether or not the scroll bar can be interacted with.
func (element *ScrollBar) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.redo()
}

// Enabled returns whether or not the element is enabled.
func (element *ScrollBar) Enabled () (enabled bool) {
	return element.enabled
}

// SetBounds sets the content and viewport bounds of the scroll bar.
func (element *ScrollBar) SetBounds (content, viewport image.Rectangle) {
	element.contentBounds  = content
	element.viewportBounds = viewport
	element.redo()
}

// SetTheme sets the element's theme.
func (element *ScrollBar) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *ScrollBar) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.redo()
}

func (element *ScrollBar) recalculate () {
	 if element.vertical {
	 	element.recalculateVertical()
	 } else {
	 	element.recalculateHorizontal()
	 }
}

func (element *ScrollBar) recalculateVertical () {
	bounds := element.Bounds()
	padding := element.theme.Padding(theme.PatternGutter)
	element.track = padding.Apply(bounds)

	contentBounds  := element.contentBounds
	viewportBounds := element.viewportBounds
	if element.Enabled() {
		element.bar.Min.X = element.track.Min.X
		element.bar.Max.X = element.track.Max.X

		scale := float64(element.track.Dy()) /
			float64(contentBounds.Dy())
		element.bar.Min.Y = int(float64(viewportBounds.Min.Y) * scale)
		element.bar.Max.Y = int(float64(viewportBounds.Max.Y) * scale)
		
		element.bar.Min.Y += element.track.Min.Y
		element.bar.Max.Y += element.track.Min.Y
	}

	// if the handle is out of bounds, don't display it
	if element.bar.Dy() >= element.track.Dy() {
		element.bar = image.Rectangle { }
	}
}

func (element *ScrollBar) recalculateHorizontal () {
	
}

func (element *ScrollBar) updateMinimumSize () {
	padding := element.theme.Padding(theme.PatternGutter)
	if element.vertical {
		element.core.SetMinimumSize (
			padding.Horizontal() + element.config.HandleWidth(),
			padding.Vertical()   + element.config.HandleWidth() * 2)
	} else {
		element.core.SetMinimumSize (
			padding.Horizontal() + element.config.HandleWidth() * 2,
			padding.Vertical()   + element.config.HandleWidth())
	}
}

func (element *ScrollBar) redo () {
	if element.core.HasImage () {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *ScrollBar) draw () {
	bounds := element.Bounds()
	state := theme.State {
		Disabled: !element.Enabled(),
		Pressed:  element.dragging,
	}
	artist.DrawBounds (
		element.core,
		element.theme.Pattern(theme.PatternGutter, state),
		bounds)
	artist.DrawBounds (
		element.core,
		element.theme.Pattern(theme.PatternHandle, state),
		element.bar)
}
