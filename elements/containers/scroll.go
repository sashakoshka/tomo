package containers

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/elements/basic"

// ScrollContainer is a container that is capable of holding a scrollable
// element.
type ScrollContainer struct {
	*core.Core
	*core.Propagator
	core core.CoreControl
	
	child      elements.Scrollable
	horizontal *basicElements.ScrollBar
	vertical   *basicElements.ScrollBar
	
	config config.Wrapped
	theme  theme.Wrapped

	onFocusRequest func () (granted bool)
	onFocusMotionRequest func (input.KeynavDirection) (granted bool)
}

// NewScrollContainer creates a new scroll container with the specified scroll
// bars.
func NewScrollContainer (horizontal, vertical bool) (element *ScrollContainer) {
	element = &ScrollContainer { }
	element.theme.Case = theme.C("containers", "scrollContainer")
	element.Core, element.core = core.NewCore(element, element.redoAll)
	element.Propagator = core.NewPropagator(element, element.core)

	if horizontal {
		element.horizontal = basicElements.NewScrollBar(false)
		element.setUpChild(element.horizontal)
		element.horizontal.OnScroll (func (viewport image.Point) {
			if element.child != nil {
				element.child.ScrollTo(viewport)
			}
			if element.vertical != nil {
				element.vertical.SetBounds (
					element.child.ScrollContentBounds(),
					element.child.ScrollViewportBounds())
			}
		})
	}
	if vertical {
		element.vertical = basicElements.NewScrollBar(true)
		element.setUpChild(element.vertical)
		element.vertical.OnScroll (func (viewport image.Point) {
			if element.child != nil {
				element.child.ScrollTo(viewport)
			}
			if element.horizontal != nil {
				element.horizontal.SetBounds (
					element.child.ScrollContentBounds(),
					element.child.ScrollViewportBounds())
			}
		})
	}
	return
}


// Adopt adds a scrollable element to the scroll container. The container can
// only contain one scrollable element at a time, and when a new one is adopted
// it replaces the last one.
func (element *ScrollContainer) Adopt (child elements.Scrollable) {
	// disown previous child if it exists
	if element.child != nil {
		element.disownChild(child)
	}

	// adopt new child
	element.child = child
	if child != nil {
		element.setUpChild(child)
	}
	
	element.updateEnabled()
	element.updateMinimumSize()
	if element.core.HasImage() {
		element.redoAll()
		element.core.DamageAll()
	}
}

func (element *ScrollContainer) setUpChild (child elements.Element) {
	child.SetParent(element)
	if child, ok := child.(elements.Themeable); ok {
		child.SetTheme(element.theme.Theme)
	}
	if child, ok := child.(elements.Configurable); ok {
		child.SetConfig(element.config.Config)
	}
}

func (element *ScrollContainer) disownChild (child elements.Scrollable) {
	child.DrawTo(nil, image.Rectangle { }, nil)
	child.SetParent(nil)
	if child, ok := child.(elements.Focusable); ok {
		if child.Focused() {
			child.HandleUnfocus()
		}
	}
}

// NotifyMinimumSizeChange notifies the container that the minimum size of a
// child element has changed.
func (element *ScrollContainer) NotifyMinimumSizeChange (child elements.Element) {
	element.redoAll()
	element.core.DamageAll()
}

// NotifyScrollBoundsChange notifies the container that the scroll bounds or
// axes of a child have changed.
func (element *ScrollContainer) NotifyScrollBoundsChange (child elements.Scrollable) {
	element.updateEnabled()
	viewportBounds := element.child.ScrollViewportBounds()
	contentBounds  := element.child.ScrollContentBounds()
	if element.horizontal != nil {
		element.horizontal.SetBounds(contentBounds, viewportBounds)
	}
	if element.vertical != nil {
		element.vertical.SetBounds(contentBounds, viewportBounds)
	}
}

// SetTheme sets the element's theme.
func (element *ScrollContainer) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.Propagator.SetTheme(new)
	element.updateMinimumSize()
	element.redoAll()
}

// SetConfig sets the element's configuration.
func (element *ScrollContainer) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.Propagator.SetConfig(new)
	element.updateMinimumSize()
	element.redoAll()
}

func (element *ScrollContainer) HandleScroll (
	x, y int,
	deltaX, deltaY float64,
) {
	element.scrollChildBy(int(deltaX), int(deltaY))
}

