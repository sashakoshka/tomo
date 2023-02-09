package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// ScrollContainer is a container that is capable of holding a scrollable
// element.
type ScrollContainer struct {
	*core.Core
	core core.CoreControl
	focused bool
	
	child elements.Scrollable
	childWidth, childHeight int
	
	horizontal struct {
		theme theme.Wrapped
		exists bool
		enabled bool
		dragging bool
		dragOffset int
		gutter image.Rectangle
		track image.Rectangle
		bar image.Rectangle
	}

	vertical struct {
		theme theme.Wrapped
		exists bool
		enabled bool
		dragging bool
		dragOffset int
		gutter image.Rectangle
		track image.Rectangle
		bar image.Rectangle
	}
	
	config config.Wrapped
	theme  theme.Wrapped

	onFocusRequest func () (granted bool)
	onFocusMotionRequest func (input.KeynavDirection) (granted bool)
}

// NewScrollContainer creates a new scroll container with the specified scroll
// bars.
func NewScrollContainer (horizontal, vertical bool) (element *ScrollContainer) {
	element = &ScrollContainer { }
	element.theme.Case = theme.C("basic", "scrollContainer")
	element.horizontal.theme.Case = theme.C("basic", "scrollBarHorizontal")
	element.vertical.theme.Case   = theme.C("basic", "scrollBarVertical")
	
	element.Core, element.core = core.NewCore(element.handleResize)
	element.horizontal.exists = horizontal
	element.vertical.exists   = vertical
	return
}

func (element *ScrollContainer) handleResize () {
	element.recalculate()
	element.resizeChildToFit()
	element.draw()
}

// Adopt adds a scrollable element to the scroll container. The container can
// only contain one scrollable element at a time, and when a new one is adopted
// it replaces the last one.
func (element *ScrollContainer) Adopt (child elements.Scrollable) {
	// disown previous child if it exists
	if element.child != nil {
		element.clearChildEventHandlers(child)
	}

	// adopt new child
	element.child = child
	if child != nil {
		if child0, ok := child.(elements.Themeable); ok {
			child0.SetTheme(element.theme.Theme)
		}
		if child0, ok := child.(elements.Configurable); ok {
			child0.SetConfig(element.config.Config)
		}
		child.OnDamage(element.childDamageCallback)
		child.OnMinimumSizeChange(element.updateMinimumSize)
		child.OnScrollBoundsChange(element.childScrollBoundsChangeCallback)
		if newChild, ok := child.(elements.Focusable); ok {
			newChild.OnFocusRequest (
				element.childFocusRequestCallback)
			newChild.OnFocusMotionRequest (
				element.childFocusMotionRequestCallback)
		}

		element.updateMinimumSize()
		
		element.horizontal.enabled,
		element.vertical.enabled = element.child.ScrollAxes()

		if element.core.HasImage() {
			element.resizeChildToFit()
		}
	}
}

// SetTheme sets the element's theme.
func (element *ScrollContainer) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	if child, ok := element.child.(elements.Themeable); ok {
		child.SetTheme(element.theme.Theme)
	}
	if element.core.HasImage() {
		element.recalculate()
		element.resizeChildToFit()
		element.draw()
	}
}

// SetConfig sets the element's configuration.
func (element *ScrollContainer) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	if child, ok := element.child.(elements.Configurable); ok {
		child.SetConfig(element.config.Config)
	}
	if element.core.HasImage() {
		element.recalculate()
		element.resizeChildToFit()
		element.draw()
	}
}

func (element *ScrollContainer) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	if child, ok := element.child.(elements.KeyboardTarget); ok {
		child.HandleKeyDown(key, modifiers)
	}
}

func (element *ScrollContainer) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	if child, ok := element.child.(elements.KeyboardTarget); ok {
		child.HandleKeyUp(key, modifiers)
	}
}

