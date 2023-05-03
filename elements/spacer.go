package elements

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"

var spacerCase = tomo.C("tomo", "spacer")

// Spacer can be used to put space between two elements..
type Spacer struct {
	entity tomo.Entity
	line bool
}

// NewSpacer creates a new spacer.
func NewSpacer () (element *Spacer) {
	element = &Spacer { }
	element.entity = tomo.GetBackend().NewEntity(element)
	element.updateMinimumSize()
	return
}

// NewLine creates a new line separator.
func NewLine () (element *Spacer) {
	element = NewSpacer()
	element.SetLine(true)
	return
}

// Entity returns this element's entity.
func (element *Spacer) Entity () tomo.Entity {
	return element.entity
}

// Draw causes the element to draw to the specified destination canvas.
func (element *Spacer) Draw (destination artist.Canvas) {
	bounds := element.entity.Bounds()

	if element.line {
		pattern := element.entity.Theme().Pattern (
			tomo.PatternLine,
			tomo.State { }, spacerCase)
		pattern.Draw(destination, bounds)
	} else {
		pattern := element.entity.Theme().Pattern (
			tomo.PatternBackground,
			tomo.State { }, spacerCase)
		pattern.Draw(destination, bounds)
	}
}

/// SetLine sets whether or not the spacer will appear as a colored line.
func (element *Spacer) SetLine (line bool) {
	if element.line == line { return }
	element.line = line
	element.updateMinimumSize()
	element.entity.Invalidate()
}

func (element *Spacer) HandleThemeChange () {
	element.entity.Invalidate()
}

func (element *Spacer) updateMinimumSize () {
	if element.line {
		padding := element.entity.Theme().Padding(tomo.PatternLine, spacerCase)
		element.entity.SetMinimumSize (
			padding.Horizontal(),
			padding.Vertical())
	} else {
		element.entity.SetMinimumSize(1, 1)
	}
}
