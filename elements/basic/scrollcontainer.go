package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// ScrollContainer is a container that is capable of holding a scrollable
// element.
type ScrollContainer struct {
	*core.Core
	core core.CoreControl
	selected bool
	
	child tomo.Scrollable
	childWidth, childHeight int
	
	horizontal struct {
		exists bool
		enabled bool
		dragging bool
		dragOffset image.Point
		gutter image.Rectangle
		bar image.Rectangle
	}

	vertical struct {
		exists bool
		enabled bool
		dragging bool
		dragOffset image.Point
		gutter image.Rectangle
		bar image.Rectangle
	}

	onSelectionRequest func () (granted bool)
	onSelectionMotionRequest func (tomo.SelectionDirection) (granted bool)
}

// NewScrollContainer creates a new scroll container with the specified scroll
// bars.
func NewScrollContainer (horizontal, vertical bool) (element *ScrollContainer) {
	element = &ScrollContainer { }
	element.Core, element.core = core.NewCore(element)
	element.updateMinimumSize()
	element.horizontal.exists = horizontal
	element.vertical.exists   = vertical
	return
}

// Resize resizes the scroll box.
func (element *ScrollContainer) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.recalculate()
	element.child.Resize (
		element.childWidth,
		element.childHeight)
	element.draw()
}

// Adopt adds a scrollable element to the scroll container. The container can
// only contain one scrollable element at a time, and when a new one is adopted
// it replaces the last one.
func (element *ScrollContainer) Adopt (child tomo.Scrollable) {
	// disown previous child if it exists
	if element.child != nil {
		element.clearChildEventHandlers(child)
	}

	// adopt new child
	element.child = child
	if child != nil {
		child.OnDamage(element.childDamageCallback)
		child.OnMinimumSizeChange(element.updateMinimumSize)
		child.OnScrollBoundsChange(element.childScrollBoundsChangeCallback)
		if newChild, ok := child.(tomo.Selectable); ok {
			newChild.OnSelectionRequest (
				element.childSelectionRequestCallback)
			newChild.OnSelectionMotionRequest (
				element.childSelectionMotionRequestCallback)
		}

		// TODO: somehow inform the core that we do not in fact want to
		// redraw the element.
		element.updateMinimumSize()
		
		element.horizontal.enabled,
		element.vertical.enabled = element.child.ScrollAxes()

		if element.core.HasImage() {
			element.child.Resize (
				element.childWidth,
				element.childHeight)
			element.core.DamageAll()
		}
	}
}

func (element *ScrollContainer) HandleKeyDown (key tomo.Key, modifiers tomo.Modifiers) {
	if child, ok := element.child.(tomo.KeyboardTarget); ok {
		child.HandleKeyDown(key, modifiers)
	}
}

func (element *ScrollContainer) HandleKeyUp (key tomo.Key, modifiers tomo.Modifiers) {
	if child, ok := element.child.(tomo.KeyboardTarget); ok {
		child.HandleKeyUp(key, modifiers)
	}
}

func (element *ScrollContainer) HandleMouseDown (x, y int, button tomo.Button) {
	point := image.Pt(x, y)
	if point.In(element.horizontal.gutter) {
		element.horizontal.dragging = true
		element.dragHorizontalBar(point)
		
	} else if point.In(element.vertical.gutter) {
		element.vertical.dragging = true
		element.dragVerticalBar(point)
		
	} else if child, ok := element.child.(tomo.MouseTarget); ok {
		child.HandleMouseDown(x, y, button)
	}
}

func (element *ScrollContainer) HandleMouseUp (x, y int, button tomo.Button) {
	if element.horizontal.dragging {
		element.horizontal.dragging = false
		
	} else if element.vertical.dragging {
		element.vertical.dragging = false
		
	} else if child, ok := element.child.(tomo.MouseTarget); ok {
		child.HandleMouseUp(x, y, button)
	}
}

