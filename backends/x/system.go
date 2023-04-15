package x

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

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
	canvas  canvas.BasicCanvas

	theme  theme.Wrapped
	config config.Wrapped

	invalidateIgnore bool
	drawingInvalid   entitySet
	anyLayoutInvalid bool
	
	drags [10]tomo.MouseTarget

	pushFunc func (image.Rectangle)
}

func (system *system) initialize () {
	system.drawingInvalid = make(entitySet)
}

func (system *system) SetTheme (theme tomo.Theme) {
	system.theme.Theme = theme
	system.propagate (func (entity *entity) bool {
		if child, ok := system.child.element.(tomo.Themeable); ok {
			child.SetTheme(theme)
		}
		return true
	})
}

func (system *system) SetConfig (config tomo.Config) {
	system.config.Config = config
	system.propagate (func (entity *entity) bool {
		if child, ok := system.child.element.(tomo.Configurable); ok {
			child.SetConfig(config)
		}
		return true
	})
}

func (system *system) focusNext () {
	// TODO
}

func (system *system) focusPrevious () {
	// TODO
}

func (system *system) propagate (callback func (*entity) bool) {
	if system.child == nil { return }
	system.child.propagate(callback)
}

func (system *system) childAt (point image.Point) *entity {
	if system.child == nil { return nil }
	return system.child.childAt(point)
}

func (system *system) resizeChildToFit () {
	system.child.bounds        = system.canvas.Bounds()
	system.child.clippedBounds = system.child.bounds
	system.child.Invalidate()
	if system.child.isContainer {
		system.child.InvalidateLayout()
	}
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
		entity.element.(tomo.Container).Layout()
		entity.layoutInvalid = false
		force = true
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
		entity.element.Draw (canvas.Cut (
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
