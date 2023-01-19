package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

// Container is an element capable of containg other elements, and arranging
// them in a layout.
type Container struct {
	*core.Core
	core core.CoreControl

	layout     tomo.Layout
	children   []tomo.LayoutEntry
	drags      [10]tomo.MouseTarget
	warping    bool
	selected   bool
	selectable bool
	flexible   bool
	
	onSelectionRequest func () (granted bool)
	onSelectionMotionRequest func (tomo.SelectionDirection) (granted bool)
	onFlexibleHeightChange func ()
}

// NewContainer creates a new container.
func NewContainer (layout tomo.Layout) (element *Container) {
	element = &Container { }
	element.Core, element.core = core.NewCore(element)
	element.SetLayout(layout)
	return
}

// SetLayout sets the layout of this container.
func (element *Container) SetLayout (layout tomo.Layout) {
	element.layout = layout
	if element.core.HasImage() {
		element.recalculate()
		element.draw()
		element.core.DamageAll()
	}
}

// Adopt adds a new child element to the container. If expand is set to true,
// the element will expand (instead of contract to its minimum size), in
// whatever way is defined by the current layout.
func (element *Container) Adopt (child tomo.Element, expand bool) {
	// set event handlers
	child.OnDamage (func (region tomo.Canvas) {
		element.drawChildRegion(child, region)
	})
	child.OnMinimumSizeChange(element.updateMinimumSize)
	if child0, ok := child.(tomo.Flexible); ok {
		child0.OnFlexibleHeightChange(element.updateMinimumSize)
	}
	if child0, ok := child.(tomo.Selectable); ok {
		child0.OnSelectionRequest (func () (granted bool) {
			return element.childSelectionRequestCallback(child0)
		})
	}
	if child0, ok := child.(tomo.Selectable); ok {
		child0.OnSelectionMotionRequest (
			func (direction tomo.SelectionDirection) (granted bool) {
				if element.onSelectionMotionRequest == nil { return }
				return element.onSelectionMotionRequest(direction)
			})
	}

	// add child
	element.children = append (element.children, tomo.LayoutEntry {
		Element: child,
		Expand:  expand,
	})

	// refresh stale data
	element.updateMinimumSize()
	element.reflectChildProperties()
	if element.core.HasImage() && !element.warping {
		element.recalculate()
		element.draw()
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
		element.recalculate()
		element.draw()
		element.core.DamageAll()
	}
}

