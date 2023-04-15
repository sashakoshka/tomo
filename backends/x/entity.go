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

func (backend *Backend) NewEntity (owner tomo.Element) tomo.Entity {
	entity := &entity { element: owner }
	if _, ok := owner.(tomo.Container); ok {
		entity.isContainer = true
		entity.InvalidateLayout()
	}
	return entity
}

func (ent *entity) unlink () {
	ent.propagate (func (child *entity) bool {
		child.window = nil
		delete(ent.window.system.drawingInvalid, child)
		return true
	})

	if ent.window != nil {
		delete(ent.window.system.drawingInvalid, ent)
	}
	ent.parent = nil
	ent.window = nil
}

func (entity *entity) link (parent *entity) {
	entity.parent = parent
	if parent.window != nil {
		entity.setWindow(parent.window)
	}
}

func (ent *entity) setWindow (window *window) {
	ent.window = window
	ent.Invalidate()
	ent.InvalidateLayout()
	ent.propagate (func (child *entity) bool {
		child.window = window
		ent.Invalidate()
		ent.InvalidateLayout()
		return true
	})
}

func (entity *entity) propagate (callback func (*entity) bool) {
	for _, child := range entity.children {
		if !callback(child) { break }
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
	if entity.window == nil { return }
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
		if entity.window != nil {
			entity.window.setMinimumSize(width, height)
		}
	} else {
		entity.parent.element.(tomo.Container).
			HandleChildMinimumSizeChange(entity.element)
	}
}

func (entity *entity) DrawBackground (destination canvas.Canvas) {
	if entity.parent != nil {
		entity.parent.element.(tomo.Container).DrawBackground(destination)
	} else if entity.window != nil {
		entity.window.system.theme.Pattern (
			tomo.PatternBackground,
			tomo.State { }).Draw (
				destination,
				entity.window.canvas.Bounds())
	}
}

// ----------- ContainerEntity ----------- //

func (entity *entity) InvalidateLayout () {
	if entity.window == nil { return }
	if !entity.isContainer { return }
	entity.layoutInvalid = true
	entity.window.system.anyLayoutInvalid = true
}

func (ent *entity) Adopt (child tomo.Element) {
	childEntity, ok := child.Entity().(*entity)
	if !ok || childEntity == nil { return }
	childEntity.link(ent)
	ent.children = append(ent.children, childEntity)
}

func (ent *entity) Insert (index int, child tomo.Element) {
	childEntity, ok := child.Entity().(*entity)
	if !ok || childEntity == nil { return }
	ent.children = append (
		ent.children[:index + 1],
		ent.children[index:]...)
	ent.children[index] = childEntity
}

func (entity *entity) Disown (index int) {
	entity.children[index].unlink()
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
	child.InvalidateLayout()
}

func (entity *entity) ChildMinimumSize (index int) (width, height int) {
	childEntity := entity.children[index]
	return childEntity.minWidth, childEntity.minHeight
}

// ----------- FocusableEntity ----------- //

func (entity *entity) Focused () bool {
	if entity.window == nil { return false }
	return entity.window.focused == entity
}

func (entity *entity) Focus () {
	if entity.window == nil { return }
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
