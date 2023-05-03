package x

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/ability"
import defaultTheme  "git.tebibyte.media/sashakoshka/tomo/default/theme"
import defaultConfig "git.tebibyte.media/sashakoshka/tomo/default/config"

type entitySet map[*entity] struct { }

func (set entitySet) Empty () bool {
	return len(set) == 0
}

func (set entitySet) Has (entity *entity) bool {
	_, ok := set[entity]
	return ok
}

func (set entitySet) Add (entity *entity) {
	set[entity] = struct { } { }
}

type system struct {
	child   *entity
	focused *entity
	canvas  artist.BasicCanvas

	theme  tomo.Theme
	config tomo.Config

	invalidateIgnore bool
	drawingInvalid   entitySet
	anyLayoutInvalid bool
	
	drags [10]*entity

	pushFunc func (image.Rectangle)
}

func (system *system) initialize () {
	system.drawingInvalid = make(entitySet)
}

func (system *system) setTheme (theme tomo.Theme) {
	if theme == nil {
		system.theme = defaultTheme.Default { }
	} else {
		system.theme = theme
	}
	system.propagate (func (entity *entity) bool {
		if child, ok := system.child.element.(ability.Themeable); ok {
			child.HandleThemeChange()
		}
		return true
	})
}

func (system *system) setConfig (config tomo.Config) {
	if config == nil {
		system.config = defaultConfig.Default { }
	} else {
		system.config = config
	}
	system.propagate (func (entity *entity) bool {
		if child, ok := system.child.element.(ability.Configurable); ok {
			child.HandleConfigChange()
		}
		return true
	})
}

func (system *system) focus (entity *entity) {
	previous := system.focused
	system.focused = entity
	if previous != nil {
		previous.element.(ability.Focusable).HandleFocusChange()
	}
	if entity != nil {
		entity.element.(ability.Focusable).HandleFocusChange()
	}
}

func (system *system) focusNext () {
	found   := system.focused == nil
	focused := false
	system.propagateAlt (func (entity *entity) bool {
		if found {
			// looking for the next element to select
			child, ok := entity.element.(ability.Focusable)
			if ok && child.Enabled() {
				// found it
				entity.Focus()
				focused = true
				return false
			}
		} else {
			// looking for the current focused element
			if entity == system.focused {
				// found it
				found = true
			}
		}
		return true
	})

	if !focused { system.focus(nil) }
}

func (system *system) focusPrevious () {
	var behind *entity
	system.propagate (func (entity *entity) bool {
		if entity == system.focused {
			return false
		}

		child, ok := entity.element.(ability.Focusable)
		if ok && child.Enabled() { behind = entity }
		return true
	})
	system.focus(behind)
}

func (system *system) propagate (callback func (*entity) bool) {
	if system.child == nil { return }
	system.child.propagate(callback)
}

func (system *system) propagateAlt (callback func (*entity) bool) {
	if system.child == nil { return }
	system.child.propagateAlt(callback)
}

func (system *system) childAt (point image.Point) *entity {
	if system.child == nil { return nil }
	return system.child.childAt(point)
}

func (system *system) scrollTargetChildAt (point image.Point) *entity {
	if system.child == nil { return nil }
	return system.child.scrollTargetChildAt(point)
}

func (system *system) resizeChildToFit () {
	system.child.bounds        = system.canvas.Bounds()
	system.child.clippedBounds = system.child.bounds
	system.child.Invalidate()
	system.child.InvalidateLayout()
}

func (system *system) afterEvent () {
	if system.anyLayoutInvalid {
		system.layout(system.child, false)
		system.anyLayoutInvalid = false
	}
	system.draw()
}

func (system *system) layout (entity *entity, force bool) {
	if entity == nil { return }
	if entity.layoutInvalid == true || force {
		if element, ok := entity.element.(ability.Layoutable); ok {
			element.Layout()
			entity.layoutInvalid = false
			force = true
		}
	}

	for _, child := range entity.children {
		system.layout(child, force)
	}
}

func (system *system) draw () {
	finalBounds := image.Rectangle { }

	// ignore invalidations that result from drawing elements, because if an
	// element decides to do that it really needs to rethink its life
	// choices.
	system.invalidateIgnore = true
	defer func () { system.invalidateIgnore = false } ()

	for entity := range system.drawingInvalid {
		if entity.clippedBounds.Empty() { continue }
		entity.element.Draw (artist.Cut (
			system.canvas,
			entity.clippedBounds))
		finalBounds = finalBounds.Union(entity.clippedBounds)
	}
	system.drawingInvalid = make(entitySet)

	// TODO: don't just union all the bounds together, we can definetly
	// consolidateupdated regions more efficiently than this.
	if !finalBounds.Empty() {
		system.pushFunc(finalBounds)
	}
}
