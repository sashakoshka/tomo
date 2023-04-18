package elements

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"

type cellEntity interface {
	tomo.ContainerEntity
	tomo.SelectableEntity
}

// Cell is a single-element container that satisfies tomo.Selectable. It
// provides styling based on whether or not it is selected.
type Cell struct {
	entity  cellEntity
	child   tomo.Element
	enabled bool
	theme   theme.Wrapped

	onSelectionChange func ()
}

// NewCell creates a new cell element. If padding is true, the cell will have
// padding on all sides. Child can be nil and added later with the Adopt()
// method.
func NewCell (child tomo.Element) (element *Cell) {
	element = &Cell { enabled: true }
	element.theme.Case = tomo.C("tomo", "cell")
	element.entity = tomo.NewEntity(element).(cellEntity)
	element.Adopt(child)
	return
}

// Entity returns this element's entity.
func (element *Cell) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Cell) Draw (destination canvas.Canvas) {
	bounds  := element.entity.Bounds()
	pattern := element.theme.Pattern(tomo.PatternTableCell, element.state())
	if element.child == nil {
		pattern.Draw(destination, bounds)
	} else {
		artist.DrawShatter (
			destination, pattern, bounds,
			element.child.Entity().Bounds())
	}
}

// Draw causes the element to perform a layout operation.
func (element *Cell) Layout () {
	if element.child == nil { return }
	
	bounds := element.entity.Bounds()
	bounds = element.theme.Padding(tomo.PatternTableCell).Apply(bounds)

	element.entity.PlaceChild(0, bounds)
}

// DrawBackground draws this element's background pattern to the specified
// destination canvas.
func (element *Cell) DrawBackground (destination canvas.Canvas) {
	element.theme.Pattern(tomo.PatternTableCell, element.state()).
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

// SetTheme sets this element's theme.
func (element *Cell) SetTheme (theme tomo.Theme) {
	if theme == element.theme.Theme { return }
	element.theme.Theme = theme
	element.updateMinimumSize()
	element.entity.Invalidate()
	element.invalidateChild()
	element.entity.InvalidateLayout()
}

// OnSelectionChange sets a function to be called when this element is selected
// or unselected.
func (element *Cell) OnSelectionChange (callback func ()) {
	element.onSelectionChange = callback
}

func (element *Cell) Selected () bool {
	return element.entity.Selected()
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
	padding := element.theme.Padding(tomo.PatternTableCell)
	width  += padding.Horizontal()
	height += padding.Vertical()
	
	element.entity.SetMinimumSize(width, height)
}

func (element *Cell) invalidateChild () {
	if element.child != nil {
		element.child.Entity().Invalidate()
	}
}
