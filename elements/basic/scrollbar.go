package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// ScrollBar is an element similar to Slider, but it has special behavior that
// makes it well suited for controlling the viewport position on one axis of a
// scrollable element. Instead of having a value from zero to one, it stores
// viewport and content boundaries. When the user drags the scroll bar handle,
// the scroll bar calls the OnScroll callback assigned to it with the position
// the user is trying to move the handle to. A program can check to see if this
// value is valid, move the viewport, and give the scroll bar the new viewport
// bounds (which will then cause it to move the handle).
//
// Typically, you wont't want to use a ScrollBar by itself. A ScrollContainer is
// better for most cases.
type ScrollBar struct {
	*core.Core
	core core.CoreControl

	vertical bool
	enabled  bool
	dragging bool
	dragOffset image.Point
	track image.Rectangle
	bar image.Rectangle

	contentBounds  image.Rectangle
	viewportBounds image.Rectangle
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onScroll func (viewport image.Point)
}

// NewScrollBar creates a new scroll bar. If vertical is set to true, the scroll
// bar will be vertical instead of horizontal.
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
	element.Core, element.core = core.NewCore(element.handleResize)
	element.updateMinimumSize()
	return
}

func (element *ScrollBar) handleResize () {
	if element.core.HasImage() {
		element.recalculate()
		element.draw()
	}
}

func (element *ScrollBar) HandleMouseDown (x, y int, button input.Button) {
	velocity := element.config.ScrollVelocity()
	point := image.Pt(x, y)

	if point.In(element.bar) {
		// the mouse is pressed down within the bar's handle
		element.dragging   = true
		element.drawAndPush()
		element.dragOffset =
			point.Sub(element.bar.Min).
			Add(element.Bounds().Min)
		element.dragTo(point)
	} else {
		// the mouse is pressed down within the bar's gutter
		switch button {
		case input.ButtonLeft:
			// start scrolling at this point, but set the offset to
			// the middle of the handle
			element.dragging = true
			element.dragOffset = element.fallbackDragOffset()
			element.dragTo(point)
			
		case input.ButtonMiddle:
			// page up/down on middle click
			viewport := 0
			if element.vertical {
				viewport = element.viewportBounds.Dy()
			} else {
				viewport = element.viewportBounds.Dx()
			}
			if element.isAfterHandle(point) {
				element.scrollBy(viewport)
			} else {
				element.scrollBy(-viewport)
			}
			
		case input.ButtonRight:
			// inch up/down on right click
			if element.isAfterHandle(point) {
				element.scrollBy(velocity)
			} else {
				element.scrollBy(-velocity)
			}
		}
	}
}

func (element *ScrollBar) HandleMouseUp (x, y int, button input.Button) {
	if element.dragging {
		element.dragging = false
		element.drawAndPush()
	}
}

func (element *ScrollBar) HandleMouseMove (x, y int) {
	if element.dragging {
		element.dragTo(image.Pt(x, y))
	}
}

func (element *ScrollBar) HandleMouseScroll (x, y int, deltaX, deltaY float64) {
	if element.vertical {
		element.scrollBy(int(deltaY))
	} else {
		element.scrollBy(int(deltaX))
	}
}

// SetEnabled sets whether or not the scroll bar can be interacted with.
func (element *ScrollBar) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.drawAndPush()
}

// Enabled returns whether or not the element is enabled.
func (element *ScrollBar) Enabled () (enabled bool) {
	return element.enabled
}

// SetBounds sets the content and viewport bounds of the scroll bar.
func (element *ScrollBar) SetBounds (content, viewport image.Rectangle) {
	element.contentBounds  = content
	element.viewportBounds = viewport
	element.recalculate()
	element.drawAndPush()
}

// OnScroll sets a function to be called when the user tries to move the scroll
// bar's handle. The callback is passed a point representing the new viewport
// position. For the scroll bar's position to visually update, the callback must
// check if the position is valid and call ScrollBar.SetBounds with the new
// viewport bounds.
func (element *ScrollBar) OnScroll (callback func (viewport image.Point)) {
	element.onScroll = callback
}

// SetTheme sets the element's theme.
func (element *ScrollBar) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.drawAndPush()
}

// SetConfig sets the element's configuration.
func (element *ScrollBar) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.updateMinimumSize()
	element.drawAndPush()
}

func (element *ScrollBar) isAfterHandle (point image.Point) bool {
	if element.vertical {
		return point.Y > element.bar.Min.Y
	} else {
		return point.X > element.bar.Min.X
	}
}

func (element *ScrollBar) fallbackDragOffset () image.Point {
	if element.vertical {
		return element.Bounds().Min.
			Add(image.Pt(0, element.bar.Dy() / 2))
	} else {
		return element.Bounds().Min.
			Add(image.Pt(element.bar.Dx() / 2, 0))
	}
}

func (element *ScrollBar) scrollBy (delta int) {
	deltaPoint := image.Point { }
	if element.vertical {
		deltaPoint.Y = delta
	} else {
		deltaPoint.X = delta
	}
	if element.onScroll != nil {
		element.onScroll(element.viewportBounds.Min.Add(deltaPoint))
	}
}

func (element *ScrollBar) dragTo (point image.Point) {
	point = point.Sub(element.dragOffset)
	var scrollX, scrollY float64

	if element.vertical {
		ratio :=
			float64(element.contentBounds.Dy()) /
			float64(element.track.Dy())
		scrollX = float64(element.viewportBounds.Min.X)
		scrollY = float64(point.Y) * ratio
	} else {
		ratio :=
			float64(element.contentBounds.Dx()) /
			float64(element.track.Dx())
		scrollX = float64(point.X) * ratio
		scrollY = float64(element.viewportBounds.Min.Y)
	}
	
	if element.onScroll != nil {
		element.onScroll(image.Pt(int(scrollX), int(scrollY)))
	}
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

		ratio :=
			float64(element.track.Dy()) /
			float64(contentBounds.Dy())
		element.bar.Min.Y = int(float64(viewportBounds.Min.Y) * ratio)
		element.bar.Max.Y = int(float64(viewportBounds.Max.Y) * ratio)
		
		element.bar.Min.Y += element.track.Min.Y
		element.bar.Max.Y += element.track.Min.Y
	}

	// if the handle is out of bounds, don't display it
	if element.bar.Dy() >= element.track.Dy() {
		element.bar = image.Rectangle { }
	}
}

func (element *ScrollBar) recalculateHorizontal () {
	bounds := element.Bounds()
	padding := element.theme.Padding(theme.PatternGutter)
	element.track = padding.Apply(bounds)

	contentBounds  := element.contentBounds
	viewportBounds := element.viewportBounds
	if element.Enabled() {
		element.bar.Min.Y = element.track.Min.Y
		element.bar.Max.Y = element.track.Max.Y

		ratio :=
			float64(element.track.Dx()) /
			float64(contentBounds.Dx())
		element.bar.Min.X = int(float64(viewportBounds.Min.X) * ratio)
		element.bar.Max.X = int(float64(viewportBounds.Max.X) * ratio)
		
		element.bar.Min.X += element.track.Min.X
		element.bar.Max.X += element.track.Min.X
	}

	// if the handle is out of bounds, don't display it
	if element.bar.Dx() >= element.track.Dx() {
		element.bar = image.Rectangle { }
	}
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

func (element *ScrollBar) drawAndPush () {
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
	element.theme.Pattern(theme.PatternGutter, state).Draw (
		element.core,
		bounds)
	element.theme.Pattern(theme.PatternHandle, state).Draw (
		element.core,
		element.bar)
}
