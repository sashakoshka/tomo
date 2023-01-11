package basic

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type Container struct {
	*core.Core
	core core.CoreControl

	layout     tomo.Layout
	children   []tomo.LayoutEntry
	selectable bool

	drags [10]tomo.Element
}

func NewContainer (layout tomo.Layout) (element *Container) {
	element = &Container { }
	element.Core, element.core = core.NewCore(element)
	element.SetLayout(layout)
	return
}

func (element *Container) SetLayout (layout tomo.Layout) {
	element.layout = layout
	if element.core.HasImage() {
		element.recalculate()
		element.draw()
		element.core.PushAll()
	}
}

func (element *Container) Adopt (child tomo.Element, expand bool) {
	child.SetParentHooks (tomo.ParentHooks {
		MinimumSizeChange: func (int, int) {
			element.updateMinimumSize()
		},
		SelectabilityChange: func (bool) {
			element.updateSelectable()
		},
		Draw: func (region tomo.Image) {
			element.drawChildRegion(child, region)
		},
	})
	element.children = append (element.children, tomo.LayoutEntry {
		Element: child,
		Expand:  expand,
	})

	element.updateMinimumSize()
	element.updateSelectable()
	if element.core.HasImage() {
		element.recalculate()
		element.draw()
		element.core.PushAll()
	}
}

// Disown removes the given child from the container if it is contained within
// it.
func (element *Container) Disown (child tomo.Element) {
	for index, entry := range element.children {
		if entry.Element == child {
			entry.SetParentHooks(tomo.ParentHooks { })
			element.children = append (
				element.children[:index],
				element.children[index + 1:]...)
				break
		}
	}

	element.updateMinimumSize()
	element.updateSelectable()
	if element.core.HasImage() {
		element.recalculate()
		element.draw()
		element.core.PushAll()
	}
}

// DisownAll removes all child elements from the container at once.
func (element *Container) DisownAll () {
	element.children = nil

	element.updateMinimumSize()
	element.updateSelectable()
	if element.core.HasImage() {
		element.recalculate()
		element.draw()
		element.core.PushAll()
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

func (element *Container) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		element.recalculate()
		element.draw()

	case tomo.EventMouseDown:
		mouseDownEvent := event.(tomo.EventMouseDown)
		child := element.ChildAt (image.Pt (
			mouseDownEvent.X,
			mouseDownEvent.Y))
		if child == nil { break }
		element.drags[mouseDownEvent.Button] = child
		childPosition := element.childPosition(child)
		child.Handle (tomo.EventMouseDown {
			Button: mouseDownEvent.Button,
			X: mouseDownEvent.X - childPosition.X,
			Y: mouseDownEvent.Y - childPosition.Y,
		})

	case tomo.EventMouseUp:
		mouseUpEvent := event.(tomo.EventMouseUp)
		child := element.drags[mouseUpEvent.Button]
		if child == nil { break }
		element.drags[mouseUpEvent.Button] = nil
		childPosition := element.childPosition(child)
		child.Handle (tomo.EventMouseUp {
			Button: mouseUpEvent.Button,
			X: mouseUpEvent.X - childPosition.X,
			Y: mouseUpEvent.Y - childPosition.Y,
		})

	case tomo.EventMouseMove:
		mouseMoveEvent := event.(tomo.EventMouseMove)
		for _, child := range element.drags {
			if child == nil { continue }
			childPosition := element.childPosition(child)
			child.Handle (tomo.EventMouseMove {
				X: mouseMoveEvent.X - childPosition.X,
				Y: mouseMoveEvent.Y - childPosition.Y,
			})
		}
	}
	return
}

func (element *Container) AdvanceSelection (direction int) (ok bool) {
	// TODO:
	return
}

func (element *Container) updateSelectable () {
	selectable := false
	for _, entry := range element.children {
		if entry.Selectable() { selectable = true }
	}
	element.core.SetSelectable(selectable)
}

func (element *Container) updateMinimumSize () {
	element.core.SetMinimumSize(element.layout.MinimumSize(element.children))
}

func (element *Container) recalculate () {
	bounds := element.Bounds()
	element.layout.Arrange(element.children, bounds.Dx(), bounds.Dy())
}

func (element *Container) draw () {
	bounds := element.core.Bounds()

	artist.Rectangle (
		element.core,
		theme.BackgroundImage(),
		nil, 0,
		bounds)

	for _, entry := range element.children {
		artist.Paste(element.core, entry, entry.Position)
	}
}

func (element *Container) drawChildRegion (child tomo.Element, region tomo.Image) {
	for _, entry := range element.children {
		if entry.Element == child {
			artist.Paste(element.core, region, entry.Position)
			element.core.PushRegion (
				region.Bounds().Add(entry.Position))
			break
		}
	}
}
