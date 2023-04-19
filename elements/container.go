package elements

import "git.tebibyte.media/sashakoshka/tomo"

type scratchEntry struct {
	expand     bool
	minSize    float64
	minBreadth float64
}

type container struct {
	entity   tomo.ContainerEntity
	scratch  map[tomo.Element] scratchEntry
	minimumSize func ()
}

func (container *container) Entity () tomo.Entity {
	return container.entity
}

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

func (container *container) Child (index int) tomo.Element {
	if index < 0 || index >= container.entity.CountChildren() { return nil }
	return container.entity.Child(index)
}

func (container *container) CountChildren () int {
	return container.entity.CountChildren()
}

func (container *container) HandleChildMinimumSizeChange (child tomo.Element) {
	container.minimumSize()
	container.entity.Invalidate()
	container.entity.InvalidateLayout()
}
