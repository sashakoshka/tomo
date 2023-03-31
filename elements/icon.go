package elements

import "image"
import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

type Icon struct {
	*core.Core
	core  core.CoreControl
	theme theme.Wrapped
	id    tomo.Icon
	size  tomo.IconSize
}

func NewIcon (id tomo.Icon, size tomo.IconSize) (element *Icon) {
	element = &Icon {
		id:   id,
		size: size,
	}
	element.theme.Case = tomo.C("tomo", "icon")
	element.Core, element.core = core.NewCore(element, element.draw)
	element.updateMinimumSize()
	return
}

func (element *Icon) SetIcon (id tomo.Icon, size tomo.IconSize) {
	element.id   = id
	element.size = size
	element.updateMinimumSize()
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

// SetTheme sets the element's theme.
func (element *Icon) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.updateMinimumSize()
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Icon) icon () artist.Icon {
	return element.theme.Icon(element.id, element.size)
}

func (element *Icon) updateMinimumSize () {
	icon := element.icon()
	if icon == nil {
		element.core.SetMinimumSize(0, 0)
	} else {
		bounds := icon.Bounds()
		element.core.SetMinimumSize(bounds.Dx(), bounds.Dy())
	}
}

func (element *Icon) draw () {
	bounds := element.Bounds()
	state  := tomo.State { }
	element.theme.
		Pattern(tomo.PatternBackground, state).
		Draw(element.core, bounds)
	icon := element.icon()
	if icon != nil {
		iconBounds := icon.Bounds()
		offset := image.Pt (
			(bounds.Dx() - iconBounds.Dx()) / 2,
			(bounds.Dy() - iconBounds.Dy()) / 2)
		icon.Draw (
			element.core,
			element.theme.Color(tomo.ColorForeground, state),
			bounds.Min.Add(offset))
	}
}