func (element *ScrollContainer) HandleMouseMove (x, y int) {
	if element.horizontal.dragging {
		element.dragHorizontalBar(image.Pt(x, y))
		
	} else if element.vertical.dragging {
		element.dragVerticalBar(image.Pt(x, y))
		
	} else if child, ok := element.child.(tomo.MouseTarget); ok {
		child.HandleMouseMove(x, y)
	}
}

func (element *ScrollContainer) HandleMouseScroll (
	x, y int,
	deltaX, deltaY float64,
) {
	scrollPoint := element.child.ScrollViewportBounds().Min.Add(image.Pt (
		int(deltaX),
		int(deltaY)))
	element.child.ScrollTo(scrollPoint)
}

func (element *ScrollContainer) Selected () (selected bool) {
	return element.selected
}

func (element *ScrollContainer) Select () {
	if element.onSelectionRequest != nil {
		element.onSelectionRequest()
	}
}

func (element *ScrollContainer) HandleSelection (
	direction tomo.SelectionDirection,
) (
	accepted bool,
) {
	if child, ok := element.child.(tomo.Selectable); ok {
		element.selected = true
		return child.HandleSelection(direction)
	} else {
		element.selected = false
		return false
	}
}

func (element *ScrollContainer) HandleDeselection () {
	if child, ok := element.child.(tomo.Selectable); ok {
		child.HandleDeselection()
	}
	element.selected = false
}

func (element *ScrollContainer) OnSelectionRequest (callback func () (granted bool)) {
	element.onSelectionRequest = callback
}

func (element *ScrollContainer) OnSelectionMotionRequest (
	callback func (direction tomo.SelectionDirection) (granted bool),
) {
	element.onSelectionMotionRequest = callback
}

func (element *ScrollContainer) childDamageCallback (region tomo.Canvas) {
	element.core.DamageRegion(artist.Paste(element, region, image.Point { }))
}

func (element *ScrollContainer) childSelectionRequestCallback () (granted bool) {
	child, ok := element.child.(tomo.Selectable)
	if !ok { return false }
	if element.onSelectionRequest != nil && element.onSelectionRequest() {
		child.HandleSelection(tomo.SelectionDirectionNeutral)
		return true
	} else {
		return false
	}
}

func (element *ScrollContainer) childSelectionMotionRequestCallback (
	direction tomo.SelectionDirection,
) (
	granted bool,
) {
	if element.onSelectionMotionRequest == nil {
		 return
	}
	return element.onSelectionMotionRequest(direction)
}

func (element *ScrollContainer) clearChildEventHandlers (child tomo.Scrollable) {
	child.OnDamage(nil)
	child.OnMinimumSizeChange(nil)
	child.OnScrollBoundsChange(nil)
	if child0, ok := child.(tomo.Selectable); ok {
		child0.OnSelectionRequest(nil)
		child0.OnSelectionMotionRequest(nil)
		if child0.Selected() {
			child0.HandleDeselection()
		}
	}
	if child0, ok := child.(tomo.Flexible); ok {
		child0.OnFlexibleHeightChange(nil)
	}
}

