package containers

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// DocumentContainer is a scrollable container capable of containing flexible
// elements.
type DocumentContainer struct {
	*core.Core
	*core.Propagator
	core core.CoreControl

	children []tomo.LayoutEntry
	scroll   image.Point
	warping  bool
	contentBounds image.Rectangle
	
	config config.Wrapped
	theme  theme.Wrapped

	onScrollBoundsChange func ()
}

// NewDocumentContainer creates a new document container.
func NewDocumentContainer () (element *DocumentContainer) {
	element = &DocumentContainer { }
	element.theme.Case = tomo.C("tomo", "documentContainer")
	element.Core, element.core = core.NewCore(element, element.redoAll)
	element.Propagator = core.NewPropagator(element, element.core)
	return
}

// Adopt adds a new child element to the container. If expand is true, then the
// element will stretch to either side of the container (much like a css block
// element). If expand is false, the element will share a line with other inline
// elements.
func (element *DocumentContainer) Adopt (child tomo.Element, expand bool) {
	// set event handlers
	if child0, ok := child.(tomo.Themeable); ok {
		child0.SetTheme(element.theme.Theme)
	}
	if child0, ok := child.(tomo.Configurable); ok {
		child0.SetConfig(element.config.Config)
	}

	// add child
	element.children = append (element.children, tomo.LayoutEntry {
		Element: child,
		Expand:  expand,
	})

	child.SetParent(element)

	// refresh stale data
	element.updateMinimumSize()
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

// Warp runs the specified callback, deferring all layout and rendering updates
// until the callback has finished executing. This allows for aplications to
// perform batch gui updates without flickering and stuff.
func (element *DocumentContainer) Warp (callback func ()) {
	if element.warping {
		callback()
		return
	}

	element.warping = true
	callback()
	element.warping = false
	
	if element.core.HasImage() {
		element.redoAll()
		element.core.DamageAll()
	}
}

// Disown removes the given child from the container if it is contained within
// it.
func (element *DocumentContainer) Disown (child tomo.Element) {
	for index, entry := range element.children {
		if entry.Element == child {
			element.clearChildEventHandlers(entry.Element)
			element.children = append (
				element.children[:index],
				element.children[index + 1:]...)
				break
		}
	}

	element.updateMinimumSize()
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

func (element *DocumentContainer) clearChildEventHandlers (child tomo.Element) {
	child.DrawTo(nil, image.Rectangle { }, nil)
	child.SetParent(nil)
	
	if child, ok := child.(tomo.Focusable); ok {
		if child.Focused() {
			child.HandleUnfocus()
		}
	}
}

// DisownAll removes all child elements from the container at once.
func (element *DocumentContainer) DisownAll () {
	for _, entry := range element.children {
		element.clearChildEventHandlers(entry.Element)
	}
	element.children = nil

	element.updateMinimumSize()
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

// Children returns a slice containing this element's children.
func (element *DocumentContainer) Children () (children []tomo.Element) {
	children = make([]tomo.Element, len(element.children))
	for index, entry := range element.children {
		children[index] = entry.Element
	}
	return
}

// CountChildren returns the amount of children contained within this element.
func (element *DocumentContainer) CountChildren () (count int) {
	return len(element.children)
}

// Child returns the child at the specified index. If the index is out of
// bounds, this method will return nil.
func (element *DocumentContainer) Child (index int) (child tomo.Element) {
	if index < 0 || index > len(element.children) { return }
	return element.children[index].Element
}

// ChildAt returns the child that contains the specified x and y coordinates. If
// there are no children at the coordinates, this method will return nil.
func (element *DocumentContainer) ChildAt (point image.Point) (child tomo.Element) {
	for _, entry := range element.children {
		if point.In(entry.Bounds) {
			child = entry.Element
		}
	}
	return
}

func (element *DocumentContainer) redoAll () {
	if !element.core.HasImage() { return }
	
	// do a layout
	element.doLayout()
	
	maxScrollHeight := element.maxScrollHeight()
	if element.scroll.Y > maxScrollHeight {
		element.scroll.Y = maxScrollHeight
		element.doLayout()
	}

	// draw a background
	rocks := make([]image.Rectangle, len(element.children))
	for index, entry := range element.children {
		rocks[index] = entry.Bounds
	}
	pattern := element.theme.Pattern (
		tomo.PatternBackground,
		tomo.State { })
	artist.DrawShatter(element.core, pattern, element.Bounds(), rocks...)

	element.partition()
	if parent, ok := element.core.Parent().(tomo.ScrollableParent); ok {
		parent.NotifyScrollBoundsChange(element)
	}
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *DocumentContainer) partition () {
	for _, entry := range element.children {
		entry.DrawTo(nil, entry.Bounds, nil)
	}

	// cut our canvas up and give peices to child elements
	for _, entry := range element.children {
		if entry.Bounds.Overlaps(element.Bounds()) {
			entry.DrawTo (	
				canvas.Cut(element.core, entry.Bounds),
				entry.Bounds, func (region image.Rectangle) {
					element.core.DamageRegion(region)
				})
		}
	}
}

func (element *DocumentContainer) Window () tomo.Window {
	return element.core.Window()
}

// NotifyMinimumSizeChange notifies the container that the minimum size of a
// child element has changed.
func (element *DocumentContainer) NotifyMinimumSizeChange (child tomo.Element) {
	element.redoAll()
	element.core.DamageAll()
}

// DrawBackground draws a portion of the container's background pattern within
// the specified bounds. The container will not push these changes.
func (element *DocumentContainer) DrawBackground (bounds image.Rectangle) {
	element.core.DrawBackgroundBounds (
		element.theme.Pattern(tomo.PatternBackground, tomo.State { }),
		bounds)
}

// NotifyFlexibleHeightChange notifies the parent that the parameters
// affecting a child's flexible height have changed. This method is
// expected to be called by flexible child element when their content
// changes.
func (element *DocumentContainer) NotifyFlexibleHeightChange (child tomo.Flexible) {
	element.redoAll()
	element.core.DamageAll()
}

// SetTheme sets the element's theme.
func (element *DocumentContainer) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.Propagator.SetTheme(new)
	element.redoAll()
}

// SetConfig sets the element's configuration.
func (element *DocumentContainer) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.Propagator.SetConfig(new)
	element.redoAll()
}

// ScrollContentBounds returns the full content size of the element.
func (element *DocumentContainer) ScrollContentBounds () image.Rectangle {
	return element.contentBounds
}

// ScrollViewportBounds returns the size and position of the element's
// viewport relative to ScrollBounds.
func (element *DocumentContainer) ScrollViewportBounds () image.Rectangle {
	padding := element.theme.Padding(tomo.PatternBackground)
	bounds  := padding.Apply(element.Bounds())
	bounds   = bounds.Sub(bounds.Min).Add(element.scroll)
	return bounds
}

// ScrollTo scrolls the viewport to the specified point relative to
// ScrollBounds.
func (element *DocumentContainer) ScrollTo (position image.Point) {
	if position.Y < 0 {
		position.Y = 0
	}
	maxScrollHeight := element.maxScrollHeight()
	if position.Y > maxScrollHeight {
		position.Y = maxScrollHeight
	}
	element.scroll = position
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

// OnScrollBoundsChange sets a function to be called when the element's viewport
// bounds, content bounds, or scroll axes change.
func (element *DocumentContainer) OnScrollBoundsChange (callback func ()) {
	element.onScrollBoundsChange = callback
}

func (element *DocumentContainer) maxScrollHeight () (height int) {
	padding := element.theme.Padding(tomo.PatternSunken)
	viewportHeight := element.Bounds().Dy() - padding.Vertical()
	height = element.contentBounds.Dy() - viewportHeight
	if height < 0 { height = 0 }
	return
}

// ScrollAxes returns the supported axes for scrolling.
func (element *DocumentContainer) ScrollAxes () (horizontal, vertical bool) {
	return false, true
}

func (element *DocumentContainer) doLayout () {
	margin := element.theme.Margin(tomo.PatternBackground)
	padding := element.theme.Padding(tomo.PatternBackground)
	bounds := padding.Apply(element.Bounds())
	element.contentBounds = image.Rectangle { }

	dot := bounds.Min.Sub(element.scroll)
	xStart := dot.X
	rowHeight := 0

	nextLine := func () {
		dot.X = xStart
		dot.Y += margin.Y
		dot.Y += rowHeight
		rowHeight = 0
	}
	
	for index, entry := range element.children {
		if dot.X > xStart && entry.Expand {
			nextLine()
		}
	
		width, height := entry.MinimumSize()
		if width + dot.X > bounds.Dx() && !entry.Expand {
			nextLine()
		}
		if width < bounds.Dx() && entry.Expand {
			width = bounds.Dx()
		}
		if typedChild, ok := entry.Element.(tomo.Flexible); ok {
			height = typedChild.FlexibleHeightFor(width)
		}
		if rowHeight < height {
			rowHeight = height
		}
		
		entry.Bounds.Min = dot
		entry.Bounds.Max = image.Pt(dot.X + width, dot.Y + height)
		element.children[index] = entry
		element.contentBounds = element.contentBounds.Union(entry.Bounds)
		
		if entry.Expand {
			nextLine()
		} else {
			dot.X += width + margin.X
		}
	}
	
	element.contentBounds =
		element.contentBounds.Sub(element.contentBounds.Min)
}

func (element *DocumentContainer) updateMinimumSize () {
	padding := element.theme.Padding(tomo.PatternBackground)
	minimumWidth := 0
	for _, entry := range element.children {
		width, _ := entry.MinimumSize()
		if width > minimumWidth {
			minimumWidth = width
		}
	}
	element.core.SetMinimumSize (
		minimumWidth + padding.Horizontal(),
		padding.Vertical())
}
