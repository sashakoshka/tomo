package containers

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

type childManager struct {
	onChange func ()
	children []tomo.LayoutEntry
	parent   tomo.Parent
	theme    theme.Wrapped
	config   config.Wrapped
}

// Adopt adds a new child element to the container. If expand is set to true,
// the element will expand (instead of contract to its minimum size), in
// whatever way is defined by the container's layout.
func (manager *childManager) Adopt (child tomo.Element, expand bool) {
	if child0, ok := child.(tomo.Themeable); ok {
		child0.SetTheme(manager.theme.Theme)
	}
	if child0, ok := child.(tomo.Configurable); ok {
		child0.SetConfig(manager.config.Config)
	}
	child.SetParent(manager.parent)

	manager.children = append (manager.children, tomo.LayoutEntry {
		Element: child,
		Expand:  expand,
	})
	
	manager.onChange()
}


// Disown removes the given child from the container if it is contained within
// it.
func (manager *childManager) Disown (child tomo.Element) {
	for index, entry := range manager.children {
		if entry.Element == child {
			manager.clearChildEventHandlers(entry.Element)
			manager.children = append (
				manager.children[:index],
				manager.children[index + 1:]...)
				break
		}
	}

	manager.onChange()
}

// DisownAll removes all child elements from the container at once.
func (manager *childManager) DisownAll () {
	for _, entry := range manager.children {
		manager.clearChildEventHandlers(entry.Element)
	}
	manager.children = nil

	manager.onChange()
}

// Children returns a slice containing this element's children.
func (manager *childManager) Children () (children []tomo.Element) {
	children = make([]tomo.Element, len(manager.children))
	for index, entry := range manager.children {
		children[index] = entry.Element
	}
	return
}

// CountChildren returns the amount of children contained within this element.
func (manager *childManager) CountChildren () (count int) {
	return len(manager.children)
}

// Child returns the child at the specified index. If the index is out of
// bounds, this method will return nil.
func (manager *childManager) Child (index int) (child tomo.Element) {
	if index < 0 || index > len(manager.children) { return }
	return manager.children[index].Element
}

// ChildAt returns the child that contains the specified x and y coordinates. If
// there are no children at the coordinates, this method will return nil.
func (manager *childManager) ChildAt (point image.Point) (child tomo.Element) {
	for _, entry := range manager.children {
		if point.In(entry.Bounds) {
			child = entry.Element
		}
	}
	return
}

func (manager *childManager) clearChildEventHandlers (child tomo.Element) {
	child.DrawTo(nil, image.Rectangle { }, nil)
	child.SetParent(nil)
	
	if child, ok := child.(tomo.Focusable); ok {
		if child.Focused() {
			child.HandleUnfocus()
		}
	}
}
