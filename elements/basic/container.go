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

// Container is an element capable of containg other elements, and arranging
// them in a layout.
type Container struct {
	*core.Core
	*core.Propagator
	core core.CoreControl

	layout    layouts.Layout
	children  []layouts.LayoutEntry
	warping   bool
	flexible  bool
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onFlexibleHeightChange func ()
	onFocusRequest func () (granted bool)
	onFocusMotionRequest func (input.KeynavDirection) (granted bool)
}

// NewContainer creates a new container.
func NewContainer (layout layouts.Layout) (element *Container) {
	element = &Container { }
	element.theme.Case = theme.C("basic", "container")
	element.Core, element.core = core.NewCore(element.redoAll)
	element.Propagator = core.NewPropagator(element)
	element.SetLayout(layout)
	return
}

// SetLayout sets the layout of this container.
func (element *Container) SetLayout (layout layouts.Layout) {
	element.layout = layout
	if element.core.HasImage() {
		element.redoAll()
		element.core.DamageAll()
	}
}

// Adopt adds a new child element to the container. If expand is set to true,
// the element will expand (instead of contract to its minimum size), in
// whatever way is defined by the current layout.
func (element *Container) Adopt (child elements.Element, expand bool) {
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
		// TODO: this could probably stand to be more efficient. I mean
		// seriously?
		element.updateMinimumSize()
		element.redoAll()
		element.core.DamageAll()
	})
	if child0, ok := child.(elements.Flexible); ok {
		child0.OnFlexibleHeightChange(element.updateMinimumSize)
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
		Expand:  expand,
	})

	// refresh stale data
	element.updateMinimumSize()
	element.reflectChildProperties()
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

// Warp runs the specified callback, deferring all layout and rendering updates
// until the callback has finished executing. This allows for aplications to
// perform batch gui updates without flickering and stuff.
func (element *Container) Warp (callback func ()) {
	if element.warping {
		callback()
		return
	}

	element.warping = true
	callback()
	element.warping = false
	
	// TODO: create some sort of task list so we don't do a full recalculate
	// and redraw every time, because although that is the most likely use
	// case, it is not the only one.
	if element.core.HasImage() {
		element.redoAll()
		element.core.DamageAll()
	}
}

// Disown removes the given child from the container if it is contained within
// it.
func (element *Container) Disown (child elements.Element) {
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
	element.reflectChildProperties()
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

func (element *Container) clearChildEventHandlers (child elements.Element) {
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
	if child0, ok := child.(elements.Flexible); ok {
		child0.OnFlexibleHeightChange(nil)
	}
}

// DisownAll removes all child elements from the container at once.
func (element *Container) DisownAll () {
	element.children = nil

	element.updateMinimumSize()
	element.reflectChildProperties()
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

// Children returns a slice containing this element's children.
func (element *Container) Children () (children []elements.Element) {
	children = make([]elements.Element, len(element.children))
	for index, entry := range element.children {
		children[index] = entry.Element
	}
	return
}

// CountChildren returns the amount of children contained within this element.
func (element *Container) CountChildren () (count int) {
	return len(element.children)
}

// Child returns the child at the specified index. If the index is out of
// bounds, this method will return nil.
func (element *Container) Child (index int) (child elements.Element) {
	if index < 0 || index > len(element.children) { return }
	return element.children[index].Element
}

// ChildAt returns the child that contains the specified x and y coordinates. If
// there are no children at the coordinates, this method will return nil.
func (element *Container) ChildAt (point image.Point) (child elements.Element) {
	for _, entry := range element.children {
		if point.In(entry.Bounds) {
			child = entry.Element
		}
	}
	return
}

func (element *Container) redoAll () {
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

	// cut our canvas up and give peices to child elements
	for _, entry := range element.children {
		entry.DrawTo(canvas.Cut(element.core, entry.Bounds))
	}
}

// SetTheme sets the element's theme.
func (element *Container) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.Propagator.SetTheme(new)
	element.updateMinimumSize()
	element.redoAll()
}

// SetConfig sets the element's configuration.
func (element *Container) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.Propagator.SetConfig(new)
	element.updateMinimumSize()
	element.redoAll()
}

func (element *Container) FlexibleHeightFor (width int) (height int) {
	margin  := element.theme.Margin(theme.PatternBackground)
	padding := element.theme.Padding(theme.PatternBackground)
	return element.layout.FlexibleHeightFor (
		element.children,
		margin, padding, width)
}

func (element *Container) OnFlexibleHeightChange (callback func ()) {
	element.onFlexibleHeightChange = callback
}

func (element *Container) OnFocusRequest (callback func () (granted bool)) {
	element.onFocusRequest = callback
	element.Propagator.OnFocusRequest(callback)
}

func (element *Container) OnFocusMotionRequest (
	callback func (direction input.KeynavDirection) (granted bool),
) {
	element.onFocusMotionRequest = callback
	element.Propagator.OnFocusMotionRequest(callback)
}

func (element *Container) forFlexible (callback func (child elements.Flexible) bool) {
	for _, entry := range element.children {
		child, flexible := entry.Element.(elements.Flexible)
		if flexible {
			if !callback(child) { break }
		}
	}
}

func (element *Container) reflectChildProperties () {
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
	
	element.flexible = false
	element.forFlexible (func (elements.Flexible) bool {
		element.flexible = true
		return false
	})
}

func (element *Container) childFocusRequestCallback (
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

func (element *Container) updateMinimumSize () {
	margin  := element.theme.Margin(theme.PatternBackground)
	padding := element.theme.Padding(theme.PatternBackground)
	width, height := element.layout.MinimumSize (
		element.children, margin, padding)
	if element.flexible {
		height = element.layout.FlexibleHeightFor (
			element.children, margin,
			padding, width)
	}
	element.core.SetMinimumSize(width, height)
}

func (element *Container) doLayout () {
	margin := element.theme.Margin(theme.PatternBackground)
	padding := element.theme.Padding(theme.PatternBackground)
	element.layout.Arrange (
		element.children, margin,
		padding, element.Bounds())
}
