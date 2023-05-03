package elements

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/artist/artutil"

var cellCase = tomo.C("tomo", "cell")

// Cell is a single-element container that satisfies tomo.Selectable. It
// provides styling based on whether or not it is selected.
type Cell struct {
	entity  tomo.Entity
	child   tomo.Element
	enabled bool

	onSelectionChange func ()
}

// NewCell creates a new cell element. If padding is true, the cell will have
// padding on all sides. Child can be nil and added later with the Adopt()
// method.
func NewCell (child tomo.Element) (element *Cell) {
	element = &Cell { enabled: true }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.Adopt(child)
	return
}

// Entity returns this element's entity.
func (element *Cell) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Cell) Draw (destination artist.Canvas) {
	bounds  := element.entity.Bounds()
	pattern := element.entity.Theme().Pattern(tomo.PatternTableCell, element.state(), cellCase)
	if element.child == nil {
		pattern.Draw(destination, bounds)
	} else {
		artutil.DrawShatter (
			destination, pattern, bounds,
			element.child.Entity().Bounds())
	}
}

// Draw causes the element to perform a layout operation.
func (element *Cell) Layout () {
	if element.child == nil { return }
	
	bounds := element.entity.Bounds()
	bounds = element.entity.Theme().Padding(tomo.PatternTableCell, cellCase).Apply(bounds)

	element.entity.PlaceChild(0, bounds)
}

// DrawBackground draws this element's background pattern to the specified
// destination canvas.
func (element *Cell) DrawBackground (destination artist.Canvas) {
	element.entity.Theme().Pattern(tomo.PatternTableCell, element.state(), cellCase).
		Draw(destination, element.entity.Bounds())
}

// Adopt sets this element's child. If nil is passed, any child is removed.
func (element *Cell) Adopt (child tomo.Element) {
	if element.child != nil {
		element.entity.Disown(element.entity.IndexOf(element.child))
	}
	if child != nil {
		element.entity.Adopt(child)
	}
	element.child = child

	element.updateMinimumSize()
	element.entity.Invalidate()
	element.invalidateChild()
	element.entity.InvalidateLayout()
}

// Child returns this element's child. If there is no child, this method will
// return nil.
func (element *Cell) Child () tomo.Element {
	return element.child
}

// Enabled returns whether this cell is enabled or not.
func (element *Cell) Enabled () bool {
	return element.enabled
}

// SetEnabled sets whether this cell can be selected or not.
func (element *Cell) SetEnabled (enabled bool) {
	if element.enabled == enabled { return }
	element.enabled = enabled
	element.entity.Invalidate()
	element.invalidateChild()
}

// OnSelectionChange sets a function to be called when this element is selected
// or unselected.
func (element *Cell) OnSelectionChange (callback func ()) {
	element.onSelectionChange = callback
}

func (element *Cell) Selected () bool {
	return element.entity.Selected()
}

func (element *Cell) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.invalidateChild()
	element.entity.InvalidateLayout()
}

func (element *Cell) HandleSelectionChange () {
	element.entity.Invalidate()
	element.invalidateChild()
	if element.onSelectionChange != nil {
		element.onSelectionChange()
	}
}

func (element *Cell) HandleChildMinimumSizeChange (tomo.Element) {
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.entity.InvalidateLayout()
}

func (element *Cell) state () tomo.State {
	return tomo.State {
		Disabled: !element.enabled,
		On:       element.entity.Selected(),
	}
}

func (element *Cell) updateMinimumSize () {
	width, height := 0, 0

	if element.child != nil {
		childWidth, childHeight := element.entity.ChildMinimumSize(0)
		width  += childWidth
		height += childHeight
	}
	padding := element.entity.Theme().Padding(tomo.PatternTableCell, cellCase)
	width  += padding.Horizontal()
	height += padding.Vertical()
	
	element.entity.SetMinimumSize(width, height)
}

func (element *Cell) invalidateChild () {
	if element.child != nil {
		element.child.Entity().Invalidate()
	}
}