// CountChildren returns the amount of children contained within this element.
func (element *ScrollContainer) CountChildren () (count int) {
	return 3
}

// Child returns the child at the specified index. If the index is out of
// bounds, this method will return nil.
func (element *ScrollContainer) Child (index int) (child elements.Element) {
	switch index {
	case 0: return element.child
	case 1:
		if element.horizontal == nil {
			return nil
		} else {
			return element.horizontal
		}
	case 2:
		if element.vertical == nil {
			return nil
		} else {
			return element.vertical
		}
	default: return nil
	}
}

func (element *ScrollContainer) redoAll () {
	if !element.core.HasImage() { return }

	zr := image.Rectangle { }
	if element.child      != nil { element.child.DrawTo(nil, zr, nil)      }
	if element.horizontal != nil { element.horizontal.DrawTo(nil, zr, nil) }
	if element.vertical   != nil { element.vertical.DrawTo(nil, zr, nil)   }
	
	childBounds, horizontalBounds, verticalBounds := element.layout()
	if element.child != nil {
		element.child.DrawTo (
			canvas.Cut(element.core, childBounds),
			childBounds, element.childDamageCallback)
	}
	if element.horizontal != nil {
		element.horizontal.DrawTo (
			canvas.Cut(element.core, horizontalBounds),
			horizontalBounds, element.childDamageCallback)
	}
	if element.vertical != nil {
		element.vertical.DrawTo (
			canvas.Cut(element.core, verticalBounds),
			verticalBounds, element.childDamageCallback)
	}
	element.draw()
}

func (element *ScrollContainer) scrollChildBy (x, y int) {
	if element.child == nil { return }
	scrollPoint :=
		element.child.ScrollViewportBounds().Min.
		Add(image.Pt(x, y))
	element.child.ScrollTo(scrollPoint)
}

func (element *ScrollContainer) childDamageCallback (region image.Rectangle) {
	element.core.DamageRegion(region)
}

func (element *ScrollContainer) layout () (
	child      image.Rectangle,
	horizontal image.Rectangle,
	vertical   image.Rectangle,
) {
	bounds := element.Bounds()
	child = bounds

	if element.horizontal != nil {
		_, hMinHeight := element.horizontal.MinimumSize()
		child.Max.Y -= hMinHeight
	}
	if element.vertical != nil {
		vMinWidth, _  := element.vertical.MinimumSize()
		child.Max.X -= vMinWidth
	}
	
	vertical.Min.X = child.Max.X
	vertical.Max.X = bounds.Max.X
	vertical.Min.Y = bounds.Min.Y
	vertical.Max.Y = child.Max.Y

	horizontal.Min.X = bounds.Min.X
	horizontal.Max.X = child.Max.X
	horizontal.Min.Y = child.Max.Y
	horizontal.Max.Y = bounds.Max.Y
	return
}

func (element *ScrollContainer) draw () {
	if element.horizontal != nil && element.vertical != nil {
		bounds := element.Bounds()
		bounds.Min = image.Pt (
			bounds.Max.X - element.vertical.Bounds().Dx(),
			bounds.Max.Y - element.horizontal.Bounds().Dy())
		state := theme.State { }
		deadArea := element.theme.Pattern(theme.PatternDead, state)
		deadArea.Draw(canvas.Cut(element.core, bounds), bounds)
	}
}

func (element *ScrollContainer) updateMinimumSize () {
	var width, height int

	if element.child != nil {
		width, height = element.child.MinimumSize()
	}
	if element.horizontal != nil {
		hMinWidth, hMinHeight := element.horizontal.MinimumSize()
		height += hMinHeight
		if hMinWidth > width {
			width = hMinWidth
		}
	}
	if element.vertical != nil {
		vMinWidth, vMinHeight := element.vertical.MinimumSize()
		width += vMinWidth
		if vMinHeight > height {
			height = vMinHeight
		}
	}
	element.core.SetMinimumSize(width, height)
}

func (element *ScrollContainer) updateEnabled () {
	horizontal, vertical := element.child.ScrollAxes()
	if element.horizontal != nil {
		element.horizontal.SetEnabled(horizontal)
	}
	if element.vertical != nil {
		element.vertical.SetEnabled(vertical)
	}
}