func (element *ScrollContainer) HandleMouseDown (x, y int, button input.Button) {
	velocity := element.config.ScrollVelocity()
	point := image.Pt(x, y)
	if point.In(element.horizontal.bar) {
		element.horizontal.dragging = true
		element.horizontal.dragOffset =
			x - element.horizontal.bar.Min.X +
			element.Bounds().Min.X
		element.dragHorizontalBar(point)
		
	} else if point.In(element.horizontal.gutter) {
		// FIXME: x backend and scroll container should pull these
		// values from the same place
		if x > element.horizontal.bar.Min.X {
			element.scrollChildBy(velocity, 0)
		} else {
			element.scrollChildBy(-velocity, 0)
		}
		
	} else if point.In(element.vertical.bar) {
		element.vertical.dragging = true
		element.vertical.dragOffset =
			y - element.vertical.bar.Min.Y +
				element.Bounds().Min.Y
		element.dragVerticalBar(point)
		
	} else if point.In(element.vertical.gutter) {
		if y > element.vertical.bar.Min.Y {
			element.scrollChildBy(0, velocity)
		} else {
			element.scrollChildBy(0, -velocity)
		}
		
	} else if child, ok := element.child.(elements.MouseTarget); ok {
		child.HandleMouseDown(x, y, button)
	}
}

func (element *ScrollContainer) HandleMouseUp (x, y int, button input.Button) {
	if element.horizontal.dragging {
		element.horizontal.dragging = false
		element.drawHorizontalBar()
		element.core.DamageRegion(element.horizontal.bar)
		
	} else if element.vertical.dragging {
		element.vertical.dragging = false
		element.drawVerticalBar()
		element.core.DamageRegion(element.vertical.bar)
		
	} else if child, ok := element.child.(elements.MouseTarget); ok {
		child.HandleMouseUp(x, y, button)
	}
}

func (element *ScrollContainer) HandleMouseMove (x, y int) {
	if element.horizontal.dragging {
		element.dragHorizontalBar(image.Pt(x, y))
		
	} else if element.vertical.dragging {
		element.dragVerticalBar(image.Pt(x, y))
		
	} else if child, ok := element.child.(elements.MouseTarget); ok {
		child.HandleMouseMove(x, y)
	}
}

func (element *ScrollContainer) HandleMouseScroll (
	x, y int,
	deltaX, deltaY float64,
) {
	element.scrollChildBy(int(deltaX), int(deltaY))
}

func (element *ScrollContainer) scrollChildBy (x, y int) {
	if element.child == nil { return }
	scrollPoint :=
		element.child.ScrollViewportBounds().Min.
		Add(image.Pt(x, y))
	element.child.ScrollTo(scrollPoint)
}

func (element *ScrollContainer) Focused () (focused bool) {
	return element.focused
}

func (element *ScrollContainer) Focus () {
	if element.onFocusRequest != nil {
		element.onFocusRequest()
	}
}

func (element *ScrollContainer) HandleFocus (
	direction input.KeynavDirection,
) (
	accepted bool,
) {
	if child, ok := element.child.(elements.Focusable); ok {
		element.focused = true
		return child.HandleFocus(direction)
	} else {
		element.focused = false
		return false
	}
}

func (element *ScrollContainer) HandleUnfocus () {
	if child, ok := element.child.(elements.Focusable); ok {
		child.HandleUnfocus()
	}
	element.focused = false
}

func (element *ScrollContainer) OnFocusRequest (callback func () (granted bool)) {
	element.onFocusRequest = callback
}

func (element *ScrollContainer) OnFocusMotionRequest (
	callback func (direction input.KeynavDirection) (granted bool),
) {
	element.onFocusMotionRequest = callback
}

func (element *ScrollContainer) childDamageCallback (region canvas.Canvas) {
	element.core.DamageRegion(artist.Paste(element, region, image.Point { }))
}

func (element *ScrollContainer) childFocusRequestCallback () (granted bool) {
	child, ok := element.child.(elements.Focusable)
	if !ok { return false }
	if element.onFocusRequest != nil && element.onFocusRequest() {
		child.HandleFocus(input.KeynavDirectionNeutral)
		return true
	} else {
		return false
	}
}

func (element *ScrollContainer) childFocusMotionRequestCallback (
	direction input.KeynavDirection,
) (
	granted bool,
) {
	if element.onFocusMotionRequest == nil { return }
	return element.onFocusMotionRequest(direction)
}

