package elements

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/artist"

// Spacer can be used to put space between two elements..
type Spacer struct {
	entity tomo.Entity
	line bool
}

// NewSpacer creates a new spacer.
func NewSpacer () (element *Spacer) {
	element = &Spacer { }
	element.entity = tomo.NewEntity(element).(spacerEntity)
	element.theme.Case = tomo.C("tomo", "spacer")
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
		pattern := element.theme.Pattern (
			tomo.PatternLine,
			tomo.State { })
		pattern.Draw(destination, bounds)
	} else {
		pattern := element.theme.Pattern (
			tomo.PatternBackground,
			tomo.State { })
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

// SetTheme sets the element's theme.
func (element *Spacer) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.entity.Invalidate()
}

// SetConfig sets the element's configuration.
func (element *Spacer) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.entity.Invalidate()
}

func (element *Spacer) updateMinimumSize () {
	if element.line {
		padding := element.theme.Padding(tomo.PatternLine)
		element.entity.SetMinimumSize (
			padding.Horizontal(),
			padding.Vertical())
	} else {
		element.entity.SetMinimumSize(1, 1)
	}
}
