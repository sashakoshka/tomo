package elements

import "image"
import "tomo"
import "tomo/input"
import "art"

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
	entity tomo.Entity

	c tomo.Case

	vertical bool
	enabled  bool
	dragging bool
	dragOffset image.Point
	track image.Rectangle
	bar image.Rectangle

	contentBounds  image.Rectangle
	viewportBounds image.Rectangle
	
	onScroll func (viewport image.Point)
}

// NewVScrollBar creates a new vertical scroll bar.
func NewVScrollBar () (element *ScrollBar) {
	element = &ScrollBar {
		vertical: true,
		enabled:  true,
	}
	element.c = tomo.C("tomo", "scrollBarVertical")
	element.entity = tomo.GetBackend().NewEntity(element)
	element.updateMinimumSize()
	return
}

// NewHScrollBar creates a new horizontal scroll bar.
func NewHScrollBar () (element *ScrollBar) {
	element = &ScrollBar {
		enabled: true,
	}
	element.c = tomo.C("tomo", "scrollBarHorizontal")
	element.entity = tomo.GetBackend().NewEntity(element)
	element.updateMinimumSize()
	return
}

// Entity returns this element's entity.
func (element *ScrollBar) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *ScrollBar) Draw (destination art.Canvas) {
	element.recalculate()

	bounds := element.entity.Bounds()
	state := tomo.State {
		Disabled: !element.Enabled(),
		Pressed:  element.dragging,
	}
	element.entity.Theme().Pattern(tomo.PatternGutter, state, element.c).Draw (
		destination,
		bounds)
	element.entity.Theme().Pattern(tomo.PatternHandle, state, element.c).Draw (
		destination,
		element.bar)
}

func (element *ScrollBar) HandleMouseDown  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	velocity := element.entity.Config().ScrollVelocity()

	if position.In(element.bar) {
		// the mouse is pressed down within the bar's handle
		element.dragging   = true
		element.entity.Invalidate()
		element.dragOffset =
			position.Sub(element.bar.Min).
			Add(element.entity.Bounds().Min)
		element.dragTo(position)
	} else {
		// the mouse is pressed down within the bar's gutter
		switch button {
		case input.ButtonLeft:
			// start scrolling at this point, but set the offset to
			// the middle of the handle
			element.dragging = true
			element.dragOffset = element.fallbackDragOffset()
			element.dragTo(position)
			
		case input.ButtonMiddle:
			// page up/down on middle click
			viewport := 0
			if element.vertical {
				viewport = element.viewportBounds.Dy()
			} else {
				viewport = element.viewportBounds.Dx()
			}
			if element.isAfterHandle(position) {
				element.scrollBy(viewport)
			} else {
				element.scrollBy(-viewport)
			}
			
		case input.ButtonRight:
			// inch up/down on right click
			if element.isAfterHandle(position) {
				element.scrollBy(velocity)
			} else {
				element.scrollBy(-velocity)
			}
		}
	}
}

func (element *ScrollBar) HandleMouseUp  (
	position image.Point,
	button input.Button,
	modifiers input.Modifiers,
) {
	if element.dragging {
		element.dragging = false
		element.entity.Invalidate()
	}
}

func (element *ScrollBar) HandleMotion (position image.Point) {
	if element.dragging {
		element.dragTo(position)
	}
}

func (element *ScrollBar) HandleScroll (
	position image.Point,
	deltaX, deltaY float64,
	modifiers input.Modifiers,
) {
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
	element.entity.Invalidate()
}

// Enabled returns whether or not the element is enabled.
func (element *ScrollBar) Enabled () (enabled bool) {
	return element.enabled
}

// SetBounds sets the content and viewport bounds of the scroll bar.
func (element *ScrollBar) SetBounds (content, viewport image.Rectangle) {
	element.contentBounds  = content
	element.viewportBounds = viewport
	element.entity.Invalidate()
}

// OnScroll sets a function to be called when the user tries to move the scroll
// bar's handle. The callback is passed a point representing the new viewport
// position. For the scroll bar's position to visually update, the callback must
// check if the position is valid and call ScrollBar.SetBounds with the new
// viewport bounds.
func (element *ScrollBar) OnScroll (callback func (viewport image.Point)) {
	element.onScroll = callback
}

func (element *ScrollBar) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
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
		return element.entity.Bounds().Min.
			Add(image.Pt(0, element.bar.Dy() / 2))
	} else {
		return element.entity.Bounds().Min.
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
	bounds := element.entity.Bounds()
	padding := element.entity.Theme().Padding(tomo.PatternGutter, element.c)
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
	bounds := element.entity.Bounds()
	padding := element.entity.Theme().Padding(tomo.PatternGutter, element.c)
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
	gutterPadding := element.entity.Theme().Padding(tomo.PatternGutter, element.c)
	handlePadding := element.entity.Theme().Padding(tomo.PatternHandle, element.c)
	if element.vertical {
		element.entity.SetMinimumSize (
			gutterPadding.Horizontal() + handlePadding.Horizontal(),
			gutterPadding.Vertical()   + handlePadding.Vertical() * 2)
	} else {
		element.entity.SetMinimumSize (
			gutterPadding.Horizontal() + handlePadding.Horizontal() * 2,
			gutterPadding.Vertical()   + handlePadding.Vertical())
	}
}