func (element *ScrollContainer) clearChildEventHandlers (child elements.Scrollable) {
	child.DrawTo(nil)
	child.OnDamage(nil)
	child.OnMinimumSizeChange(nil)
	child.OnScrollBoundsChange(nil)
	if child0, ok := child.(elements.Focusable); ok {
		child0.OnFocusRequest(nil)
		child0.OnFocusMotionRequest(nil)
		if child0.Focused() {
			child0.HandleUnfocus()
		}
	}
	if child0, ok := child.(elements.Flexible); ok {
		child0.OnFlexibleHeightChange(nil)
	}
}

func (element *ScrollContainer) resizeChildToFit () {
	childBounds := image.Rect (
		0, 0,
		element.childWidth,
		element.childHeight).Add(element.Bounds().Min)
	element.child.DrawTo(canvas.Cut(element, childBounds))
}

func (element *ScrollContainer) recalculate () {
	horizontal := &element.horizontal
	vertical   := &element.vertical
	
	gutterInsetHorizontal := horizontal.theme.Inset(theme.PatternGutter)
	gutterInsetVertical   := vertical.theme.Inset(theme.PatternGutter)

	bounds     := element.Bounds()
	thicknessHorizontal :=
		element.config.HandleWidth() +
		gutterInsetHorizontal[3] +
		gutterInsetHorizontal[1]
	thicknessVertical :=
		element.config.HandleWidth() +
		gutterInsetVertical[3] +
		gutterInsetVertical[1]

	// calculate child size
	element.childWidth  = bounds.Dx()
	element.childHeight = bounds.Dy()

	// reset bounds
	horizontal.gutter = image.Rectangle { }
	vertical.gutter   = image.Rectangle { }
	horizontal.bar    = image.Rectangle { }
	vertical.bar      = image.Rectangle { }

	// if enabled, give substance to the gutters
	if horizontal.exists {
		horizontal.gutter.Min.X = bounds.Min.X
		horizontal.gutter.Min.Y = bounds.Max.Y - thicknessHorizontal
		horizontal.gutter.Max.X = bounds.Max.X
		horizontal.gutter.Max.Y = bounds.Max.Y
		if vertical.exists {
			horizontal.gutter.Max.X -= thicknessVertical
		}
		element.childHeight -= thicknessHorizontal
		horizontal.track = gutterInsetHorizontal.Apply(horizontal.gutter)
	}
	if vertical.exists {
		vertical.gutter.Min.X = bounds.Max.X - thicknessVertical
		vertical.gutter.Max.X = bounds.Max.X
		vertical.gutter.Min.Y = bounds.Min.Y
		vertical.gutter.Max.Y = bounds.Max.Y
		if horizontal.exists {
			vertical.gutter.Max.Y -= thicknessHorizontal
		}
		element.childWidth -= thicknessVertical
		vertical.track = gutterInsetVertical.Apply(vertical.gutter)
	}

	// if enabled, calculate the positions of the bars
	contentBounds  := element.child.ScrollContentBounds()
	viewportBounds := element.child.ScrollViewportBounds()
	if horizontal.exists && horizontal.enabled {
		horizontal.bar.Min.Y = horizontal.track.Min.Y
		horizontal.bar.Max.Y = horizontal.track.Max.Y

		scale := float64(horizontal.track.Dx()) /
			float64(contentBounds.Dx())
		horizontal.bar.Min.X = int(float64(viewportBounds.Min.X) * scale)
		horizontal.bar.Max.X = int(float64(viewportBounds.Max.X) * scale)
		
		horizontal.bar.Min.X += horizontal.track.Min.X
		horizontal.bar.Max.X += horizontal.track.Min.X
	}
	if vertical.exists && vertical.enabled {
		vertical.bar.Min.X = vertical.track.Min.X
		vertical.bar.Max.X = vertical.track.Max.X

		scale := float64(vertical.track.Dy()) /
			float64(contentBounds.Dy())
		vertical.bar.Min.Y = int(float64(viewportBounds.Min.Y) * scale)
		vertical.bar.Max.Y = int(float64(viewportBounds.Max.Y) * scale)
		
		vertical.bar.Min.Y += vertical.track.Min.Y
		vertical.bar.Max.Y += vertical.track.Min.Y
	}

	// if the scroll bars are out of bounds, don't display them.
	if horizontal.bar.Dx() >= horizontal.track.Dx() {
		horizontal.bar = image.Rectangle { }
	}
	if vertical.bar.Dy() >= vertical.track.Dy() {
		vertical.bar = image.Rectangle { }
	}
}

