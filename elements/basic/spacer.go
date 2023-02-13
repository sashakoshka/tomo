package basicElements

import "git.tebibyte.media/sashakoshka/tomo/theme"
import "git.tebibyte.media/sashakoshka/tomo/config"
import "git.tebibyte.media/sashakoshka/tomo/artist"
import "git.tebibyte.media/sashakoshka/tomo/elements/core"

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
	element.theme.Case = theme.C("basic", "spacer")
	element.Core, element.core = core.NewCore(element.draw)
	element.core.SetMinimumSize(1, 1)
	return
}

/// SetLine sets whether or not the spacer will appear as a colored line.
func (element *Spacer) SetLine (line bool) {
	if element.line == line { return }
	element.line = line
	if element.core.HasImage() {
		element.draw()
		element.core.DamageAll()
	}
}

// SetTheme sets the element's theme.
func (element *Spacer) SetTheme (new theme.Theme) {
	if new == element.theme.Theme { return }
	element.theme.Theme = new
	element.redo()
}

// SetConfig sets the element's configuration.
func (element *Spacer) SetConfig (new config.Config) {
	if new == element.config.Config { return }
	element.config.Config = new
	element.redo()
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
			theme.PatternForeground,
			theme.PatternState { })
		artist.FillRectangle(element.core, pattern, bounds)
	} else {
		pattern := element.theme.Pattern (
			theme.PatternBackground,
			theme.PatternState { })
		artist.FillRectangle(element.core, pattern, bounds)
	}
}
