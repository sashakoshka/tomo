package elements

import "git.tebibyte.media/sashakoshka/tomo"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"
import "git.tebibyte.media/sashakoshka/tomo/default/theme"
import "git.tebibyte.media/sashakoshka/tomo/default/config"

// Spacer can be used to put space between two elements..
type Spacer struct {
	*core.Core
	core core.CoreControl
	line bool
	
	config config.Wrapped
	theme  theme.Wrapped
}

// NewSpacer creates a new spacer. If line is set to true, the spacer will be
// filled with a line color, and if compressed to its minimum width or height,
// will appear as a line.
func NewSpacer (line bool) (element *Spacer) {
	element = &Spacer { line: line }
	element.theme.Case = tomo.C("tomo", "spacer")
	element.Core, element.core = core.NewCore(element, element.draw)
	element.updateMinimumSize()
	return
}

/// SetLine sets whether or not the spacer will appear as a colored line.
func (element *Spacer) SetLine (line bool) {
	if element.line == line { return }
	element.line = line
	element.updateMinimumSize()
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

// SetTheme sets the element's theme.
func (element *Spacer) SetTheme (new tomo.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *Spacer) SetConfig (new tomo.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.redo()
}

func (element *Spacer) updateMinimumSize () {
	if element.line {
		padding := element.theme.Padding(tomo.PatternLine)
		element.core.SetMinimumSize (
			padding.Horizontal(),
			padding.Vertical())
	} else {
		element.core.SetMinimumSize(1, 1)
	}
}

func (element *Spacer) redo () {
	if !element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

func (element *Spacer) draw () {
	bounds := element.Bounds()

	if element.line {
		pattern := element.theme.Pattern (
			tomo.PatternLine,
			tomo.State { })
		pattern.Draw(element.core, bounds)
	} else {
		pattern := element.theme.Pattern (
			tomo.PatternBackground,
			tomo.State { })
		pattern.Draw(element.core, bounds)
	}
}