func (element *ScrollContainer) draw () {
	artist.Paste(element, element.child, image.Point { })
	deadPattern := element.theme.Pattern (
		theme.PatternDead, theme.PatternState { })
	artist.FillRectangle (
		element, deadPattern,
		image.Rect (
			element.vertical.gutter.Min.X,
			element.horizontal.gutter.Min.Y,
			element.vertical.gutter.Max.X,
			element.horizontal.gutter.Max.Y))
	element.drawHorizontalBar()
	element.drawVerticalBar()
}

func (element *ScrollContainer) drawHorizontalBar () {
	state := theme.PatternState {
		Disabled: !element.horizontal.enabled,
		Pressed:  element.horizontal.dragging,
	}
	gutterPattern := element.horizontal.theme.Pattern(theme.PatternGutter, state)
	artist.FillRectangle(element, gutterPattern, element.horizontal.gutter)
	
	handlePattern := element.horizontal.theme.Pattern(theme.PatternHandle, state)
	artist.FillRectangle(element, handlePattern, element.horizontal.bar)
}

func (element *ScrollContainer) drawVerticalBar () {
	state := theme.PatternState {
		Disabled: !element.vertical.enabled,
		Pressed:  element.vertical.dragging,
	}
	gutterPattern := element.vertical.theme.Pattern(theme.PatternGutter, state)
	artist.FillRectangle(element, gutterPattern, element.vertical.gutter)
	
	handlePattern := element.vertical.theme.Pattern(theme.PatternHandle, state)
	artist.FillRectangle(element, handlePattern, element.vertical.bar)
}

func (element *ScrollContainer) dragHorizontalBar (mousePosition image.Point) {
	scrollX :=
		float64(element.child.ScrollContentBounds().Dx()) /
		float64(element.horizontal.track.Dx()) *
		float64(mousePosition.X - element.horizontal.dragOffset)
	scrollY := element.child.ScrollViewportBounds().Min.Y
	element.child.ScrollTo(image.Pt(int(scrollX), scrollY))
}

func (element *ScrollContainer) dragVerticalBar (mousePosition image.Point) {
	scrollY :=
		float64(element.child.ScrollContentBounds().Dy()) /
		float64(element.vertical.track.Dy()) *
		float64(mousePosition.Y - element.vertical.dragOffset)
	scrollX := element.child.ScrollViewportBounds().Min.X
	element.child.ScrollTo(image.Pt(scrollX, int(scrollY)))
}

func (element *ScrollContainer) updateMinimumSize () {
	gutterInsetHorizontal := element.horizontal.theme.Inset(theme.PatternGutter)
	gutterInsetVertical   := element.vertical.theme.Inset(theme.PatternGutter)

	thicknessHorizontal :=
		element.config.HandleWidth() +
		gutterInsetHorizontal[3] +
		gutterInsetHorizontal[1]
	thicknessVertical :=
		element.config.HandleWidth() +
		gutterInsetVertical[3] +
		gutterInsetVertical[1]
	
	width  := thicknessHorizontal
	height := thicknessVertical
	if element.child != nil {
		childWidth, childHeight := element.child.MinimumSize()
		width  += childWidth
		height += childHeight
	}
	element.core.SetMinimumSize(width, height)
}

func (element *ScrollContainer) childScrollBoundsChangeCallback () {
	element.horizontal.enabled,
	element.vertical.enabled = element.child.ScrollAxes()
	if element.core.HasImage() {
		element.recalculate()
		element.drawHorizontalBar()
		element.drawVerticalBar()
		element.core.DamageRegion(element.horizontal.gutter)
		element.core.DamageRegion(element.vertical.gutter)
	}
}
