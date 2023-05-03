package elements

import "tomo"

type scratchEntry struct {
	expand     bool
	minSize    float64
	minBreadth float64
}

type container struct {
	entity   tomo.Entity
	scratch  map[tomo.Element] scratchEntry
	minimumSize func ()
}

// Entity returns this element's entity.
func (container *container) Entity () tomo.Entity {
	return container.entity
}

// Adopt adds one or more elements to the container.
func (container *container) Adopt (children ...tomo.Element) {
	container.adopt(false, children...)
}

func (container *container) init () {
	container.scratch = make(map[tomo.Element] scratchEntry)
}

func (container *container) adopt (expand bool, children ...tomo.Element) {
	for _, child := range children {
		container.entity.Adopt(child)
		container.scratch[child] = scratchEntry { expand: expand }
	}
	container.minimumSize()
	container.entity.Invalidate()
	container.entity.InvalidateLayout()
}

// Disown removes one or more elements from the container.
func (container *container) Disown (children ...tomo.Element) {
	for _, child := range children {
		index := container.entity.IndexOf(child)
		if index < 0 { continue }
		container.entity.Disown(index)
		delete(container.scratch, child)
	}
	container.minimumSize()
	container.entity.Invalidate()
	container.entity.InvalidateLayout()
}

// DisownAll removes all elements from the container.
func (container *container) DisownAll () {
	func () {
		for index := 0; index < container.entity.CountChildren(); index ++ {
			index := index
			defer container.entity.Disown(index)
		}
	} ()
	container.scratch = make(map[tomo.Element] scratchEntry)
	container.minimumSize()
	container.entity.Invalidate()
	container.entity.InvalidateLayout()
}

// Child returns the child at the specified index.
func (container *container) Child (index int) tomo.Element {
	if index < 0 || index >= container.entity.CountChildren() { return nil }
	return container.entity.Child(index)
}

// CountChildren returns the amount of children in this container.
func (container *container) CountChildren () int {
	return container.entity.CountChildren()
}

func (container *container) HandleChildMinimumSizeChange (child tomo.Element) {
	container.minimumSize()
	container.entity.Invalidate()
	container.entity.InvalidateLayout()
}
