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
		element.child.SetParentHooks (tomo.ParentHooks { })
		if previousChild, ok := element.child.(tomo.Selectable); ok {
			if previousChild.Selected() {
				previousChild.HandleDeselection()
			}
		}
	}

	// adopt new child
	element.child = child
	if child != nil {
		child.SetParentHooks (tomo.ParentHooks {
			// Draw: window.childDrawCallback,
			// MinimumSizeChange: window.childMinimumSizeChangeCallback,
			// FlexibleHeightChange: window.resizeChildToFit,
			// SelectionRequest: window.childSelectionRequestCallback,
		})

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
