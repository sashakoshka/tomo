package basic

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
}

func NewContainer (layout tomo.Layout) (element *Container) {
	element = &Container { }
	element.Core, element.core = core.NewCore(element)
	element.SetLayout(layout)
	return
}

func (element *Container) SetLayout (layout tomo.Layout) {
	element.layout = layout
	element.recalculate()
}

func (element *Container) Adopt (child tomo.Element, expand bool) {
	child.SetParentHooks (tomo.ParentHooks {
		MinimumSizeChange:
			func (int, int) { element.updateMinimumSize() },
		SelectabilityChange:
			func (bool) { element.updateSelectable() },
	})
	element.children = append (element.children, tomo.LayoutEntry {
		Element: child,
	})

	element.updateMinimumSize()
	element.updateSelectable()
	element.recalculate()
	if element.core.HasImage() { element.draw() }
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
	element.recalculate()
	if element.core.HasImage() { element.draw() }
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

func (element *Container) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		resizeEvent := event.(tomo.EventResize)
		element.core.AllocateCanvas (
			resizeEvent.Width,
			resizeEvent.Height)
		element.recalculate()
		element.draw()
	
	// TODO:
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

	// TODO
	for _, entry := range element.children {
		artist.Paste(element.core, entry, entry.Position)
	}
}
