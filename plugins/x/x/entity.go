package x

import "image"
import "tomo"
import "tomo/artist"
import "tomo/ability"

type entity struct {
	backend     *backend
	window      *window
	parent      *entity
	children    []*entity
	element     tomo.Element
	
	bounds        image.Rectangle
	clippedBounds image.Rectangle
	minWidth      int
	minHeight     int

	selected      bool
	layoutInvalid bool
}

func (backend *backend) NewEntity (owner tomo.Element) tomo.Entity {
	entity := &entity { element: owner, backend: backend }
	entity.InvalidateLayout()
	return entity
}

func (ent *entity) unlink () {
	ent.propagate (func (child *entity) bool {
		if child.window != nil {
			delete(ent.window.system.drawingInvalid, child)
		}
		child.window = nil
		return true
	})

	if ent.window != nil {
		delete(ent.window.system.drawingInvalid, ent)
	}
	ent.parent = nil
	ent.window = nil
	
	if element, ok := ent.element.(ability.Selectable); ok {
		ent.selected = false
		element.HandleSelectionChange()
	}
}

func (entity *entity) link (parent *entity) {
	entity.parent = parent
	entity.clip(parent.clippedBounds)
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

func (entity *entity) propagate (callback func (*entity) bool) bool {
	for _, child := range entity.children {
		if !child.propagate(callback) {
			return false
		}
	}
	return callback(entity)
}


func (entity *entity) propagateAlt (callback func (*entity) bool) bool {
	if !callback(entity) {
		return false
	}
	
	for _, child := range entity.children {
		if !child.propagate(callback) {
			return false
		}
	}

	return true
}
func (entity *entity) childAt (point image.Point) *entity {
	for _, child := range entity.children {
		if point.In(child.bounds) {
			return child.childAt(point)
		}
	}
	return entity
}

func (entity *entity) scrollTargetChildAt (point image.Point) *entity {
	for _, child := range entity.children {
		if point.In(child.bounds) {
			result := child.scrollTargetChildAt(point)
			if result != nil { return result }
			break
		}
	}

	if _, ok := entity.element.(ability.ScrollTarget); ok {
		return entity
	}
	return nil
}

func (entity *entity) forMouseTargetContainers (callback func (ability.MouseTargetContainer, tomo.Element)) {
	if entity.parent == nil { return }
	if parent, ok := entity.parent.element.(ability.MouseTargetContainer); ok {
		callback(parent, entity.element)
	}
	entity.parent.forMouseTargetContainers(callback)
}

func (entity *entity) clip (bounds image.Rectangle) {
	entity.clippedBounds = entity.bounds.Intersect(bounds)
	for _, child := range entity.children {
		child.clip(entity.clippedBounds)
	}
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
		entity.parent.element.(ability.Container).
			HandleChildMinimumSizeChange(entity.element)
	}
}

func (entity *entity) DrawBackground (destination artist.Canvas) {
	if entity.parent != nil {
		entity.parent.element.(ability.Container).DrawBackground(destination)
	} else if entity.window != nil {
		entity.backend.theme.Pattern (
			tomo.PatternBackground,
			tomo.State { },
			tomo.C("tomo", "window")).Draw (
				destination,
				entity.window.canvas.Bounds())
	}
}

// ----------- ContainerEntity ----------- //

func (entity *entity) InvalidateLayout () {
	if entity.window == nil { return }
	if _, ok := entity.element.(ability.Layoutable); !ok { return }
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
	child.clip(entity.clippedBounds)
	child.Invalidate()
	child.InvalidateLayout()
}

func (entity *entity) SelectChild (index int, selected bool) {
	child := entity.children[index]
	if element, ok := child.element.(ability.Selectable); ok {
		if child.selected == selected { return }
		child.selected = selected
		element.HandleSelectionChange()
	}
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
	entity.window.system.focus(entity)
}

func (entity *entity) FocusNext () {
	entity.window.system.focusNext()
}

func (entity *entity) FocusPrevious () {
	entity.window.system.focusPrevious()
}

// ----------- SelectableEntity ----------- //

func (entity *entity) Selected () bool {
	return entity.selected
}

// ----------- FlexibleEntity ----------- //

func (entity *entity) NotifyFlexibleHeightChange () {
	if entity.parent == nil { return }
	if parent, ok := entity.parent.element.(ability.FlexibleContainer); ok {
		parent.HandleChildFlexibleHeightChange (
			entity.element.(ability.Flexible))
	}
}

// ----------- ScrollableEntity ----------- //

func (entity *entity) NotifyScrollBoundsChange () {
	if entity.parent == nil { return }
	if parent, ok := entity.parent.element.(ability.ScrollableContainer); ok {
		parent.HandleChildScrollBoundsChange (
			entity.element.(ability.Scrollable))
	}
}

// ----------- ThemeableEntity ----------- //

func (entity *entity) Theme () tomo.Theme {
	return entity.backend.theme
}

// ----------- ConfigurableEntity ----------- //

func (entity *entity) Config () tomo.Config {
	return entity.backend.config
}
