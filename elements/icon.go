package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/canvas"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"

// Icon is an element capable of displaying a singular icon.
type Icon struct {
	entity tomo.Entity
	theme  theme.Wrapped
	id     tomo.Icon
	size   tomo.IconSize
}

// Icon creates a new icon element.
func NewIcon (id tomo.Icon, size tomo.IconSize) (element *Icon) {
	element = &Icon {
		id:   id,
		size: size,
	}
	element.entity = tomo.NewEntity(element)
	element.theme.Case = tomo.C("tomo", "icon")
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

// SetTheme sets the element's theme.
func (element *Icon) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	if element.entity == nil { return }
	element.updateMinimumSize()
	element.entity.Invalidate()
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Icon) Draw (destination canvas.Canvas) {
	if element.entity == nil { return }
	
	bounds := element.entity.Bounds()
	state  := tomo.State { }
	element.theme.
		Pattern(tomo.PatternBackground, state).
		Draw(destination, bounds)
	icon := element.icon()
	if icon != nil {
		iconBounds := icon.Bounds()
		offset := image.Pt (
			(bounds.Dx() - iconBounds.Dx()) / 2,
			(bounds.Dy() - iconBounds.Dy()) / 2)
		icon.Draw (
			destination,
			element.theme.Color(tomo.ColorForeground, state),
			bounds.Min.Add(offset))
	}
}

func (element *Icon) icon () artist.Icon {
	return element.theme.Icon(element.id, element.size)
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
