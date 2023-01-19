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

	horizontal struct {
		enabled bool
		bounds image.Rectangle
	}

	vertical struct {
		enabled bool
		bounds image.Rectangle
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
	element.horizontal.enabled = horizontal
	element.vertical.enabled   = vertical
	return
}

// Resize resizes the scroll box.
func (element *ScrollContainer) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.recalculate()
	element.child.Resize (
		element.vertical.bounds.Min.X,
		element.horizontal.bounds.Min.Y)
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
		if newChild, ok := child.(tomo.Selectable); ok {
			newChild.OnSelectionRequest (
				element.childSelectionRequestCallback)
			newChild.OnSelectionMotionRequest (
				element.childSelectionMotionRequestCallback)
		}

		// TODO: somehow inform the core that we do not in fact want to
		// redraw the element.
		element.updateMinimumSize()

		if element.core.HasImage() {
			element.recalculate()
			element.child.Resize (
				element.horizontal.bounds.Min.X,
				element.vertical.bounds.Min.X)
			element.draw()
		}
	}
}

func (element *ScrollContainer) HandleKeyDown (
	key tomo.Key,
	modifiers tomo.Modifiers,
	repeated bool,
) {
	if child, ok := element.child.(tomo.KeyboardTarget); ok {
		child.HandleKeyDown(key, modifiers, repeated)
	}
}

func (element *ScrollContainer) HandleKeyUp (key tomo.Key, modifiers tomo.Modifiers) {
	if child, ok := element.child.(tomo.KeyboardTarget); ok {
		child.HandleKeyUp(key, modifiers)
	}
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

func (element *ScrollContainer) clearChildEventHandlers (child tomo.Element) {
	child.OnDamage(nil)
	child.OnMinimumSizeChange(nil)
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

	// reset bounds
	horizontal.bounds = image.Rectangle { }
	vertical.bounds   = image.Rectangle { }

	// if enabled, give substance to the bars
	if horizontal.enabled {
		horizontal.bounds.Max.X = bounds.Max.X - thickness
		horizontal.bounds.Max.Y = thickness
	}
	if vertical.enabled {
		vertical.bounds.Max.X = thickness
		vertical.bounds.Max.Y = bounds.Max.Y - thickness
	}

	// move the bars to the edge of the element
	horizontal.bounds = horizontal.bounds.Add (
		bounds.Max.Sub(horizontal.bounds.Max))
	vertical.bounds = vertical.bounds.Add (
		bounds.Max.Sub(vertical.bounds.Max))
}

func (element *ScrollContainer) draw () {
	artist.Paste(element.core, element.child, image.Point { })
	element.drawHorizontalBar()
	element.drawVerticalBar()
}

func (element *ScrollContainer) drawHorizontalBar () {
	
}

func (element *ScrollContainer) drawVerticalBar () {
	
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
