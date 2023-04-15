package x

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"

type entity struct {
	window      *window
	parent      *entity
	children    []*entity
	element     tomo.Element
	
	bounds        image.Rectangle
	clippedBounds image.Rectangle
	minWidth      int
	minHeight     int

	layoutInvalid bool
	isContainer   bool
}

func bind (parent *entity, window *window, element tomo.Element) *entity {
	entity := &entity {
		window:  window,
		parent:  parent,
		element: element,
	}
	entity.Invalidate()
	if _, ok := element.(tomo.Container); ok {
		entity.isContainer = true
		entity.InvalidateLayout()
	}

	element.Bind(entity)
	return entity
}

func (entity *entity) unbind () {
	entity.element.Bind(nil)
	for _, childEntity := range entity.children {
		childEntity.unbind()
	}
}

func (entity *entity) propagate (callback func (*entity) bool) {
	for _, child := range entity.children {
		if callback(child) { break }
		child.propagate(callback)
	}
}

func (entity *entity) childAt (point image.Point) *entity {
	for _, child := range entity.children {
		if point.In(child.bounds) {
			return child
		}
	}
	return entity
}

// ----------- Entity ----------- //

func (entity *entity) Invalidate () {
	if entity.window.system.invalidateIgnore { return }
	entity.window.drawingInvalid.Add(entity)
}

func (entity *entity) Bounds () image.Rectangle {
	return entity.bounds
}

func (entity *entity) Window () tomo.Window {
	return entity.window
}

func (entity *entity) SetMinimumSize (width, height int) {
	entity.minWidth  = width
	entity.minHeight = height
	if entity.parent == nil {
		entity.window.setMinimumSize(width, height)
	} else {
		entity.parent.element.(tomo.Container).
			HandleChildMinimumSizeChange(entity.element)
	}
}

func (entity *entity) DrawBackground (destination canvas.Canvas, bounds image.Rectangle) {
	if entity.parent == nil { return }
	entity.parent.element.(tomo.Container).DrawBackground(destination, bounds)
}

// ----------- ContainerEntity ----------- //

func (entity *entity) InvalidateLayout () {
	if !entity.isContainer { return }
	entity.layoutInvalid = true
	entity.window.system.anyLayoutInvalid = true
}

func (entity *entity) Adopt (child tomo.Element) {
	entity.children = append(entity.children, bind(entity, entity.window, child))
}

func (entity *entity) Insert (index int, child tomo.Element) {
	entity.children = append (
		entity.children[:index + 1],
		entity.children[index:]...)
	entity.children[index] = bind(entity, entity.window, child)
}

func (entity *entity) Disown (index int) {
	entity.children[index].unbind()
	entity.children = append (
		entity.children[:index],
		entity.children[index + 1:]...)
}

func (entity *entity) IndexOf (child tomo.Element) int {
	for index, childEntity := range entity.children {
		if childEntity.element == child {
			return index
		}
	}

	return -1
}

func (entity *entity) Child (index int) tomo.Element {
	return entity.children[index].element
}

func (entity *entity) CountChildren () int {
	return len(entity.children)
}

func (entity *entity) PlaceChild (index int, bounds image.Rectangle) {
	child := entity.children[index]
	child.bounds = bounds
	child.clippedBounds = entity.bounds.Intersect(bounds)
	child.Invalidate()
	if child.isContainer {
		child.InvalidateLayout()
	}
}

func (entity *entity) ChildMinimumSize (index int) (width, height int) {
	childEntity := entity.children[index]
	return childEntity.minWidth, childEntity.minHeight
}

// ----------- FocusableEntity ----------- //

func (entity *entity) Focused () bool {
	return entity.window.focused == entity
}

func (entity *entity) Focus () {
	previous := entity.window.focused
	entity.window.focused = entity
	if previous != nil {
		previous.element.(tomo.Focusable).HandleFocusChange()
	}
	entity.element.(tomo.Focusable).HandleFocusChange()
}

func (entity *entity) FocusNext () {
	entity.window.system.focusNext()
}

func (entity *entity) FocusPrevious () {
	entity.window.system.focusPrevious()
}

// ----------- FlexibleEntity ----------- //

func (entity *entity) NotifyFlexibleHeightChange () {
	if entity.parent == nil { return }
	if parent, ok := entity.parent.element.(tomo.FlexibleContainer); ok {
		parent.HandleChildFlexibleHeightChange (
			entity.element.(tomo.Flexible))
	}
}

// ----------- ScrollableEntity ----------- //

func (entity *entity) NotifyScrollBoundsChange () {
	if entity.parent == nil { return }
	if parent, ok := entity.parent.element.(tomo.ScrollableContainer); ok {
		parent.HandleChildScrollBoundsChange (
			entity.element.(tomo.Scrollable))
	}
}
