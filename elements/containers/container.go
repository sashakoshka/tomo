package containers

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
	
	config config.Wrapped
	theme  theme.Wrapped
	
	onFocusRequest func () (granted bool)
	onFocusMotionRequest func (input.KeynavDirection) (granted bool)
}

// NewContainer creates a new container.
func NewContainer (layout layouts.Layout) (element *Container) {
	element = &Container { }
	element.theme.Case = theme.C("containers", "container")
	element.Core, element.core = core.NewCore(element, element.redoAll)
	element.Propagator = core.NewPropagator(element, element.core)
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
	if child0, ok := child.(elements.Themeable); ok {
		child0.SetTheme(element.theme.Theme)
	}
	if child0, ok := child.(elements.Configurable); ok {
		child0.SetConfig(element.config.Config)
	}
	child.SetParent(element)

	// add child
	element.children = append (element.children, layouts.LayoutEntry {
		Element: child,
		Expand:  expand,
	})

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
	if element.core.HasImage() && !element.warping {
		element.redoAll()
		element.core.DamageAll()
	}
}

func (element *Container) clearChildEventHandlers (child elements.Element) {
	child.DrawTo(nil, image.Rectangle { }, nil)
	child.SetParent(nil)
	
	if child, ok := child.(elements.Focusable); ok {
		if child.Focused() {
			child.HandleUnfocus()
		}
	}
}

// DisownAll removes all child elements from the container at once.
func (element *Container) DisownAll () {
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

	// remove child canvasses so that any operations done in here will not
	// cause a child to draw to a wack ass canvas.
	for _, entry := range element.children {
		entry.DrawTo(nil, entry.Bounds, nil)
	}
	
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
	artist.DrawShatter(element.core, pattern, element.Bounds(), rocks...)

	// cut our canvas up and give peices to child elements
	for _, entry := range element.children {
		entry.DrawTo (
			canvas.Cut(element.core, entry.Bounds),
			entry.Bounds, func (region image.Rectangle) {
				element.core.DamageRegion(region)
			})
	}
}

// NotifyMinimumSizeChange notifies the container that the minimum size of a
// child element has changed.
func (element *Container) NotifyMinimumSizeChange (child elements.Element) {
	element.updateMinimumSize()
	element.redoAll()
	element.core.DamageAll()
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

func (element *Container) updateMinimumSize () {
	margin  := element.theme.Margin(theme.PatternBackground)
	padding := element.theme.Padding(theme.PatternBackground)
	width, height := element.layout.MinimumSize (
		element.children, margin, padding)
	element.core.SetMinimumSize(width, height)
}

func (element *Container) doLayout () {
	margin := element.theme.Margin(theme.PatternBackground)
	padding := element.theme.Padding(theme.PatternBackground)
	element.layout.Arrange (
		element.children, margin,
		padding, element.Bounds())
}