// Disown removes the given child from the container if it is contained within
// it.
func (element *Container) Disown (child tomo.Element) {
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
		element.recalculate()
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Container) clearChildEventHandlers (child tomo.Element) {
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

// DisownAll removes all child elements from the container at once.
func (element *Container) DisownAll () {
	element.children = nil

	element.updateMinimumSize()
	element.reflectChildProperties()
	if element.core.HasImage() && !element.warping {
		element.recalculate()
		element.draw()
		element.core.DamageAll()
	}
}

// Children returns a slice containing this element's children.
func (element *Container) Children () (children []tomo.Element) {
	children = make([]tomo.Element, len(element.children))
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
func (element *Container) Child (index int) (child tomo.Element) {
	if index < 0 || index > len(element.children) { return }
	return element.children[index].Element
}

// ChildAt returns the child that contains the specified x and y coordinates. If
// there are no children at the coordinates, this method will return nil.
func (element *Container) ChildAt (point image.Point) (child tomo.Element) {
	for _, entry := range element.children {
		if point.In(entry.Bounds().Add(entry.Position)) {
			child = entry.Element
		}
	}
	return
}

func (element *Container) childPosition (child tomo.Element) (position image.Point) {
	for _, entry := range element.children {
		if entry.Element == child {
			position = entry.Position
			break
		}
	}

	return
}

func (element *Container) Resize (width, height int) {
	element.core.AllocateCanvas(width, height)
	element.recalculate()
	element.draw()
}

func (element *Container) HandleMouseDown (x, y int, button tomo.Button) {
	child, handlesMouse := element.ChildAt(image.Pt(x, y)).(tomo.MouseTarget)
	if !handlesMouse { return }
	element.drags[button] = child
	childPosition := element.childPosition(child)
	child.HandleMouseDown(x - childPosition.X, y - childPosition.Y, button)
}

func (element *Container) HandleMouseUp (x, y int, button tomo.Button) {
	child := element.drags[button]
	if child == nil { return }
	element.drags[button] = nil
	childPosition := element.childPosition(child)
	child.HandleMouseUp(x - childPosition.X, y - childPosition.Y, button)
}

func (element *Container) HandleMouseMove (x, y int) {
	for _, child := range element.drags {
		if child == nil { continue }
		childPosition := element.childPosition(child)
		child.HandleMouseMove(x - childPosition.X, y - childPosition.Y)
	}
}

func (element *Container) HandleScroll (x, y int, deltaX, deltaY float64) {
	child, handlesMouse := element.ChildAt(image.Pt(x, y)).(tomo.MouseTarget)
	if !handlesMouse { return }
	childPosition := element.childPosition(child)
	child.HandleMouseScroll(x - childPosition.X, y - childPosition.Y, deltaX, deltaY)
}

func (element *Container) HandleKeyDown (
	key tomo.Key,
	modifiers tomo.Modifiers,
	repeated bool,
) {
	element.forSelected (func (child tomo.Selectable) bool {
		child0, handlesKeyboard := child.(tomo.KeyboardTarget)
		if handlesKeyboard {
			child0.HandleKeyDown(key, modifiers, repeated)
		}
		return true
	})
}

func (element *Container) HandleKeyUp (key tomo.Key, modifiers tomo.Modifiers) {
	element.forSelected (func (child tomo.Selectable) bool {
		child0, handlesKeyboard := child.(tomo.KeyboardTarget)
		if handlesKeyboard {
			child0.HandleKeyUp(key, modifiers)
		}
		return true
	})
}

func (element *Container) Selected () (selected bool) {
	return element.selected
}

func (element *Container) Select () {
	if element.onSelectionRequest != nil {
		element.onSelectionRequest()
	}
}

func (element *Container) HandleSelection (direction tomo.SelectionDirection) (ok bool) {
	if !element.selectable { return false }
	direction = direction.Canon()

	firstSelected := element.firstSelected()
	if firstSelected < 0 {
		found := false
		switch direction {
		case tomo.SelectionDirectionBackward:
			element.forSelectableBackward (func (child tomo.Selectable) bool {
				if child.HandleSelection(direction) {
					element.selected = true
					found = true
					return false
				}
				return true
			})
			return true
		
		case tomo.SelectionDirectionNeutral, tomo.SelectionDirectionForward:
			element.forSelectable (func (child tomo.Selectable) bool {
				if child.HandleSelection(direction) {
					element.selected = true
					found = true
					return false
				}
				return true
			})
		}
		return found
	} else {
		firstSelectedChild :=
			element.children[firstSelected].Element.(tomo.Selectable)
		
		for index := firstSelected + int(direction);
			index < len(element.children) && index >= 0;
			index += int(direction) {

			child, selectable :=
				element.children[index].
				Element.(tomo.Selectable)
			if selectable && child.HandleSelection(direction) {
				firstSelectedChild.HandleDeselection()
				element.selected = true
				return true
			}
		}
	}
	
	return false
}

func (element *Container) FlexibleHeightFor (width int) (height int) {
	return element.layout.FlexibleHeightFor(element.children, width)
}

func (element *Container) OnFlexibleHeightChange (callback func ()) {
	element.onFlexibleHeightChange = callback
}

func (element *Container) HandleDeselection () {
	element.selected = false
	element.forSelected (func (child tomo.Selectable) bool {
		child.HandleDeselection()
		return true
	})
}

func (element *Container) OnSelectionRequest (callback func () (granted bool)) {
	element.onSelectionRequest = callback
}

func (element *Container) OnSelectionMotionRequest (
	callback func (direction tomo.SelectionDirection) (granted bool),
) {
	element.onSelectionMotionRequest = callback
}

func (element *Container) forSelected (callback func (child tomo.Selectable) bool) {
	for _, entry := range element.children {
		child, selectable := entry.Element.(tomo.Selectable)
		if selectable && child.Selected() {
			if !callback(child) { break }
		}
	}
}

func (element *Container) forSelectable (callback func (child tomo.Selectable) bool) {
	for _, entry := range element.children {
		child, selectable := entry.Element.(tomo.Selectable)
		if selectable {
			if !callback(child) { break }
		}
	}
}

func (element *Container) forFlexible (callback func (child tomo.Flexible) bool) {
	for _, entry := range element.children {
		child, selectable := entry.Element.(tomo.Flexible)
		if selectable {
			if !callback(child) { break }
		}
	}
}

func (element *Container) forSelectableBackward (callback func (child tomo.Selectable) bool) {
	for index := len(element.children) - 1; index >= 0; index -- {
		child, selectable := element.children[index].Element.(tomo.Selectable)
		if selectable {
			if !callback(child) { break }
		}
	}
}

func (element *Container) firstSelected () (index int) {
	for currentIndex, entry := range element.children {
		child, selectable := entry.Element.(tomo.Selectable)
		if selectable && child.Selected() {
			return currentIndex
		}
	}
	return -1
}

func (element *Container) reflectChildProperties () {
	element.selectable = false
	element.forSelectable (func (tomo.Selectable) bool {
		element.selectable = true
		return false
	})
	element.flexible = false
	element.forFlexible (func (tomo.Flexible) bool {
		element.flexible = true
		return false
	})
	if !element.selectable {
		element.selected = false
	}
}

func (element *Container) childSelectionRequestCallback (
	child tomo.Selectable,
) (
	granted bool,
) {
	if element.onSelectionRequest != nil && element.onSelectionRequest() {
		element.forSelected (func (child tomo.Selectable) bool {
			child.HandleDeselection()
			return true
		})
		child.HandleSelection(tomo.SelectionDirectionNeutral)
		return true
	} else {
		return false
	}
}

func (element *Container) updateMinimumSize () {
	width, height := element.layout.MinimumSize(element.children)
	if element.flexible {
		height = element.layout.FlexibleHeightFor(element.children, width)
	}
	element.core.SetMinimumSize(width, height)
}

func (element *Container) recalculate () {
	bounds := element.Bounds()
	element.layout.Arrange(element.children, bounds.Dx(), bounds.Dy())
}

func (element *Container) draw () {
	bounds := element.core.Bounds()

	artist.FillRectangle (
		element.core,
		theme.BackgroundPattern(),
		bounds)

	for _, entry := range element.children {
		artist.Paste(element.core, entry, entry.Position)
	}
}

func (element *Container) drawChildRegion (child tomo.Element, region tomo.Canvas) {
	if element.warping { return }
	for _, entry := range element.children {
		if entry.Element == child {
			artist.Paste(element.core, region, entry.Position)
			element.core.DamageRegion (
				region.Bounds().Add(entry.Position))
			break
		}
	}
}
