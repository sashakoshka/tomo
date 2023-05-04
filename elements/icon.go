package elements

import "image"
import "tomo"
import "art"

var iconCase = tomo.C("tomo", "icon")

// Icon is an element capable of displaying a singular icon.
type Icon struct {
	entity tomo.Entity
	id     tomo.Icon
	size   tomo.IconSize
}

// Icon creates a new icon element.
func NewIcon (id tomo.Icon, size tomo.IconSize) (element *Icon) {
	element = &Icon {
		id:   id,
		size: size,
	}
	element.entity = tomo.GetBackend().NewEntity(element)
	element.updateMinimumSize()
	return
}

// Entity returns this element's entity.
func (element *Icon) Entity () tomo.Entity {
	return element.entity
}

// SetIcon sets the element's icon.
func (element *Icon) SetIcon (id tomo.Icon, size tomo.IconSize) {
	element.id   = id
	element.size = size
	if element.entity == nil { return }
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Icon) HandleThemeChange () {
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Icon) Draw (destination art.Canvas) {
	if element.entity == nil { return }
	
	bounds := element.entity.Bounds()
	state  := tomo.State { }
	element.entity.Theme().
		Pattern(tomo.PatternBackground, state, iconCase).
		Draw(destination, bounds)
	icon := element.icon()
	if icon != nil {
		iconBounds := icon.Bounds()
		offset := image.Pt (
			(bounds.Dx() - iconBounds.Dx()) / 2,
			(bounds.Dy() - iconBounds.Dy()) / 2)
		icon.Draw (
			destination,
			element.entity.Theme().Color(tomo.ColorForeground, state, iconCase),
			bounds.Min.Add(offset))
	}
}

func (element *Icon) icon () art.Icon {
	return element.entity.Theme().Icon(element.id, element.size, iconCase)
}

func (element *Icon) updateMinimumSize () {
	icon := element.icon()
	if icon == nil {
		element.entity.SetMinimumSize(0, 0)
	} else {
		bounds := icon.Bounds()
		element.entity.SetMinimumSize(bounds.Dx(), bounds.Dy())
	}
}