func (element *ScrollContainer) recalculate () {
	horizontal := &element.horizontal
	vertical   := &element.vertical
	bounds     := element.Bounds()
	thickness  := theme.Padding() * 2

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
		horizontal.gutter.Min.Y = bounds.Max.Y - thickness
		horizontal.gutter.Max.X = bounds.Max.X
		horizontal.gutter.Max.Y = bounds.Max.Y
		if vertical.exists {
			horizontal.gutter.Max.X -= thickness
		}
		element.childHeight -= thickness
	}
	if vertical.exists {
		vertical.gutter.Min.X = bounds.Max.X - thickness
		vertical.gutter.Max.X = bounds.Max.X
		vertical.gutter.Max.Y = bounds.Max.Y
		if horizontal.exists {
			vertical.gutter.Max.Y -= thickness
		}
		element.childWidth -= thickness
	}

	// if enabled, calculate the positions of the bars
	contentBounds  := element.child.ScrollContentBounds()
	viewportBounds := element.child.ScrollViewportBounds()
	if horizontal.exists && horizontal.enabled {
		horizontal.bar.Min.Y = horizontal.gutter.Min.Y
		horizontal.bar.Max.Y = horizontal.gutter.Max.Y

		scale := float64(horizontal.gutter.Dx()) /
			float64(contentBounds.Dx())
		horizontal.bar.Min.X = int(float64(viewportBounds.Min.X) * scale)
		horizontal.bar.Max.X = int(float64(viewportBounds.Max.X) * scale)
		
		horizontal.bar.Min.X += horizontal.gutter.Min.X
		horizontal.bar.Max.X += horizontal.gutter.Min.X
	}
	if vertical.exists && vertical.enabled {
		vertical.bar.Min.X = vertical.gutter.Min.X
		vertical.bar.Max.X = vertical.gutter.Max.X

		scale := float64(vertical.gutter.Dy()) /
			float64(contentBounds.Dy())
		vertical.bar.Min.Y = int(float64(viewportBounds.Min.Y) * scale)
		vertical.bar.Max.Y = int(float64(viewportBounds.Max.Y) * scale)
		
		vertical.bar.Min.Y += vertical.gutter.Min.Y
		vertical.bar.Max.Y += vertical.gutter.Min.Y
	}

	// if the scroll bars are out of bounds, don't display them.
	if horizontal.bar.Dx() >= horizontal.gutter.Dx() {
		horizontal.bar = image.Rectangle { }
	}
	if vertical.bar.Dy() >= vertical.gutter.Dy() {
		vertical.bar = image.Rectangle { }
	}
}

func (element *ScrollContainer) draw () {
	artist.Paste(element.core, element.child, image.Point { })
	artist.FillRectangle (
		element, theme.DeadPattern(),
		image.Rect (
			element.vertical.gutter.Min.X,
			element.horizontal.gutter.Min.Y,
			element.vertical.gutter.Max.X,
			element.horizontal.gutter.Max.Y))
	element.drawHorizontalBar()
	element.drawVerticalBar()
}

func (element *ScrollContainer) drawHorizontalBar () {
	artist.FillRectangle (
		element,
		theme.ScrollGutterPattern(true, element.horizontal.enabled),
		element.horizontal.gutter)
	artist.FillRectangle (
		element,
		theme.ScrollBarPattern(true, element.horizontal.enabled),
		element.horizontal.bar)
}

func (element *ScrollContainer) drawVerticalBar () {
	artist.FillRectangle (
		element,
		theme.ScrollGutterPattern(false, element.vertical.enabled),
		element.vertical.gutter)
	artist.FillRectangle (
		element,
		theme.ScrollBarPattern(false, element.vertical.enabled),
		element.vertical.bar)
}

func (element *ScrollContainer) dragHorizontalBar (mousePosition image.Point) {
	scrollX :=
		float64(element.child.ScrollContentBounds().Dx()) /
		float64(element.horizontal.gutter.Dx()) *
		float64 (
			mousePosition.X - element.horizontal.bar.Dx() / 2)
	scrollY := element.child.ScrollViewportBounds().Min.Y
	element.child.ScrollTo(image.Pt(int(scrollX), scrollY))
}

func (element *ScrollContainer) dragVerticalBar (mousePosition image.Point) {
	scrollY :=
		float64(element.child.ScrollContentBounds().Dy()) /
		float64(element.vertical.gutter.Dy()) *
		float64 (
			mousePosition.Y - element.vertical.bar.Dy() / 2)
	scrollX := element.child.ScrollViewportBounds().Min.X
	element.child.ScrollTo(image.Pt(scrollX, int(scrollY)))
}

func (element *ScrollContainer) updateMinimumSize () {
	width  := theme.Padding() * 2
	height := theme.Padding() * 2
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
