package layouts

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/theme"
// import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type verticalEntry struct {
	y int
	minHeight int
	element tomo.Element
}

// Vertical lays its children out vertically. It can contain any number of
// children. When an child is added to the layout, it can either be set to
// contract to its minimum height or expand to fill the remaining space (space
// that is not taken up by other children or padding is divided equally among
// these). Child elements will all have the same width.
type Vertical struct {
	*core.Core
	core core.CoreControl

	gap, pad bool
	children []verticalEntry
	selectable bool
}

// NewVertical creates a new vertical layout. If gap is set to true, a gap will
// be placed between each child element. If pad is set to true, padding will be
// be placed around the inside of this element's border. Usually, you will want
// these to be true.
func NewVertical (gap, pad bool) (element *Vertical) {
	element = &Vertical { }
	element.Core, element.core = core.NewCore(element)
	element.gap = gap
	element.pad = pad
	element.recalculate()
	return
}

// SetPad sets whether or not padding will be placed around the inside of this
// element's border.
func (element *Vertical) SetPad (pad bool) {
	changed := element.pad != pad
	element.pad = pad
	if changed { element.recalculate() }
}

// SetGap sets whether or not a gap will be placed in between child elements.
func (element *Vertical) SetGap (gap bool) {
	changed := element.gap != gap
	element.gap = gap
	if changed { element.recalculate() }
}

// Adopt adds a child element to the vertical layout. If expand is set to true,
// the element will be expanded to fill a portion of the remaining space in the
// layout.
func (element *Vertical) Adopt (child tomo.Element, expand bool) {
	_, minHeight := child.MinimumSize()
	child.SetParentHooks (tomo.ParentHooks {
		// TODO
	})
	element.children = append (element.children, verticalEntry {
		element: child,
		minHeight: minHeight,
	})
	if child.Selectable() { element.core.SetSelectable(true) }

	element.recalculate()
}

// Disown removes the given child from the layout if it is contained within it.
func (element *Vertical) Disown (child tomo.Element) {
	for index, entry := range element.children {
		if entry.element == child {
			entry.element.SetParentHooks(tomo.ParentHooks { })
			element.children = append (
				element.children[:index],
				element.children[index + 1:]...)
				break
		}
	}

	selectable := false
	for _, entry := range element.children {
		if entry.element.Selectable() { selectable = true }
	}
	element.core.SetSelectable(selectable)
}

// Children returns a slice containing this element's children.
func (element *Vertical) Children () (children []tomo.Element) {
	children = make([]tomo.Element, len(element.children))
	for index, entry := range element.children {
		children[index] = entry.element
	}
	return
}

// CountChildren returns the amount of children contained within this element.
func (element *Vertical) CountChildren () (count int) {
	return len(element.children)
}

// Child returns the child at the specified index. If the index is out of
// bounds, this method will return nil.
func (element *Vertical) Child (index int) (child tomo.Element) {
	if index < 0 || index > len(element.children) { return }
	return element.children[index].element
}

func (element *Vertical) Handle (event tomo.Event) {
	switch event.(type) {
	case tomo.EventResize:
		element.recalculate()
		// TODO:
	
	// TODO:
	}
	return
}

func (element *Vertical) AdvanceSelection (direction int) (ok bool) {
	// TODO:
	return
}

func (element *Vertical) recalculate () {
	var x, y int
	if element.pad {
		x += theme.Padding()
		y += theme.Padding()
	}
	// TODO
}
