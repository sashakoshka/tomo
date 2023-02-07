package basicElements

import "image"
import "git.tebibyte.media/sashakoshka/tomo/input"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/layouts"
import "git.tebibyte.media/sashakoshka/tomo/elements"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Container is an element capable of containg other elements, and arranging
// them in a layout.
type Container struct {
	*core.Core
	core core.CoreControl

	layout    layouts.Layout
	children  []layouts.LayoutEntry
	drags     [10]elements.MouseTarget
	warping   bool
	focused   bool
	focusable bool
	flexible  bool
	
	onFocusRequest func () (granted bool)
	onFocusMotionRequest func (input.KeynavDirection) (granted bool)
	onFlexibleHeightChange func ()
}

// NewContainer creates a new container.
func NewContainer (layout layouts.Layout) (element *Container) {
	element = &Container { }
	element.Core, element.core = core.NewCore (
		element.redoAll,
		element.handleConfigChange,
		element.handleThemeChange,
		theme.C("basic", "container"))
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
	child.OnDamage (func (region canvas.Canvas) {
		element.core.DamageRegion(region.Bounds())
	})
	child.OnMinimumSizeChange(element.updateMinimumSize)
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

func (element *Container) childPosition (child elements.Element) (position image.Point) {
	for _, entry := range element.children {
		if entry.Element == child {
			position = entry.Bounds.Min
			break
		}
	}

	return
}

func (element *Container) redoAll () {
	// do a layout
	element.recalculate()

	// draw a background
	bounds := element.Bounds()
	pattern := element.core.Pattern (theme.PatternBackground, theme.PatternState { })
	artist.FillRectangle(element, pattern, bounds)

	// cut our canvas up and give peices to child elements
	for _, entry := range element.children {
		entry.DrawTo(canvas.Cut(element, entry.Bounds))
	}
}

func (element *Container) handleConfigChange () {
	for _, child := range element.children {
		if child0, ok := child.Element.(elements.Configurable); ok {
			child0.SetConfig(element.core.Config())
		}
	}
	element.redoAll()
}

func (element *Container) handleThemeChange () {
	for _, child := range element.children {
		if child0, ok := child.Element.(elements.Themeable); ok {
			child0.SetTheme(element.core.Theme())
		}
	}
	element.redoAll()
}

func (element *Container) HandleMouseDown (x, y int, button input.Button) {
	child, handlesMouse := element.ChildAt(image.Pt(x, y)).(elements.MouseTarget)
	if !handlesMouse { return }
	element.drags[button] = child
	child.HandleMouseDown(x, y, button)
}

func (element *Container) HandleMouseUp (x, y int, button input.Button) {
	child := element.drags[button]
	if child == nil { return }
	element.drags[button] = nil
	child.HandleMouseUp(x, y, button)
}

func (element *Container) HandleMouseMove (x, y int) {
	for _, child := range element.drags {
		if child == nil { continue }
		child.HandleMouseMove(x, y)
	}
}

func (element *Container) HandleMouseScroll (x, y int, deltaX, deltaY float64) {
	child, handlesMouse := element.ChildAt(image.Pt(x, y)).(elements.MouseTarget)
	if !handlesMouse { return }
	child.HandleMouseScroll(x, y, deltaX, deltaY)
}

func (element *Container) HandleKeyDown (key input.Key, modifiers input.Modifiers) {
	element.forFocused (func (child elements.Focusable) bool {
		child0, handlesKeyboard := child.(elements.KeyboardTarget)
		if handlesKeyboard {
			child0.HandleKeyDown(key, modifiers)
		}
		return true
	})
}

func (element *Container) HandleKeyUp (key input.Key, modifiers input.Modifiers) {
	element.forFocused (func (child elements.Focusable) bool {
		child0, handlesKeyboard := child.(elements.KeyboardTarget)
		if handlesKeyboard {
			child0.HandleKeyUp(key, modifiers)
		}
		return true
	})
}

func (element *Container) FlexibleHeightFor (width int) (height int) {
	return element.layout.FlexibleHeightFor (
		element.children,
		element.core.Config().Margin(), width)
}

func (element *Container) OnFlexibleHeightChange (callback func ()) {
	element.onFlexibleHeightChange = callback
}

func (element *Container) Focused () (focused bool) {
	return element.focused
}

func (element *Container) Focus () {
	if element.onFocusRequest != nil {
		element.onFocusRequest()
	}
}

func (element *Container) HandleFocus (direction input.KeynavDirection) (ok bool) {
	if !element.focusable { return false }
	direction = direction.Canon()

	firstFocused := element.firstFocused()
	if firstFocused < 0 {
		// no element is currently focused, so we need to focus either
		// the first or last focusable element depending on the
		// direction.
		switch direction {
		case input.KeynavDirectionNeutral, input.KeynavDirectionForward:
			// if we recieve a neutral or forward direction, focus
			// the first focusable element.
			return element.focusFirstFocusableElement(direction)
		
		case input.KeynavDirectionBackward:
			// if we recieve a backward direction, focus the last
			// focusable element.
			return element.focusLastFocusableElement(direction)
		}
	} else {
		// an element is currently focused, so we need to move the
		// focus in the specified direction
		firstFocusedChild :=
			element.children[firstFocused].Element.(elements.Focusable)

		// before we move the focus, the currently focused child
		// may also be able to move its focus. if the child is able
		// to do that, we will let it and not move ours.
		if firstFocusedChild.HandleFocus(direction) {
			return true
		}

		// find the previous/next focusable element relative to the
		// currently focused element, if it exists.
		for index := firstFocused + int(direction);
			index < len(element.children) && index >= 0;
			index += int(direction) {

			child, focusable :=
				element.children[index].
				Element.(elements.Focusable)
			if focusable && child.HandleFocus(direction) {
				// we have found one, so we now actually move
				// the focus.
				firstFocusedChild.HandleUnfocus()
				element.focused = true
				return true
			}
		}
	}
	
	return false
}

func (element *Container) focusFirstFocusableElement (
	direction input.KeynavDirection,
) (
	ok bool,
) {
	element.forFocusable (func (child elements.Focusable) bool {
		if child.HandleFocus(direction) {
			element.focused = true
			ok = true
			return false
		}
		return true
	})
	return
}

func (element *Container) focusLastFocusableElement (
	direction input.KeynavDirection,
) (
	ok bool,
) {
	element.forFocusableBackward (func (child elements.Focusable) bool {
		if child.HandleFocus(direction) {
			element.focused = true
			ok = true
			return false
		}
		return true
	})
	return
}

func (element *Container) HandleUnfocus () {
	element.focused = false
	element.forFocused (func (child elements.Focusable) bool {
		child.HandleUnfocus()
		return true
	})
}

func (element *Container) OnFocusRequest (callback func () (granted bool)) {
	element.onFocusRequest = callback
}

func (element *Container) OnFocusMotionRequest (
	callback func (direction input.KeynavDirection) (granted bool),
) {
	element.onFocusMotionRequest = callback
}

func (element *Container) forFocused (callback func (child elements.Focusable) bool) {
	for _, entry := range element.children {
		child, focusable := entry.Element.(elements.Focusable)
		if focusable && child.Focused() {
			if !callback(child) { break }
		}
	}
}

func (element *Container) forFocusable (callback func (child elements.Focusable) bool) {
	for _, entry := range element.children {
		child, focusable := entry.Element.(elements.Focusable)
		if focusable {
			if !callback(child) { break }
		}
	}
}

func (element *Container) forFlexible (callback func (child elements.Flexible) bool) {
	for _, entry := range element.children {
		child, flexible := entry.Element.(elements.Flexible)
		if flexible {
			if !callback(child) { break }
		}
	}
}

func (element *Container) forFocusableBackward (callback func (child elements.Focusable) bool) {
	for index := len(element.children) - 1; index >= 0; index -- {
		child, focusable := element.children[index].Element.(elements.Focusable)
		if focusable {
			if !callback(child) { break }
		}
	}
}

func (element *Container) firstFocused () (index int) {
	for currentIndex, entry := range element.children {
		child, focusable := entry.Element.(elements.Focusable)
		if focusable && child.Focused() {
			return currentIndex
		}
	}
	return -1
}

func (element *Container) reflectChildProperties () {
	element.focusable = false
	element.forFocusable (func (elements.Focusable) bool {
		element.focusable = true
		return false
	})
	element.flexible = false
	element.forFlexible (func (elements.Flexible) bool {
		element.flexible = true
		return false
	})
	if !element.focusable {
		element.focused = false
	}
}

func (element *Container) childFocusRequestCallback (
	child elements.Focusable,
) (
	granted bool,
) {
	if element.onFocusRequest != nil && element.onFocusRequest() {
		element.forFocused (func (child elements.Focusable) bool {
			child.HandleUnfocus()
			return true
		})
		child.HandleFocus(input.KeynavDirectionNeutral)
		return true
	} else {
		return false
	}
}

func (element *Container) updateMinimumSize () {
	width, height := element.layout.MinimumSize (
		element.children, element.core.Config().Margin())
	if element.flexible {
		height = element.layout.FlexibleHeightFor (
			element.children, element.core.Config().Margin(), width)
	}
	element.core.SetMinimumSize(width, height)
}

func (element *Container) recalculate () {
	element.layout.Arrange (
		element.children, element.core.Config().Margin(), element.Bounds())
}
