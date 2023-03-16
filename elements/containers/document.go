package containers

import "image"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type DocumentContainer struct {
	*core.Core
	*core.Propagator
	core core.CoreControl

	children []layouts.LayoutEntry
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
	element.theme.Case = theme.C("containers", "documentContainer")
	element.Core, element.core = core.NewCore(element, element.redoAll)
	element.Propagator = core.NewPropagator(element, element.core)
	return
}

// Adopt adds a new child element to the container.
func (element *DocumentContainer) Adopt (child elements.Element) {
	// set event handlers
	if child0, ok := child.(elements.Themeable); ok {
		child0.SetTheme(element.theme.Theme)
	}
	if child0, ok := child.(elements.Configurable); ok {
		child0.SetConfig(element.config.Config)
	}

	// add child
	element.children = append (element.children, layouts.LayoutEntry {
		Element: child,
	})

	child.SetParent(element)

	// refresh stale data
	element.reflectChildProperties()
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
func (element *DocumentContainer) Disown (child elements.Element) {
	for index, entry := range element.children {
		if entry.Element == child {
			element.clearChildEventHandlers(entry.Element)
			element.children = append (
				element.children[:index],
				element.children[index + 1:]...)
				break
		}
	}

	element.reflectChildProperties()
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

func (element *DocumentContainer) clearChildEventHandlers (child elements.Element) {
	child.DrawTo(nil, image.Rectangle { }, nil)
	child.SetParent(nil)
	
	if child, ok := child.(elements.Focusable); ok {
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

	element.reflectChildProperties()
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

// Children returns a slice containing this element's children.
func (element *DocumentContainer) Children () (children []elements.Element) {
	children = make([]elements.Element, len(element.children))
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
func (element *DocumentContainer) Child (index int) (child elements.Element) {
	if index < 0 || index > len(element.children) { return }
	return element.children[index].Element
}

// ChildAt returns the child that contains the specified x and y coordinates. If
// there are no children at the coordinates, this method will return nil.
func (element *DocumentContainer) ChildAt (point image.Point) (child elements.Element) {
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
		theme.PatternBackground,
		theme.State { })
	artist.DrawShatter(element.core, pattern, element.Bounds(), rocks...)

	element.partition()
	if parent, ok := element.core.Parent().(elements.ScrollableParent); ok {
		parent.NotifyScrollBoundsChange(element)
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

// NotifyMinimumSizeChange notifies the container that the minimum size of a
// child element has changed.
func (element *DocumentContainer) NotifyMinimumSizeChange (child elements.Element) {
	element.redoAll()
	element.core.DamageAll()
}

// NotifyFlexibleHeightChange notifies the parent that the parameters
// affecting a child's flexible height have changed. This method is
// expected to be called by flexible child element when their content
// changes.
func (element *DocumentContainer) NotifyFlexibleHeightChange (child elements.Flexible) {
	element.redoAll()
	element.core.DamageAll()
}


// SetTheme sets the element's theme.
func (element *DocumentContainer) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.Propagator.SetTheme(new)
	element.redoAll()
}

// SetConfig sets the element's configuration.
func (element *DocumentContainer) SetConfig (new config.Config) {
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
	padding := element.theme.Padding(theme.PatternBackground)
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
	padding := element.theme.Padding(theme.PatternSunken)
	viewportHeight := element.Bounds().Dy() - padding.Vertical()
	height = element.contentBounds.Dy() - viewportHeight
	if height < 0 { height = 0 }
	return
}

// ScrollAxes returns the supported axes for scrolling.
func (element *DocumentContainer) ScrollAxes () (horizontal, vertical bool) {
	return false, true
}

func (element *DocumentContainer) reflectChildProperties () {
	focusable := false
	for _, entry := range element.children {
		_, focusable := entry.Element.(elements.Focusable)
		if focusable {
			focusable = true
			break
		}
	}
	if !focusable && element.Focused() {
		element.Propagator.HandleUnfocus()
	}
}

func (element *DocumentContainer) doLayout () {
	margin := element.theme.Margin(theme.PatternBackground)
	padding := element.theme.Padding(theme.PatternBackground)
	bounds := padding.Apply(element.Bounds())
	element.contentBounds = image.Rectangle { }

	minimumWidth := 0
	dot := bounds.Min.Sub(element.scroll)
	for index, entry := range element.children {
		if index > 0 {
			dot.Y += margin.Y
		}
	
		width, height := entry.MinimumSize()
		if width > minimumWidth {
			minimumWidth = width
		}
		if width < bounds.Dx() {
			width = bounds.Dx()
		}
		if typedChild, ok := entry.Element.(elements.Flexible); ok {
			height = typedChild.FlexibleHeightFor(width)
		}
		
		entry.Bounds.Min = dot
		entry.Bounds.Max = image.Pt(dot.X + width, dot.Y + height)
		element.children[index] = entry
		element.contentBounds = element.contentBounds.Union(entry.Bounds)
		dot.Y += height
	}
	
	element.contentBounds =
		element.contentBounds.Sub(element.contentBounds.Min)	
	element.core.SetMinimumSize (
		minimumWidth + padding.Horizontal(),
		padding.Vertical())
}
