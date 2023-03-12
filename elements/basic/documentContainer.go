package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
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
	
	onFocusRequest       func () (granted bool)
	onFocusMotionRequest func (input.KeynavDirection) (granted bool)
	onScrollBoundsChange func ()
}

// NewDocumentContainer creates a new document container.
func NewDocumentContainer () (element *DocumentContainer) {
	element = &DocumentContainer { }
	element.theme.Case = theme.C("basic", "documentContainer")
	element.Core, element.core = core.NewCore(element.redoAll)
	element.Propagator = core.NewPropagator(element)
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
	child.OnDamage (func (region canvas.Canvas) {
		element.core.DamageRegion(region.Bounds())
	})
	child.OnMinimumSizeChange (func () {
		element.redoAll()
		element.core.DamageAll()
	})
	if child0, ok := child.(elements.Flexible); ok {
		child0.OnFlexibleHeightChange (func () {
			element.redoAll()
			element.core.DamageAll()
		})
	}
	if child0, ok := child.(elements.Focusable); ok {
		child0.OnFocusRequest (func () (granted bool) {
			return element.childFocusRequestCallback(child0)
		})
		child0.OnFocusMotionRequest (
			func (direction input.KeynavDirection) (granted bool) {
				if element.onFocusMotionRequest == nil { return }
				return element.onFocusMotionRequest(direction)
			})
	}

	// add child
	element.children = append (element.children, layouts.LayoutEntry {
		Element: child,
	})

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
	child.DrawTo(nil)
	child.OnDamage(nil)
	child.OnMinimumSizeChange(nil)
	if child0, ok := child.(elements.Focusable); ok {
		child0.OnFocusRequest(nil)
		child0.OnFocusMotionRequest(nil)
		if child0.Focused() {
			child0.HandleUnfocus()
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

	// draw a background
	rocks := make([]image.Rectangle, len(element.children))
	for index, entry := range element.children {
		rocks[index] = entry.Bounds
	}
	pattern := element.theme.Pattern (
		theme.PatternBackground,
		theme.State { })
	artist.DrawShatter (
		element.core, pattern, rocks...)

	element.partition()
	if element.onScrollBoundsChange != nil {
		element.onScrollBoundsChange()
	}
}

func (element *DocumentContainer) partition () {
	for _, entry := range element.children {
		entry.DrawTo(nil)
	}

	// cut our canvas up and give peices to child elements
	for _, entry := range element.children {
		if entry.Bounds.Overlaps(element.Bounds()) {
			entry.DrawTo(canvas.Cut(element.core, entry.Bounds))
		}
	}
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

func (element *DocumentContainer) OnFocusRequest (callback func () (granted bool)) {
	element.onFocusRequest = callback
	element.Propagator.OnFocusRequest(callback)
}

func (element *DocumentContainer) OnFocusMotionRequest (
	callback func (direction input.KeynavDirection) (granted bool),
) {
	element.onFocusMotionRequest = callback
	element.Propagator.OnFocusMotionRequest(callback)
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

// OnScrollBoundsChange sets a function to be called when the element's
// ScrollContentBounds, ScrollViewportBounds, or ScrollAxes are changed.
func (element *DocumentContainer) OnScrollBoundsChange(callback func()) {
	element.onScrollBoundsChange = callback
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

func (element *DocumentContainer) childFocusRequestCallback (
	child elements.Focusable,
) (
	granted bool,
) {
	if element.onFocusRequest != nil && element.onFocusRequest() {
		element.Propagator.HandleUnfocus()
		element.Propagator.HandleFocus(input.KeynavDirectionNeutral)
		return true
	} else {
		return false
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
